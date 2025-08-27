package validation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("username", validateUsername)
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("email", validateEmail)
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("postal_code", validatePostalCode)
	_ = validate.RegisterValidation("date_of_birth", validateDateOfBirth)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors is a slice of validation errors
type ValidationErrors []ValidationError

// Error returns the error message
func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// ValidateStruct validates a struct using go-playground/validator
func ValidateStruct(s interface{}) ValidationErrors {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		value := err.Value()

		message := getErrorMessage(field, tag, value)

		errors = append(errors, ValidationError{
			Field:   field,
			Tag:     tag,
			Value:   fmt.Sprintf("%v", value),
			Message: message,
		})
	}

	return errors
}

// Custom validation functions
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Username must be 3-30 characters, alphanumeric + underscore + hyphen
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	// Check if username contains only allowed characters
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Password must be at least 6 characters
	if len(password) < 6 {
		return false
	}

	// Password must contain at least one uppercase, one lowercase, and one number
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Remove all non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Phone number should be 10-15 digits
	return len(digits) >= 10 && len(digits) <= 15
}

func validatePostalCode(fl validator.FieldLevel) bool {
	postalCode := fl.Field().String()

	// Basic postal code validation (5 digits for US, more flexible for international)
	postalRegex := regexp.MustCompile(`^[0-9]{5}(-[0-9]{4})?$`)
	return postalRegex.MatchString(postalCode)
}

func validateDateOfBirth(fl validator.FieldLevel) bool {
	dob := fl.Field().Interface().(time.Time)

	// Date of birth should not be in the future
	if dob.After(time.Now()) {
		return false
	}

	// Date of birth should not be more than 150 years ago
	minDate := time.Now().AddDate(-150, 0, 0)
	if dob.Before(minDate) {
		return false
	}

	// User should be at least 13 years old
	minAge := time.Now().AddDate(-13, 0, 0)
	return dob.Before(minAge)
}

// getErrorMessage returns a user-friendly error message
func getErrorMessage(field, tag string, value interface{}) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %v characters", field, value)
	case "max":
		return fmt.Sprintf("%s must be at most %v characters", field, value)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "username":
		return fmt.Sprintf("%s must be 3-30 characters and contain only letters, numbers, underscores, and hyphens", field)
	case "password":
		return fmt.Sprintf("%s must be at least 6 characters with uppercase, lowercase, and number", field)
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "postal_code":
		return fmt.Sprintf("%s must be a valid postal code", field)
	case "date_of_birth":
		return fmt.Sprintf("%s must be a valid date and user must be at least 13 years old", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Manual validation functions for specific use cases
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 6 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email must be a valid email address")
	}

	return nil
}

func ValidateDateOfBirth(dob time.Time) error {
	if dob.IsZero() {
		return fmt.Errorf("date of birth is required")
	}

	if dob.After(time.Now()) {
		return fmt.Errorf("date of birth cannot be in the future")
	}

	minAge := time.Now().AddDate(-13, 0, 0)
	if dob.After(minAge) {
		return fmt.Errorf("user must be at least 13 years old")
	}

	maxAge := time.Now().AddDate(-150, 0, 0)
	if dob.Before(maxAge) {
		return fmt.Errorf("date of birth is invalid")
	}

	return nil
}
