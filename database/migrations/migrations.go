package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Up      func(*sql.DB) error
	Down    func(*sql.DB) error
}

var migrations []Migration

// RegisterMigration registers a new migration
func RegisterMigration(migration Migration) {
	migrations = append(migrations, migration)
}

// GetMigrations returns all registered migrations
func GetMigrations() []Migration {
	return migrations
}

// CreateMigrationsTable creates the migrations tracking table
func CreateMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// GetAppliedMigrations returns a list of applied migration versions
func GetAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	applied := make(map[int]bool)

	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

// RecordMigration records a migration as applied
func RecordMigration(db *sql.DB, version int, name string) error {
	query := `
		INSERT INTO schema_migrations (version, name) 
		VALUES ($1, $2)
		ON CONFLICT (version) DO NOTHING
	`

	_, err := db.Exec(query, version, name)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// RunMigrations runs all pending migrations
func RunMigrations(db *sql.DB) error {
	if err := CreateMigrationsTable(db); err != nil {
		return err
	}

	applied, err := GetAppliedMigrations(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d (%s) already applied, skipping", migration.Version, migration.Name)
			continue
		}

		log.Printf("Running migration %d (%s)", migration.Version, migration.Name)

		if err := migration.Up(db); err != nil {
			return fmt.Errorf("migration %d (%s) failed: %w", migration.Version, migration.Name, err)
		}

		if err := RecordMigration(db, migration.Version, migration.Name); err != nil {
			return err
		}

		log.Printf("Migration %d (%s) completed successfully", migration.Version, migration.Name)
	}

	return nil
}
