package validation

import (
	"testing"
	"time"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid username", "john_doe", false},
		{"valid username with numbers", "user123", false},
		{"valid username with hyphens", "user-name", false},
		{"empty username", "", true},
		{"too short", "ab", true},
		{"too long", "very_long_username_that_exceeds_thirty_characters", true},
		{"contains spaces", "user name", true},
		{"contains special chars", "user@name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "SecurePass123", false},
		{"valid password with special chars", "Secure@Pass123", false},
		{"empty password", "", true},
		{"too short", "Pass1", true},
		{"no uppercase", "securepass123", true},
		{"no lowercase", "SECUREPASS123", true},
		{"no number", "SecurePass", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid email with subdomain", "user@sub.example.com", false},
		{"valid email with plus", "user+tag@example.com", false},
		{"empty email", "", true},
		{"no @ symbol", "userexample.com", true},
		{"no domain", "user@", true},
		{"no TLD", "user@example", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDateOfBirth(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		dob     time.Time
		wantErr bool
	}{
		{"valid DOB", now.AddDate(-25, 0, 0), false},
		{"13 years old", now.AddDate(-13, 0, 0), false},
		{"future date", now.AddDate(1, 0, 0), true},
		{"too young", now.AddDate(-12, 0, 0), true},
		{"too old", now.AddDate(-151, 0, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDateOfBirth(tt.dob)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDateOfBirth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	type TestUser struct {
		Username string `validate:"required,username"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,password"`
	}

	validUser := TestUser{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "SecurePass123",
	}

	invalidUser := TestUser{
		Username: "jo",            // too short
		Email:    "invalid-email", // invalid email
		Password: "weak",          // invalid password
	}

	// Test valid user
	if errors := ValidateStruct(validUser); len(errors) > 0 {
		t.Errorf("ValidateStruct() should not return errors for valid user, got %d errors", len(errors))
	}

	// Test invalid user
	errors := ValidateStruct(invalidUser)
	if len(errors) == 0 {
		t.Error("ValidateStruct() should return errors for invalid user")
	}

	// Check specific errors
	expectedErrors := 3
	if len(errors) != expectedErrors {
		t.Errorf("ValidateStruct() returned %d errors, expected %d", len(errors), expectedErrors)
	}
}

func TestValidationErrors_Error(t *testing.T) {
	errors := ValidationErrors{
		{Field: "username", Message: "username is required"},
		{Field: "email", Message: "email is invalid"},
	}

	expected := "username is required; email is invalid"
	if errors.Error() != expected {
		t.Errorf("ValidationErrors.Error() = %v, want %v", errors.Error(), expected)
	}
}
