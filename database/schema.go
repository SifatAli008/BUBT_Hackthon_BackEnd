package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// InitSchema initializes the database schema from schema.sql file
func InitSchema() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	// Read schema.sql file
	schemaPath := filepath.Join("schema.sql")
	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	_, err = DB.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// CheckTableExists checks if a table exists in the database
func CheckTableExists(tableName string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("database connection is not initialized")
	}

	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)
	`

	var exists bool
	err := DB.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check table existence: %w", err)
	}

	return exists, nil
}

// GetTableCount returns the number of tables in the database
func GetTableCount() (int, error) {
	if DB == nil {
		return 0, fmt.Errorf("database connection is not initialized")
	}

	query := `
		SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`

	var count int
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get table count: %w", err)
	}

	return count, nil
}

// ExecuteMigration executes a migration SQL file
func ExecuteMigration(migrationSQL string) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	_, err := DB.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// BeginTransaction starts a new database transaction
func BeginTransaction() (*sql.Tx, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return DB.Begin()
}
