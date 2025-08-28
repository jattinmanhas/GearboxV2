package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/validation"
)

func main() {
	fmt.Println("🚀 Validation System Demo")
	fmt.Println("========================")

	// Demo 1: Valid user registration
	fmt.Println("1. Valid User Registration:")
	validUser := struct {
		Username    string    `json:"username" validate:"required,username"`
		Email       string    `json:"email" validate:"required,email"`
		Password    string    `json:"password" validate:"required,password"`
		FirstName   string    `json:"first_name" validate:"required,min=1,max=50"`
		LastName    string    `json:"last_name" validate:"required,min=1,max=50"`
		DateOfBirth time.Time `json:"date_of_birth" validate:"required,date_of_birth"`
	}{
		Username:    "john_doe",
		Email:       "john@example.com",
		Password:    "SecurePass123",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if errors := validation.ValidateStruct(validUser); len(errors) > 0 {
		fmt.Printf("❌ Validation failed: %v\n", errors)
	} else {
		fmt.Println("✅ User validation passed!")
	}

	// Demo 2: Invalid user registration
	fmt.Println("\n2. Invalid User Registration:")
	invalidUser := struct {
		Username    string    `json:"username" validate:"required,username"`
		Email       string    `json:"email" validate:"required,email"`
		Password    string    `json:"password" validate:"required,password"`
		FirstName   string    `json:"first_name" validate:"required,min=1,max=50"`
		LastName    string    `json:"last_name" validate:"required,min=1,max=50"`
		DateOfBirth time.Time `json:"date_of_birth" validate:"required,date_of_birth"`
	}{
		Username:    "jo",            // too short
		Email:       "invalid-email", // invalid email
		Password:    "weak",          // invalid password
		FirstName:   "",              // empty
		LastName:    "Doe",
		DateOfBirth: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // future date
	}

	if errors := validation.ValidateStruct(invalidUser); len(errors) > 0 {
		fmt.Println("❌ Validation failed with the following errors:")
		for i, err := range errors {
			fmt.Printf("   %d. %s: %s\n", i+1, err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ User validation passed!")
	}

	// Demo 3: Individual field validation
	fmt.Println("\n3. Individual Field Validation:")

	// Username validation
	if err := validation.ValidateUsername(""); err != nil {
		fmt.Printf("❌ Username validation: %v\n", err)
	}
	if err := validation.ValidateUsername("john_doe"); err != nil {
		fmt.Printf("❌ Username validation: %v\n", err)
	} else {
		fmt.Println("✅ Username validation passed!")
	}

	// Password validation
	if err := validation.ValidatePassword("weak"); err != nil {
		fmt.Printf("❌ Password validation: %v\n", err)
	}
	if err := validation.ValidatePassword("SecurePass123"); err != nil {
		fmt.Printf("❌ Password validation: %v\n", err)
	} else {
		fmt.Println("✅ Password validation passed!")
	}

	// Email validation
	if err := validation.ValidateEmail("invalid-email"); err != nil {
		fmt.Printf("❌ Email validation: %v\n", err)
	}
	if err := validation.ValidateEmail("user@example.com"); err != nil {
		fmt.Printf("❌ Email validation: %v\n", err)
	} else {
		fmt.Println("✅ Email validation passed!")
	}

	// Demo 4: Business user validation
	fmt.Println("\n4. Business User Validation:")
	businessUser := struct {
		CompanyName string `json:"company_name" validate:"required,min=2,max=100"`
		TaxID       string `json:"tax_id" validate:"required,min=9,max=11"`
		Website     string `json:"website" validate:"omitempty,url"`
		Phone       string `json:"phone" validate:"required,phone"`
		Address     struct {
			Street     string `json:"street" validate:"required"`
			City       string `json:"city" validate:"required"`
			State      string `json:"state" validate:"required"`
			PostalCode string `json:"postal_code" validate:"required,postal_code"`
			Country    string `json:"country" validate:"required"`
		} `validate:"required"`
	}{
		CompanyName: "Acme Corp",
		TaxID:       "12-3456789",
		Website:     "https://acme.com",
		Phone:       "+1-555-123-4567",
		Address: struct {
			Street     string `json:"street" validate:"required"`
			City       string `json:"city" validate:"required"`
			State      string `json:"state" validate:"required"`
			PostalCode string `json:"postal_code" validate:"required,postal_code"`
			Country    string `json:"country" validate:"required"`
		}{
			Street:     "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "USA",
		},
	}

	if errors := validation.ValidateStruct(businessUser); len(errors) > 0 {
		fmt.Println("❌ Business user validation failed:")
		for i, err := range errors {
			fmt.Printf("   %d. %s: %s\n", i+1, err.Field, err.Message)
		}
	} else {
		fmt.Println("✅ Business user validation passed!")
	}

	// Demo 5: Error handling
	fmt.Println("\n5. Error Handling Example:")
	type Product struct {
		Name        string  `json:"name" validate:"required,min=1,max=100"`
		Price       float64 `json:"price" validate:"required,min=0"`
		SKU         string  `json:"sku" validate:"required,min=3,max=20"`
		Description string  `json:"description" validate:"omitempty,max=500"`
	}

	product := Product{
		Name:  "",    // This will fail required validation
		Price: -10.0, // This will fail min validation
		SKU:   "AB",  // This will fail min validation
	}

	errors := validation.ValidateStruct(product)
	if len(errors) > 0 {
		fmt.Println("❌ Product validation failed:")

		// Convert to JSON for API response
		errorJSON, _ := json.MarshalIndent(errors, "", "  ")
		fmt.Printf("JSON Error Response:\n%s\n", string(errorJSON))

		// Or use the error message
		fmt.Printf("Error Message: %s\n", errors.Error())
	} else {
		fmt.Println("✅ Product validation passed!")
	}

	fmt.Println("\n🎉 Validation Demo Complete!")
}
