package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewConnection opens a PostgreSQL connection using DATABASE_URL
func NewConnection(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
