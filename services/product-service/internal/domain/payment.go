package domain

import (
	"time"
)

// Payment represents a payment transaction
type Payment struct {
	ID              int64      `json:"id" db:"id"`
	OrderID         int64      `json:"order_id" db:"order_id"`
	PaymentMethodID int64      `json:"payment_method_id" db:"payment_method_id"`
	TransactionID   string     `json:"transaction_id" db:"transaction_id"`
	GatewayID       string     `json:"gateway_id" db:"gateway_id"`
	Amount          float64    `json:"amount" db:"amount"`
	Currency        string     `json:"currency" db:"currency"`
	Status          string     `json:"status" db:"status"` // pending, processing, completed, failed, cancelled, refunded
	GatewayStatus   string     `json:"gateway_status" db:"gateway_status"`
	GatewayResponse string     `json:"gateway_response" db:"gateway_response"`
	FailureReason   string     `json:"failure_reason" db:"failure_reason"`
	ProcessedAt     *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// PaymentMethod represents available payment methods
type PaymentMethod struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Code        string `json:"code" db:"code"`
	Type        string `json:"type" db:"type"` // credit_card, debit_card, paypal, bank_transfer, etc.
	IsActive    bool   `json:"is_active" db:"is_active"`
	IsDefault   bool   `json:"is_default" db:"is_default"`
	SortOrder   int    `json:"sort_order" db:"sort_order"`
	Description string `json:"description" db:"description"`
	Icon        string `json:"icon" db:"icon"`
}

// PaymentGateway represents payment gateway configurations
type PaymentGateway struct {
	ID         int64  `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Code       string `json:"code" db:"code"`
	IsActive   bool   `json:"is_active" db:"is_active"`
	IsTestMode bool   `json:"is_test_mode" db:"is_test_mode"`
	Config     string `json:"config" db:"config"` // JSON configuration
	WebhookURL string `json:"webhook_url" db:"webhook_url"`
	SortOrder  int    `json:"sort_order" db:"sort_order"`
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

// Coupon represents discount coupons
type Coupon struct {
	ID                int64      `json:"id" db:"id"`
	Code              string     `json:"code" db:"code"`
	Name              string     `json:"name" db:"name"`
	Description       string     `json:"description" db:"description"`
	Type              string     `json:"type" db:"type"` // percentage, fixed_amount, free_shipping
	Value             float64    `json:"value" db:"value"`
	MinOrderAmount    float64    `json:"min_order_amount" db:"min_order_amount"`
	MaxDiscountAmount float64    `json:"max_discount_amount" db:"max_discount_amount"`
	UsageLimit        int        `json:"usage_limit" db:"usage_limit"`
	UsedCount         int        `json:"used_count" db:"used_count"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	StartsAt          time.Time  `json:"starts_at" db:"starts_at"`
	ExpiresAt         *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// CouponUsage represents coupon usage tracking
type CouponUsage struct {
	ID        int64     `json:"id" db:"id"`
	CouponID  int64     `json:"coupon_id" db:"coupon_id"`
	OrderID   int64     `json:"order_id" db:"order_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
