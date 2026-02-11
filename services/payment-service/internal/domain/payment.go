package domain

import (
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID              int64                  `json:"id" db:"id"`
	OrderID         int64                  `json:"order_id" db:"order_id"`
	PaymentMethod   string                 `json:"payment_method" db:"payment_method"`
	TransactionID   string                 `json:"transaction_id" db:"transaction_id"`
	GatewayID       string                 `json:"gateway_id" db:"gateway_id"`
	Amount          float64                `json:"amount" db:"amount"`
	Currency        string                 `json:"currency" db:"currency"`
	Status          string                 `json:"status" db:"status"` // pending, processing, completed, failed, cancelled, refunded
	GatewayStatus   string                 `json:"gateway_status" db:"gateway_status"`
	GatewayResponse string                 `json:"gateway_response" db:"gateway_response"`
	FailureReason   string                 `json:"failure_reason" db:"failure_reason"`
	ProcessedAt     *time.Time             `json:"processed_at" db:"processed_at"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata" db:"metadata"`
}

// PaymentRefund represents refund transactions
type PaymentRefund struct {
	ID              int64      `json:"id" db:"id"`
	PaymentID       int64      `json:"payment_id" db:"payment_id"`
	RefundID        string     `json:"refund_id" db:"refund_id"`
	Amount          float64    `json:"amount" db:"amount"`
	Reason          string     `json:"reason" db:"reason"`
	Status          string     `json:"status" db:"status"` // pending, processed, failed
	GatewayResponse string     `json:"gateway_response" db:"gateway_response"`
	ProcessedAt     *time.Time `json:"processed_at" db:"processed_at"`
	CreatedBy       int64      `json:"created_by" db:"created_by"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// PaymentWebhook represents webhook events from payment gateways
type PaymentWebhook struct {
	ID          int64      `json:"id" db:"id"`
	GatewayID   string     `json:"gateway_id" db:"gateway_id"`
	EventType   string     `json:"event_type" db:"event_type"`
	EventID     string     `json:"event_id" db:"event_id"`
	Payload     string     `json:"payload" db:"payload"` // JSON payload
	Signature   string     `json:"signature" db:"signature"`
	IsProcessed bool       `json:"is_processed" db:"is_processed"`
	ProcessedAt *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// PaymentSummary represents payment analytics
type PaymentSummary struct {
	TotalPayments      int64   `json:"total_payments" db:"total_payments"`
	SuccessfulPayments int64   `json:"successful_payments" db:"successful_payments"`
	FailedPayments     int64   `json:"failed_payments" db:"failed_payments"`
	PendingPayments    int64   `json:"pending_payments" db:"pending_payments"`
	TotalAmount        float64 `json:"total_amount" db:"total_amount"`
	RefundedAmount     float64 `json:"refunded_amount" db:"refunded_amount"`
	NetAmount          float64 `json:"net_amount" db:"net_amount"`
	AverageAmount      float64 `json:"average_amount" db:"average_amount"`
	SuccessRate        float64 `json:"success_rate" db:"success_rate"`
}

// PaymentFilter represents filters for payment queries
type PaymentFilter struct {
	OrderID       *int64     `json:"order_id"`
	PaymentMethod *string    `json:"payment_method"`
	Status        *string    `json:"status"`
	GatewayID     *string    `json:"gateway_id"`
	Currency      *string    `json:"currency"`
	DateFrom      *time.Time `json:"date_from"`
	DateTo        *time.Time `json:"date_to"`
	Search        string     `json:"search"`
	Page          int        `json:"page"`
	Limit         int        `json:"limit"`
	SortBy        string     `json:"sort_by"`
	SortOrder     string     `json:"sort_order"`
}

// Payment Status Constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusCompleted  = "completed"
	PaymentStatusFailed     = "failed"
	PaymentStatusCancelled  = "cancelled"
	PaymentStatusRefunded   = "refunded"
)

// Payment Gateway Codes
const (
	GatewayStripe = "stripe"
)

// Refund Status Constants
const (
	RefundStatusPending   = "pending"
	RefundStatusProcessed = "processed"
	RefundStatusFailed    = "failed"
)
