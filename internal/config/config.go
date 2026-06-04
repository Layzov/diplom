package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment (.env optional).
type Config struct {
	Env string `env:"APP_ENV" env-default:"local"`

	HTTP HTTPConfig `env-prefix:"HTTP_"`
	DB   DBConfig   `env-prefix:"POSTGRES_"`
	JWT  JWTConfig  `env-prefix:"JWT_"`
}

type JWTConfig struct {
	Secret    string        `env:"SECRET" env-required:"true"`
	AccessTTL time.Duration `env:"ACCESS_TTL" env-default:"24h"`
}

type HTTPConfig struct {
	Address         string        `env:"ADDRESS" env-default:":8080"`
	Timeout         time.Duration `env:"TIMEOUT" env-default:"30s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"15s"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-default:"localhost"`
	Port     string `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	DBName   string `env:"DB" env-required:"true"`
	SSLMode  string `env:"SSLMODE" env-default:"disable"`
}

// MustLoad reads configuration from environment variables.
// If a file `.env` exists in the working directory, it is loaded first.
func MustLoad() *Config {
	_ = godotenv.Load()

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("config: read env: %v", err)
	}

	return &cfg
}

// DatabaseDSN builds a libpq connection string for GORM / pgx.
func (c *DBConfig) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// WorkingDir returns current working directory (for logs).
func WorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
