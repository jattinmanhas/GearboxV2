package domain

import "time"

// PasswordResetToken represents a password reset token in the database
type PasswordResetToken struct {
	ID        uint      `json:"id" db:"id"`
	UserID    uint      `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

