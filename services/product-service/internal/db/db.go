package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB wraps the database connection
type DB struct {
	*sqlx.DB
}

// NewDB creates a new database connection
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := cfg.GetDSN()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MinConns)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// Health checks the database connection health
func (db *DB) Health(ctx context.Context) error {
	return db.PingContext(ctx)
}
