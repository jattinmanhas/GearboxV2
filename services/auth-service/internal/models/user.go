package models

import "time"

type User struct {
	ID uint `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	Email string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
	LastName string `json:"last_name" db:"last_name"`
	Avatar string `json:"avatar" db:"avatar"`
	Gender string `json:"gender" db:"gender"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool `json:"is_deleted" db:"is_deleted"`
}

