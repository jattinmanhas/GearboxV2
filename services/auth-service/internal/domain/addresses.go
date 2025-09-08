package domain

import "time"

// UserAddress represents a user's saved address
type UserAddress struct {
	ID           uint      `json:"id" db:"id"`
	UserID       uint      `json:"user_id" db:"user_id"`
	AddressType  string    `json:"address_type" db:"address_type"` // home, work, billing, shipping, other
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Company      string    `json:"company" db:"company"`
	AddressLine1 string    `json:"address_line_1" db:"address_line_1"`
	AddressLine2 string    `json:"address_line_2" db:"address_line_2"`
	City         string    `json:"city" db:"city"`
	State        string    `json:"state" db:"state"`
	Country      string    `json:"country" db:"country"`
	PostalCode   string    `json:"postal_code" db:"postal_code"`
	Phone        string    `json:"phone" db:"phone"`
	Email        string    `json:"email" db:"email"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsDefault    bool      `json:"is_default" db:"is_default"`
	IsDeleted    bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserPhoneNumber represents a user's phone number
type UserPhoneNumber struct {
	ID          uint      `json:"id" db:"id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	PhoneType   string    `json:"phone_type" db:"phone_type"` // mobile, home, work, fax, other
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	CountryCode string    `json:"country_code" db:"country_code"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	IsPrimary   bool      `json:"is_primary" db:"is_primary"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Address constants
const (
	AddressTypeHome     = "home"
	AddressTypeWork     = "work"
	AddressTypeBilling  = "billing"
	AddressTypeShipping = "shipping"
	AddressTypeOther    = "other"
)

// Phone type constants
const (
	PhoneTypeMobile = "mobile"
	PhoneTypeHome   = "home"
	PhoneTypeWork   = "work"
	PhoneTypeFax    = "fax"
	PhoneTypeOther  = "other"
)
