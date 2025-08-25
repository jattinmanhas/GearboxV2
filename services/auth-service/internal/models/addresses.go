package models

import "time"

type Address struct {
	ID uint `json:"id" db:"id"`
	UserID uint `json:"user_id" db:"user_id"`
	AddressType string `json:"address_type" db:"address_type"`
	AddressLine1 string `json:"address_line_1" db:"address_line_1"`
	AddressLine2 string `json:"address_line_2" db:"address_line_2"`
	City string `json:"city" db:"city"`
	State string `json:"state" db:"state"`
	Country string `json:"country" db:"country"`
	PostalCode string `json:"postal_code" db:"postal_code"`
	IsVerified bool `json:"is_verified" db:"is_verified"`
	IsDefault bool `json:"is_default" db:"is_default"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
}