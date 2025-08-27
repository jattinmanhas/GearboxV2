package validation

import (
	"fmt"
	"time"
)

// ExampleUsage demonstrates how to use the validation package
func ExampleUsage() {
	// Example 1: Using struct tags for validation
	type User struct {
		Username    string    `validate:"required,username"`
		Email       string    `validate:"required,email"`
		Password    string    `validate:"required,password"`
		Age         int       `validate:"required,min=13,max=120"`
		DateOfBirth time.Time `validate:"required,date_of_birth"`
	}

	user := User{
		Username:    "john_doe",
		Email:       "john@example.com",
		Password:    "SecurePass123",
		Age:         25,
		DateOfBirth: time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Validate the struct
	if errors := ValidateStruct(user); len(errors) > 0 {
		fmt.Printf("Validation errors: %v\n", errors)
	} else {
		fmt.Println("User is valid!")
	}

	// Example 2: Using individual validation functions
	if err := ValidateUsername("jane"); err != nil {
		fmt.Printf("Username error: %v\n", err)
	}

	if err := ValidatePassword("weak"); err != nil {
		fmt.Printf("Password error: %v\n", err)
	}

	if err := ValidateEmail("invalid-email"); err != nil {
		fmt.Printf("Email error: %v\n", err)
	}

	// Example 3: Custom validation with business rules
	type BusinessUser struct {
		CompanyName string `validate:"required,min=2,max=100"`
		TaxID       string `validate:"required,min=9,max=11"`
		Website     string `validate:"omitempty,url"`
		Phone       string `validate:"required,phone"`
		Address     struct {
			Street     string `validate:"required"`
			City       string `validate:"required"`
			State      string `validate:"required"`
			PostalCode string `validate:"required,postal_code"`
			Country    string `validate:"required"`
		} `validate:"required"`
	}

	businessUser := BusinessUser{
		CompanyName: "Acme Corp",
		TaxID:       "12-3456789",
		Website:     "https://acme.com",
		Phone:       "+1-555-123-4567",
		Address: struct {
			Street     string `validate:"required"`
			City       string `validate:"required"`
			State      string `validate:"required"`
			PostalCode string `validate:"required,postal_code"`
			Country    string `validate:"required"`
		}{
			Street:     "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "USA",
		},
	}

	if errors := ValidateStruct(businessUser); len(errors) > 0 {
		fmt.Printf("Business user validation errors: %v\n", errors)
	} else {
		fmt.Println("Business user is valid!")
	}
}

// ExampleValidationErrors shows how to handle validation errors
func ExampleValidationErrors() {
	type Product struct {
		Name        string  `validate:"required,min=1,max=100"`
		Price       float64 `validate:"required,min=0"`
		SKU         string  `validate:"required,min=3,max=20"`
		Description string  `validate:"omitempty,max=500"`
	}

	product := Product{
		Name:  "",    // This will fail required validation
		Price: -10.0, // This will fail min validation
		SKU:   "AB",  // This will fail min validation
	}

	errors := ValidateStruct(product)
	if len(errors) > 0 {
		fmt.Println("Product validation failed:")
		for i, err := range errors {
			fmt.Printf("%d. Field: %s, Error: %s\n", i+1, err.Field, err.Message)
		}
	}
}
