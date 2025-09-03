package db

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// RunMigrations runs all pending migrations
func (db *DB) RunMigrations(migrationsPath string) error {
	// Create migrations table if it doesn't exist
	if err := db.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := db.loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := db.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if !db.isMigrationApplied(appliedMigrations, migration.Version) {
			log.Printf("Applying migration %d: %s", migration.Version, migration.Name)

			if err := db.applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}

			log.Printf("Successfully applied migration %d: %s", migration.Version, migration.Name)
		}
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *DB) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`

	_, err := db.Exec(query)
	return err
}

// loadMigrations loads all migration files from the migrations directory
func (db *DB) loadMigrations(migrationsPath string) ([]Migration, error) {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".up.sql") {
			continue
		}

		// Extract version and name from filename (e.g., "000001_initial_schema.up.sql")
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".up.sql")

		// Read up migration
		upPath := filepath.Join(migrationsPath, file.Name())
		upSQL, err := os.ReadFile(upPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read up migration %s: %w", upPath, err)
		}

		// Read down migration
		downPath := filepath.Join(migrationsPath, strings.Replace(file.Name(), ".up.sql", ".down.sql", 1))
		downSQL, err := os.ReadFile(downPath)
		if err != nil {
			// Down migration is optional
			downSQL = []byte("-- No down migration available")
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(upSQL),
			DownSQL: string(downSQL),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a list of applied migration versions
func (db *DB) getAppliedMigrations() ([]int, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

// isMigrationApplied checks if a migration version has been applied
func (db *DB) isMigrationApplied(appliedMigrations []int, version int) bool {
	for _, applied := range appliedMigrations {
		if applied == version {
			return true
		}
	}
	return false
}

// applyMigration applies a single migration
func (db *DB) applyMigration(migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute the migration SQL
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record the migration as applied
	query := `INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`
	if _, err := tx.Exec(query, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}
