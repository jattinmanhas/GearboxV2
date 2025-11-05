package dto

import "time"

// Payment Request DTOs

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	OrderID         int64                  `json:"order_id" validate:"required,min=1"`
	PaymentMethodID int64                  `json:"payment_method_id" validate:"required,min=1"`
	Amount          float64                `json:"amount" validate:"required,min=0.01"`
	Currency        string                 `json:"currency" validate:"required,len=3"`
	GatewayID       string                 `json:"gateway_id" validate:"required,min=1,max=50"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessPaymentRequest represents a request to process a payment
type ProcessPaymentRequest struct {
	PaymentID   int64                  `json:"payment_id" validate:"required,min=1"`
	PaymentData map[string]interface{} `json:"payment_data" validate:"required"`
	ReturnURL   string                 `json:"return_url,omitempty"`
	CancelURL   string                 `json:"cancel_url,omitempty"`
}

// RefundPaymentRequest represents a request to refund a payment
type RefundPaymentRequest struct {
	PaymentID int64   `json:"payment_id" validate:"required,min=1"`
	Amount    float64 `json:"amount" validate:"required,min=0.01"`
	Reason    string  `json:"reason" validate:"required,min=1,max=500"`
}

// UpdatePaymentStatusRequest represents a request to update payment status
type UpdatePaymentStatusRequest struct {
	Status          string `json:"status" validate:"required,oneof=pending processing completed failed cancelled refunded"`
	GatewayStatus   string `json:"gateway_status,omitempty"`
	GatewayResponse string `json:"gateway_response,omitempty"`
	FailureReason   string `json:"failure_reason,omitempty"`
}

// PaymentMethodRequest represents a request to create/update payment method
type PaymentMethodRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Code        string `json:"code" validate:"required,min=1,max=50"`
	Type        string `json:"type" validate:"required,oneof=credit_card debit_card paypal bank_transfer digital_wallet"`
	IsActive    bool   `json:"is_active"`
	IsDefault   bool   `json:"is_default"`
	SortOrder   int    `json:"sort_order" validate:"min=0,max=999"`
	Description string `json:"description" validate:"max=500"`
	Icon        string `json:"icon" validate:"max=255"`
}

// PaymentGatewayRequest represents a request to create/update payment gateway
type PaymentGatewayRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	Code       string `json:"code" validate:"required,min=1,max=50"`
	IsActive   bool   `json:"is_active"`
	IsTestMode bool   `json:"is_test_mode"`
	Config     string `json:"config" validate:"required"`
	WebhookURL string `json:"webhook_url" validate:"max=500"`
	SortOrder  int    `json:"sort_order" validate:"min=0,max=999"`
}

// Payment Response DTOs

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID              int64                  `json:"id"`
	OrderID         int64                  `json:"order_id"`
	PaymentMethodID int64                  `json:"payment_method_id"`
	TransactionID   string                 `json:"transaction_id"`
	GatewayID       string                 `json:"gateway_id"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	Status          string                 `json:"status"`
	GatewayStatus   string                 `json:"gateway_status"`
	GatewayResponse string                 `json:"gateway_response"`
	FailureReason   string                 `json:"failure_reason"`
	ProcessedAt     *time.Time             `json:"processed_at"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentMethodResponse represents a payment method response
type PaymentMethodResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Type        string `json:"type"`
	IsActive    bool   `json:"is_active"`
	IsDefault   bool   `json:"is_default"`
	SortOrder   int    `json:"sort_order"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// PaymentGatewayResponse represents a payment gateway response
type PaymentGatewayResponse struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	IsActive   bool    `json:"is_active"`
	IsTestMode bool    `json:"is_test_mode"`
	Config     string  `json:"config"`
	WebhookURL *string `json:"webhook_url,omitempty"` // Nullable webhook URL
	SortOrder  int     `json:"sort_order"`
}

// PaymentRefundResponse represents a refund response
type PaymentRefundResponse struct {
	ID              int64      `json:"id"`
	PaymentID       int64      `json:"payment_id"`
	RefundID        string     `json:"refund_id"`
	Amount          float64    `json:"amount"`
	Reason          string     `json:"reason"`
	Status          string     `json:"status"`
	GatewayResponse string     `json:"gateway_response"`
	ProcessedAt     *time.Time `json:"processed_at"`
	CreatedBy       int64      `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
}

// PaymentWebhookResponse represents a webhook response
type PaymentWebhookResponse struct {
	ID          int64      `json:"id"`
	GatewayID   string     `json:"gateway_id"`
	EventType   string     `json:"event_type"`
	EventID     string     `json:"event_id"`
	Payload     string     `json:"payload"`
	Signature   string     `json:"signature"`
	IsProcessed bool       `json:"is_processed"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Payment Summary Response
type PaymentSummaryResponse struct {
	TotalPayments      int64   `json:"total_payments"`
	SuccessfulPayments int64   `json:"successful_payments"`
	FailedPayments     int64   `json:"failed_payments"`
	PendingPayments    int64   `json:"pending_payments"`
	TotalAmount        float64 `json:"total_amount"`
	RefundedAmount     float64 `json:"refunded_amount"`
	NetAmount          float64 `json:"net_amount"`
}

// Payment Filter Request
type PaymentFilterRequest struct {
	OrderID         *int64  `json:"order_id"`
	PaymentMethodID *int64  `json:"payment_method_id"`
	Status          *string `json:"status"`
	GatewayID       *string `json:"gateway_id"`
	Currency        *string `json:"currency"`
	DateFrom        *string `json:"date_from"`
	DateTo          *string `json:"date_to"`
	Search          string  `json:"search"`
	Page            int     `json:"page" validate:"min=1"`
	Limit           int     `json:"limit" validate:"min=1,max=100"`
	SortBy          string  `json:"sort_by" validate:"oneof=id order_id amount created_at updated_at"`
	SortOrder       string  `json:"sort_order" validate:"oneof=asc desc"`
}

// Payment List Response
type PaymentListResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Pages    int               `json:"pages"`
}

// Gateway-specific DTOs

// StripePaymentRequest represents Stripe-specific payment data
type StripePaymentRequest struct {
	PaymentMethodID    string `json:"payment_method_id" validate:"required"`
	CustomerID         string `json:"customer_id,omitempty"`
	ConfirmationMethod string `json:"confirmation_method,omitempty"`
	ReturnURL          string `json:"return_url,omitempty"`
}

// PayPalPaymentRequest represents PayPal-specific payment data
type PayPalPaymentRequest struct {
	ApprovalURL string `json:"approval_url,omitempty"`
	PaymentID   string `json:"payment_id,omitempty"`
	PayerID     string `json:"payer_id,omitempty"`
}

// Webhook Event DTOs

// WebhookEventRequest represents an incoming webhook event
type WebhookEventRequest struct {
	GatewayID string                 `json:"gateway_id" validate:"required"`
	EventType string                 `json:"event_type" validate:"required"`
	EventID   string                 `json:"event_id" validate:"required"`
	Payload   map[string]interface{} `json:"payload" validate:"required"`
	Signature string                 `json:"signature" validate:"required"`
}

// WebhookEventResponse represents a webhook event response
type WebhookEventResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
