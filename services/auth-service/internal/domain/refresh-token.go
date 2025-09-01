package domain

import "time"

type RefreshToken struct {
	ID           uint      `json:"id" db:"id"`
	UserID       uint      `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	UserAgent    string    `json:"user_agent" db:"user_agent"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	IsRevoked    bool      `json:"is_revoked" db:"is_revoked"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}