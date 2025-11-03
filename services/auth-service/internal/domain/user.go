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
	Password    string         `json:"-" db:"password"` // Now nullable for OAuth users
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

// HasPassword checks if the user has a password set
func (u *User) HasPassword() bool {
	return u.Password != ""
}

// UserAnalytics represents analytics data for users
type UserAnalytics struct {
	TotalUsers            int64                  `json:"total_users"`
	ActiveUsers           int64                  `json:"active_users"`
	NewUsersToday         int64                  `json:"new_users_today"`
	NewUsersThisWeek      int64                  `json:"new_users_this_week"`
	NewUsersThisMonth     int64                  `json:"new_users_this_month"`
	UsersByRole           []UserRoleCount        `json:"users_by_role"`
	UserRegistrationTrend []UserRegistrationData `json:"user_registration_trend"`
}

// UserRoleCount represents user count by role
type UserRoleCount struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

// UserRegistrationData represents user registration data for trends
type UserRegistrationData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
