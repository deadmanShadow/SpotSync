package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// LoadConfig reads configuration from a .env file (if present) and environment
// variables, applies sensible defaults for development, and returns a populated
// *Config instance.
func LoadConfig() (*Config, error) {
	// Load .env if it exists; ignore the error when the file is not present so
	// production deployments that rely on real env vars still work.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to OS environment variables")
	}

	cfg := &Config{
		Port:        getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretjwtkey"),
	}

	if cfg.DatabaseURL == "" {
		log.Println("WARNING: DATABASE_URL is not set; database-dependent features will fail")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
