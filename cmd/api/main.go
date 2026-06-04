package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"diplom/internal/auth"
	"diplom/internal/config"
	"diplom/internal/db"
	"diplom/internal/handler"
	"diplom/internal/service"
	slogpretty "diplom/pkg/handlers/slogPretty"
	"diplom/pkg/middleware/mwLogger"
	"diplom/pkg/sl"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting API",
		slog.String("env", cfg.Env),
		slog.String("cwd", config.WorkingDir()),
	)

	database, err := db.Open(cfg.DB.DatabaseDSN())
	if err != nil {
		log.Error("database init failed", sl.Err(err))
		os.Exit(1)
	}
	if err := db.Ping(database); err != nil {
		log.Error("database ping failed", sl.Err(err))
		os.Exit(1)
	}
	log.Info("database connected",
		slog.String("host", cfg.DB.Host),
		slog.String("db", cfg.DB.DBName),
	)

	if err := db.Migrate(database); err != nil {
		log.Error("database migrate failed", sl.Err(err))
		os.Exit(1)
	}
	log.Info("database schema migrated")

	sqlDB, err := database.DB()
	if err != nil {
		log.Error("database pool", sl.Err(err))
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Error("database close failed", sl.Err(err))
		}
	}()

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(chimw.Recoverer)
	r.Use(chimw.URLFormat)

	r.Get("/health", handler.Health())
	
	// Mount Swagger docs handler at /docs path
	// This will handle both /docs and /docs/* paths
	r.Get("/docs", handler.DocsHandler())
	r.Get("/docs/*", handler.DocsHandler())

	jwtMgr, err := auth.NewManager(auth.JWTConfig{
		Secret:    cfg.JWT.Secret,
		AccessTTL: cfg.JWT.AccessTTL,
	})
	if err != nil {
		log.Error("jwt init failed", sl.Err(err))
		os.Exit(1)
	}

	services := service.New(database, jwtMgr)
	handler.RegisterAPI(r, services, jwtMgr)
	log.Info("api routes registered", slog.String("prefix", "/api/v1"))

	srv := &http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.Timeout,
		WriteTimeout: cfg.HTTP.Timeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("http listening", slog.String("addr", cfg.HTTP.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Info("shutdown signal", slog.String("signal", sig.String()))
	case err := <-errCh:
		if err != nil {
			log.Error("http server error", sl.Err(err))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("http shutdown", sl.Err(err))
	} else {
		log.Info("http stopped")
	}
}

func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
		}
		return slog.New(opts.NewPrettyHandler(os.Stdout))
	case envDev:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
