package db

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect opens a PostgreSQL connection using DATABASE_URL
func Connect() (*sqlx.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	return sqlx.Connect("postgres", dsn)
}
