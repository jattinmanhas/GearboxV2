package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	_ = validate.RegisterValidation("name", validateName)
	_ = validate.RegisterValidation("slug", validateSlug)
	_ = validate.RegisterValidation("description", validateDescription)
	_ = validate.RegisterValidation("image_url", validateImageURL)
	_ = validate.RegisterValidation("meta_title", validateMetaTitle)
	_ = validate.RegisterValidation("meta_description", validateMetaDescription)
	_ = validate.RegisterValidation("sort_order", validateSortOrder)

	// Product validation functions
	_ = validate.RegisterValidation("sku", validateSKU)
	_ = validate.RegisterValidation("price", validatePrice)
	_ = validate.RegisterValidation("weight", validateWeight)
	_ = validate.RegisterValidation("dimensions", validateDimensions)
	_ = validate.RegisterValidation("tags", validateTags)

	// Cart validation functions
	_ = validate.RegisterValidation("currency", validateCurrency)
	_ = validate.RegisterValidation("session_id", validateSessionID)
	_ = validate.RegisterValidation("coupon_code", validateCouponCode)
	_ = validate.RegisterValidation("shipping_method", validateShippingMethod)
}

type ValidatorError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidatorErrors []ValidatorError

func (v ValidatorErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

func ValidateStruct(s any) ValidatorErrors {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors ValidatorErrors
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		value := err.Value()

		message := getErrorMessage(field, tag, value)

		errors = append(errors, ValidatorError{
			Field:   field,
			Tag:     tag,
			Value:   fmt.Sprintf("%v", value),
			Message: message,
		})
	}

	return errors
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if strings.TrimSpace(name) == "" {
		return false
	}

	if len(name) < 1 || len(name) > 255 {
		return false
	}

	// Check for invalid characters
	if strings.ContainsAny(name, "<>\"'&") {
		return false
	}

	return true
}

// validateSlug validates category slug
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	if strings.TrimSpace(slug) == "" {
		return false
	}

	if len(slug) < 1 {
		return false
	}

	if len(slug) > 255 {
		return false
	}

	// Slug should only contain lowercase letters, numbers, and hyphens
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(slug) {
		return false
	}

	// Slug should not start or end with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return false
	}

	// Slug should not contain consecutive hyphens
	if strings.Contains(slug, "--") {
		return false
	}

	return true
}

// validateDescription validates category description
func validateDescription(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) <= 2000
}

// validateImageURL validates image URL
func validateImageURL(fl validator.FieldLevel) bool {
	imageURL := fl.Field().String()
	if imageURL == "" {
		return true // Optional field
	}

	if len(imageURL) > 500 {
		return false
	}

	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

	return urlRegex.MatchString(imageURL)
}

// validateMetaTitle validates meta title
func validateMetaTitle(fl validator.FieldLevel) bool {
	metaTitle := fl.Field().String()
	if metaTitle == "" {
		return true
	}

	if len(metaTitle) > 255 {
		return false
	}

	return true
}

// validateMetaDescription validates meta description
func validateMetaDescription(fl validator.FieldLevel) bool {
	metaDescription := fl.Field().String()
	if metaDescription == "" {
		return true
	}

	if len(metaDescription) > 500 {
		return false
	}

	return true
}

// validateSortOrder validates sort order
func validateSortOrder(fl validator.FieldLevel) bool {
	sortOrder := fl.Field().Int()

	return sortOrder >= 0 && sortOrder <= 999999
}

func getErrorMessage(field, tag string, value interface{}) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %v characters", field, value)
	case "max":
		return fmt.Sprintf("%s must be at most %v characters", field, value)
	case "name":
		return fmt.Sprintf("%s must be between 1 and 255 characters long", field)
	case "slug":
		return fmt.Sprintf("%s must be between 1 and 255 characters long", field)
	case "description":
		return fmt.Sprintf("%s must be at most 2000 characters long", field)
	case "image_url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "meta_title":
		return fmt.Sprintf("%s must be at most 255 characters long", field)
	case "meta_description":
		return fmt.Sprintf("%s must be at most 500 characters long", field)
	case "sort_order":
		return fmt.Sprintf("%s must be between 0 and 999999", field)
	case "sku":
		return fmt.Sprintf("%s must be alphanumeric with hyphens and underscores", field)
	case "price":
		return fmt.Sprintf("%s must be non-negative", field)
	case "weight":
		return fmt.Sprintf("%s must be non-negative", field)
	case "dimensions":
		return fmt.Sprintf("%s must be in the format 'length x width x height'", field)
	case "tags":
		return fmt.Sprintf("%s must be a comma-separated list of tags", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func ValidateSlug(slug string) error {
	if strings.TrimSpace(slug) == "" {
		return fmt.Errorf("slug is required")
	}

	if len(slug) < 1 {
		return fmt.Errorf("slug must be at least 1 character long")
	}

	if len(slug) > 255 {
		return fmt.Errorf("slug must be at most 255 characters long")
	}

	// Slug should only contain lowercase letters, numbers, and hyphens
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(slug) {
		return fmt.Errorf("slug should only contain lowercase letters, numbers, and hyphens")
	}

	// Slug should not start or end with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return fmt.Errorf("slug should not start or end with hyphen")
	}

	// Slug should not contain consecutive hyphens
	if strings.Contains(slug, "--") {
		return fmt.Errorf("slug should not contain consecutive hyphens")
	}

	return nil
}

// Product validation functions

// validateSKU validates SKU format (alphanumeric, hyphens, underscores allowed)
func validateSKU(fl validator.FieldLevel) bool {
	sku := fl.Field().String()
	if sku == "" {
		return false
	}

	// SKU should be alphanumeric with hyphens and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, sku)
	return matched
}

// validatePrice validates price (must be non-negative)
func validatePrice(fl validator.FieldLevel) bool {
	price := fl.Field().Float()
	return price >= 0
}

// validateWeight validates weight (must be non-negative)
func validateWeight(fl validator.FieldLevel) bool {
	weight := fl.Field().Float()
	return weight >= 0
}

// validateDimensions validates dimensions format (e.g., "10x20x30", "10 x 20 x 30")
func validateDimensions(fl validator.FieldLevel) bool {
	dimensions := fl.Field().String()
	if dimensions == "" {
		return true // Optional field
	}

	// Allow formats like "10x20x30", "10 x 20 x 30", "10cm x 20cm x 30cm"
	matched, _ := regexp.MatchString(`^[\d\s\.xX×\*]+(cm|in|mm|m)?$`, dimensions)
	return matched
}

// validateTags validates tags format (comma-separated, no special characters)
func validateTags(fl validator.FieldLevel) bool {
	tags := fl.Field().String()
	if tags == "" {
		return true // Optional field
	}

	// Tags should be alphanumeric with spaces, commas, and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\s,-]+$`, tags)
	return matched
}

// ValidateProductFilters validates product filter parameters
func ValidateProductFilters(filters map[string]interface{}) ValidatorErrors {
	var errors ValidatorErrors

	// Validate sort_by parameter
	if sortBy, ok := filters["sort_by"].(string); ok && sortBy != "" {
		validSortFields := []string{"name", "price", "created_at", "updated_at", "sku"}
		if !contains(validSortFields, sortBy) {
			errors = append(errors, ValidatorError{
				Field:   "sort_by",
				Tag:     "invalid",
				Value:   sortBy,
				Message: "sort_by must be one of: name, price, created_at, updated_at, sku",
			})
		}
	}

	// Validate sort_order parameter
	if sortOrder, ok := filters["sort_order"].(string); ok && sortOrder != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			errors = append(errors, ValidatorError{
				Field:   "sort_order",
				Tag:     "invalid",
				Value:   sortOrder,
				Message: "sort_order must be either 'asc' or 'desc'",
			})
		}
	}

	// Validate price range
	if minPrice, ok := filters["min_price"].(float64); ok && minPrice < 0 {
		errors = append(errors, ValidatorError{
			Field:   "min_price",
			Tag:     "invalid",
			Value:   fmt.Sprintf("%.2f", minPrice),
			Message: "min_price must be non-negative",
		})
	}

	if maxPrice, ok := filters["max_price"].(float64); ok && maxPrice < 0 {
		errors = append(errors, ValidatorError{
			Field:   "max_price",
			Tag:     "invalid",
			Value:   fmt.Sprintf("%.2f", maxPrice),
			Message: "max_price must be non-negative",
		})
	}

	// Validate price range logic
	if minPrice, ok := filters["min_price"].(float64); ok {
		if maxPrice, ok := filters["max_price"].(float64); ok && maxPrice < minPrice {
			errors = append(errors, ValidatorError{
				Field:   "price_range",
				Tag:     "invalid",
				Value:   fmt.Sprintf("%.2f - %.2f", minPrice, maxPrice),
				Message: "max_price must be greater than or equal to min_price",
			})
		}
	}

	return errors
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Cart validation functions

// validateCurrency validates currency code (ISO 4217)
func validateCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	if currency == "" {
		return false
	}

	// Common currency codes
	validCurrencies := []string{
		"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD",
		"INR", "BRL", "RUB", "KRW", "SGD", "HKD", "NOK", "MXN", "TRY", "ZAR",
	}

	return contains(validCurrencies, strings.ToUpper(currency))
}

// validateSessionID validates session ID format
func validateSessionID(fl validator.FieldLevel) bool {
	sessionID := fl.Field().String()
	if sessionID == "" {
		return false
	}

	// Session ID should be alphanumeric with hyphens and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, sessionID)
	return matched
}

// validateCouponCode validates coupon code format
func validateCouponCode(fl validator.FieldLevel) bool {
	couponCode := fl.Field().String()
	if couponCode == "" {
		return false
	}

	// Coupon code should be alphanumeric with hyphens and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, couponCode)
	return matched
}

// validateShippingMethod validates shipping method name
func validateShippingMethod(fl validator.FieldLevel) bool {
	shippingMethod := fl.Field().String()
	if shippingMethod == "" {
		return false
	}

	// Shipping method should contain only letters, numbers, spaces, and common punctuation
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\s\-\.]+$`, shippingMethod)
	return matched
}

// ValidateCartFilters validates cart filter parameters
func ValidateCartFilters(filters map[string]interface{}) ValidatorErrors {
	var errors ValidatorErrors

	// Validate currency parameter
	if currency, ok := filters["currency"].(string); ok && currency != "" {
		// Simple currency validation without using validator.FieldLevel
		validCurrencies := []string{
			"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "SEK", "NZD",
			"INR", "BRL", "RUB", "KRW", "SGD", "HKD", "NOK", "MXN", "TRY", "ZAR",
		}
		if !contains(validCurrencies, strings.ToUpper(currency)) {
			errors = append(errors, ValidatorError{
				Field:   "currency",
				Tag:     "invalid",
				Value:   currency,
				Message: "currency must be a valid ISO 4217 currency code",
			})
		}
	}

	// Validate user_id parameter
	if userID, ok := filters["user_id"].(int64); ok && userID <= 0 {
		errors = append(errors, ValidatorError{
			Field:   "user_id",
			Tag:     "invalid",
			Value:   fmt.Sprintf("%d", userID),
			Message: "user_id must be a positive integer",
		})
	}

	return errors
}