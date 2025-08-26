package domain

import "time"

type PhoneNumber struct {
	ID uint `json:"id" db:"id"`
	UserID uint `json:"user_id" db:"user_id"`
	PhoneType string `json:"phone_type" db:"phone_type"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	IsDefault bool `json:"is_default" db:"is_default"`
	IsVerified bool `json:"is_verified" db:"is_verified"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
}
