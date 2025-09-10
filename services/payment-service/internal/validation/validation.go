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
	_ = validate.RegisterValidation("currency", validateCurrency)
	_ = validate.RegisterValidation("payment_method_type", validatePaymentMethodType)
	_ = validate.RegisterValidation("payment_status", validatePaymentStatus)
	_ = validate.RegisterValidation("gateway_code", validateGatewayCode)
	_ = validate.RegisterValidation("transaction_id", validateTransactionID)
	_ = validate.RegisterValidation("refund_id", validateRefundID)
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

// Currency validation
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

// Payment method type validation
func validatePaymentMethodType(fl validator.FieldLevel) bool {
	paymentType := fl.Field().String()
	if paymentType == "" {
		return false
	}

	validTypes := []string{
		"credit_card", "debit_card", "paypal", "bank_transfer", "digital_wallet",
	}

	return contains(validTypes, paymentType)
}

// Payment status validation
func validatePaymentStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	if status == "" {
		return false
	}

	validStatuses := []string{
		"pending", "processing", "completed", "failed", "cancelled", "refunded",
	}

	return contains(validStatuses, status)
}

// Gateway code validation
func validateGatewayCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if code == "" {
		return false
	}

	validCodes := []string{
		"stripe", "paypal", "razorpay",
	}

	return contains(validCodes, code)
}

// Transaction ID validation
func validateTransactionID(fl validator.FieldLevel) bool {
	transactionID := fl.Field().String()
	if transactionID == "" {
		return false
	}

	// Transaction ID should be alphanumeric with underscores and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, transactionID)
	return matched
}

// Refund ID validation
func validateRefundID(fl validator.FieldLevel) bool {
	refundID := fl.Field().String()
	if refundID == "" {
		return false
	}

	// Refund ID should be alphanumeric with underscores and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, refundID)
	return matched
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

func getErrorMessage(field, tag string, value interface{}) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %v", field, value)
	case "max":
		return fmt.Sprintf("%s must be at most %v", field, value)
	case "len":
		return fmt.Sprintf("%s must be exactly %v characters long", field, value)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %v", field, value)
	case "currency":
		return fmt.Sprintf("%s must be a valid ISO 4217 currency code", field)
	case "payment_method_type":
		return fmt.Sprintf("%s must be one of: credit_card, debit_card, paypal, bank_transfer, digital_wallet", field)
	case "payment_status":
		return fmt.Sprintf("%s must be one of: pending, processing, completed, failed, cancelled, refunded", field)
	case "gateway_code":
		return fmt.Sprintf("%s must be one of: stripe, paypal, razorpay", field)
	case "transaction_id":
		return fmt.Sprintf("%s must be alphanumeric with underscores and hyphens", field)
	case "refund_id":
		return fmt.Sprintf("%s must be alphanumeric with underscores and hyphens", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// ValidatePaymentFilters validates payment filter parameters
func ValidatePaymentFilters(filters map[string]interface{}) ValidatorErrors {
	var errors ValidatorErrors

	// Validate status parameter
	if status, ok := filters["status"].(string); ok && status != "" {
		validStatuses := []string{"pending", "processing", "completed", "failed", "cancelled", "refunded"}
		if !contains(validStatuses, status) {
			errors = append(errors, ValidatorError{
				Field:   "status",
				Tag:     "invalid",
				Value:   status,
				Message: "status must be one of: pending, processing, completed, failed, cancelled, refunded",
			})
		}
	}

	// Validate currency parameter
	if currency, ok := filters["currency"].(string); ok && currency != "" {
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

	// Validate gateway_id parameter
	if gatewayID, ok := filters["gateway_id"].(string); ok && gatewayID != "" {
		validGateways := []string{"stripe", "paypal", "razorpay"}
		if !contains(validGateways, gatewayID) {
			errors = append(errors, ValidatorError{
				Field:   "gateway_id",
				Tag:     "invalid",
				Value:   gatewayID,
				Message: "gateway_id must be one of: stripe, paypal, razorpay",
			})
		}
	}

	// Validate sort_by parameter
	if sortBy, ok := filters["sort_by"].(string); ok && sortBy != "" {
		validSortFields := []string{"id", "order_id", "amount", "created_at", "updated_at", "status"}
		if !contains(validSortFields, sortBy) {
			errors = append(errors, ValidatorError{
				Field:   "sort_by",
				Tag:     "invalid",
				Value:   sortBy,
				Message: "sort_by must be one of: id, order_id, amount, created_at, updated_at, status",
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

	// Validate amount range
	if minAmount, ok := filters["min_amount"].(float64); ok && minAmount < 0 {
		errors = append(errors, ValidatorError{
			Field:   "min_amount",
			Tag:     "invalid",
			Value:   fmt.Sprintf("%.2f", minAmount),
			Message: "min_amount must be non-negative",
		})
	}

	if maxAmount, ok := filters["max_amount"].(float64); ok && maxAmount < 0 {
		errors = append(errors, ValidatorError{
			Field:   "max_amount",
			Tag:     "invalid",
			Value:   fmt.Sprintf("%.2f", maxAmount),
			Message: "max_amount must be non-negative",
		})
	}

	// Validate amount range logic
	if minAmount, ok := filters["min_amount"].(float64); ok {
		if maxAmount, ok := filters["max_amount"].(float64); ok && maxAmount < minAmount {
			errors = append(errors, ValidatorError{
				Field:   "amount_range",
				Tag:     "invalid",
				Value:   fmt.Sprintf("%.2f - %.2f", minAmount, maxAmount),
				Message: "max_amount must be greater than or equal to min_amount",
			})
		}
	}

	return errors
}
