package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"diplom/internal/config"
	"diplom/internal/db"
	"diplom/internal/seeder"
	slogpretty "diplom/pkg/handlers/slogPretty"
)

func main() {
	// Parse flags
	users := flag.Int("users", 5, "number of users to generate")
	subjects := flag.Int("subjects", 2, "subjects per user")
	topics := flag.Int("topics", 3, "topics per subject")
	tasks := flag.Int("tasks", 5, "tasks per topic")
	sessions := flag.Int("sessions", 5, "sessions per user")
	flag.Parse()

	// Setup logger
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
	}
	log := slog.New(opts.NewPrettyHandler(os.Stdout))

	log.Info("seed: starting",
		slog.Int("users", *users),
		slog.Int("subjects_per_user", *subjects),
		slog.Int("topics_per_subject", *topics),
		slog.Int("tasks_per_topic", *tasks),
		slog.Int("sessions_per_user", *sessions),
	)

	// Load config
	cfg := config.MustLoad()
	log.Info("config loaded", slog.String("env", cfg.Env))

	// Connect to database
	database, err := db.Open(cfg.DB.DatabaseDSN())
	if err != nil {
		log.Error("database init failed", "error", err)
		os.Exit(1)
	}

	if err := db.Ping(database); err != nil {
		log.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	log.Info("database connected")

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Error("database migrate failed", "error", err)
		os.Exit(1)
	}
	log.Info("database schema migrated")

	// Seed data
	seedCfg := seeder.Config{
		Users:            *users,
		SubjectsPerUser:  *subjects,
		TopicsPerSubject: *topics,
		TasksPerTopic:    *tasks,
		SessionsPerUser:  *sessions,
	}

	if err := seeder.Seed(database, seedCfg, log); err != nil {
		log.Error("seed failed", "error", err)
		os.Exit(1)
	}

	log.Info("seed: completed successfully")
	fmt.Println("\n✅ Database seeded successfully!")
	fmt.Printf("📊 Generated data:\n")
	fmt.Printf("  - %d users\n", *users)
	fmt.Printf("  - %d subjects per user\n", *subjects)
	fmt.Printf("  - %d topics per subject\n", *topics)
	fmt.Printf("  - %d tasks per topic\n", *tasks)
	fmt.Printf("  - %d sessions per user\n", *sessions)
	fmt.Printf("\n💡 Next: run 'go run ./cmd/api' to start the server\n")
}
