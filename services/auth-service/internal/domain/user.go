package domain

import (
	"database/sql"
	"time"
)

// Helper function to create sql.NullString from string
func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Helper function to create sql.NullTime from time.Time
func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

type User struct {
	ID          uint           `json:"id" db:"id"`
	Username    string         `json:"username" db:"username"`
	Password    string         `json:"-" db:"password"`
	Email       string         `json:"email" db:"email"`
	FirstName   string         `json:"first_name" db:"first_name"`
	MiddleName  sql.NullString `json:"middle_name" db:"middle_name"`
	LastName    sql.NullString `json:"last_name" db:"last_name"`
	PhoneNumber sql.NullString `json:"phone_number" db:"phone_number"`
	Avatar      sql.NullString `json:"avatar" db:"avatar"`
	Gender      sql.NullString `json:"gender" db:"gender"`
	DateOfBirth sql.NullTime   `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
	IsDeleted   bool           `json:"is_deleted" db:"is_deleted"`
	IsActive    bool           `json:"is_active" db:"is_active"`

	// Role information
	RoleID uint   `json:"role_id" db:"role_id"`
	Role   string `json:"role" db:"-"` // Role name, populated when needed
}
