package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"spotsync-backend/config"
)

// ConnectDB opens a GORM connection to PostgreSQL using the supplied
// DATABASE_URL, configures a sensible connection pool, and returns the
// *gorm.DB handle.
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not configured")
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve underlying sql.DB: %w", err)
	}

	// Connection pool tuning
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)

	log.Println("Database connection established successfully")

	return db, nil
}
