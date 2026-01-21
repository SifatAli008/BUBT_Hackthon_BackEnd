package database

import (
	"database/sql"
	"fmt"
	"log"

	"foodlink_backend/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// Init initializes the database connection
func Init(cfg *config.Config) error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("database URL is not configured")
	}

	var err error
	DB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	if cfg.DBMaxPoolSize > 0 {
		DB.SetMaxOpenConns(cfg.DBMaxPoolSize)
	}
	if cfg.DBMaxIdleConns > 0 {
		DB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxLifetime > 0 {
		DB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	}
	if cfg.DBIdleTimeout > 0 {
		DB.SetConnMaxIdleTime(cfg.DBIdleTimeout)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// HealthCheck checks the database connection health
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return DB.Ping()
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}
