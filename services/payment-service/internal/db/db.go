package db

import (
	"fmt"
	"log"

	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB holds the database connection
var DB *sqlx.DB

// Connect establishes a connection to the database
func Connect(cfg *config.Config) error {
	dsn := cfg.GetDatabaseDSN()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db
	log.Println("Database connected successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return DB
}
