package domain

import (
	"time"
)

// Helper function to create *string from string
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function to create *time.Time from time.Time
func TimePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

type User struct {
	ID          uint       `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	Password    string     `json:"-" db:"password"` // Now nullable for OAuth users
	Email       string     `json:"email" db:"email"`
	FirstName   string     `json:"first_name" db:"first_name"`
	MiddleName  *string    `json:"middle_name" db:"middle_name"`
	LastName    *string    `json:"last_name" db:"last_name"`
	PhoneNumber *string    `json:"phone_number" db:"phone_number"`
	Avatar      *string    `json:"avatar" db:"avatar"`
	Gender      *string    `json:"gender" db:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	IsDeleted   bool       `json:"is_deleted" db:"is_deleted"`
	IsActive    bool       `json:"is_active" db:"is_active"`

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
