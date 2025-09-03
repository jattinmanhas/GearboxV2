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
