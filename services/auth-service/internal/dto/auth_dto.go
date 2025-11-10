package dto

import "time"

type RegisterRequest struct {
	Username    string    `json:"username" validate:"required,username"`
	Password    string    `json:"password" validate:"required,password"`
	Email       string    `json:"email" validate:"required,email"`
	FirstName   string    `json:"first_name" validate:"required,min=1,max=50"`
	MiddleName  string    `json:"middle_name" validate:"omitempty,max=50"`
	LastName    string    `json:"last_name" validate:"omitempty,max=50"` // Made optional
	Avatar      string    `json:"avatar" validate:"omitempty,url"`
	Gender      string    `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"` // Made optional
	DateOfBirth time.Time `json:"date_of_birth" validate:"omitempty,date_of_birth"`                      // Made optional
}

type UpdateUserRequest struct {
	FirstName   string     `json:"first_name" validate:"omitempty,min=1,max=50"`
	MiddleName  string     `json:"middle_name" validate:"omitempty,max=50"`
	LastName    string     `json:"last_name" validate:"omitempty,min=1,max=50"`
	Avatar      string     `json:"avatar" validate:"omitempty,url"`
	Gender      string     `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"omitempty,date_of_birth"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Profile DTOs
type GetProfileResponse struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	FirstName   string     `json:"first_name" validate:"omitempty,min=1,max=50"`
	MiddleName  string     `json:"middle_name" validate:"omitempty,max=50"`
	LastName    string     `json:"last_name" validate:"omitempty,min=1,max=50"`
	PhoneNumber string     `json:"phone_number" validate:"omitempty,phone"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"omitempty,date_of_birth"`
	Avatar      string     `json:"avatar" validate:"omitempty,url"`
}

// Address DTOs

type CreateAddressRequest struct {
	AddressType  string `json:"address_type" validate:"required,oneof=home work billing shipping other"`
	FirstName    string `json:"first_name" validate:"required,min=1,max=100"`
	LastName     string `json:"last_name" validate:"required,min=1,max=100"`
	Company      string `json:"company" validate:"omitempty,max=100"`
	AddressLine1 string `json:"address_line_1" validate:"required,min=1,max=255"`
	AddressLine2 string `json:"address_line_2" validate:"omitempty,max=255"`
	City         string `json:"city" validate:"required,min=1,max=100"`
	State        string `json:"state" validate:"required,min=1,max=100"`
	Country      string `json:"country" validate:"required,min=1,max=100"`
	PostalCode   string `json:"postal_code" validate:"required,min=1,max=20"`
	Phone        string `json:"phone" validate:"omitempty,phone"`
	Email        string `json:"email" validate:"omitempty,email"`
	IsDefault    bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	AddressType  string `json:"address_type" validate:"omitempty,oneof=home work billing shipping other"`
	FirstName    string `json:"first_name" validate:"omitempty,min=1,max=100"`
	LastName     string `json:"last_name" validate:"omitempty,min=1,max=100"`
	Company      string `json:"company" validate:"omitempty,max=100"`
	AddressLine1 string `json:"address_line_1" validate:"omitempty,min=1,max=255"`
	AddressLine2 string `json:"address_line_2" validate:"omitempty,max=255"`
	City         string `json:"city" validate:"omitempty,min=1,max=100"`
	State        string `json:"state" validate:"omitempty,min=1,max=100"`
	Country      string `json:"country" validate:"omitempty,min=1,max=100"`
	PostalCode   string `json:"postal_code" validate:"omitempty,min=1,max=20"`
	Phone        string `json:"phone" validate:"omitempty,phone"`
	Email        string `json:"email" validate:"omitempty,email"`
	IsDefault    *bool  `json:"is_default"`
}

type CreatePhoneNumberRequest struct {
	PhoneType   string `json:"phone_type" validate:"required,oneof=mobile home work fax other"`
	PhoneNumber string `json:"phone_number" validate:"required,phone"`
	CountryCode string `json:"country_code" validate:"required,min=1,max=5"`
	IsPrimary   bool   `json:"is_primary"`
}

type UpdatePhoneNumberRequest struct {
	PhoneType   string `json:"phone_type" validate:"omitempty,oneof=mobile home work fax other"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone"`
	CountryCode string `json:"country_code" validate:"omitempty,min=1,max=5"`
	IsPrimary   *bool  `json:"is_primary"`
}

type SetDefaultAddressRequest struct {
	AddressID uint `json:"address_id" validate:"required"`
}

type SetPrimaryPhoneRequest struct {
	PhoneID uint `json:"phone_id" validate:"required"`
}

// Password Reset DTOs
type ForgotPasswordRequest struct {
	// Either email or username can be provided
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,username"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password"`
}
