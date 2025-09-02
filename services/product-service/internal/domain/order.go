package domain

import (
	"time"
)

// Order represents a customer order
type Order struct {
	ID                int64      `json:"id" db:"id"`
	OrderNumber       string     `json:"order_number" db:"order_number"`
	UserID            int64      `json:"user_id" db:"user_id"`
	Status            string     `json:"status" db:"status"`                         // pending, confirmed, processing, shipped, delivered, cancelled, refunded
	PaymentStatus     string     `json:"payment_status" db:"payment_status"`         // pending, paid, failed, refunded, partially_refunded
	FulfillmentStatus string     `json:"fulfillment_status" db:"fulfillment_status"` // unfulfilled, partial, fulfilled
	Subtotal          float64    `json:"subtotal" db:"subtotal"`
	TaxAmount         float64    `json:"tax_amount" db:"tax_amount"`
	ShippingAmount    float64    `json:"shipping_amount" db:"shipping_amount"`
	DiscountAmount    float64    `json:"discount_amount" db:"discount_amount"`
	TotalAmount       float64    `json:"total_amount" db:"total_amount"`
	Currency          string     `json:"currency" db:"currency"`
	Notes             string     `json:"notes" db:"notes"`
	InternalNotes     string     `json:"internal_notes" db:"internal_notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	ConfirmedAt       *time.Time `json:"confirmed_at" db:"confirmed_at"`
	ShippedAt         *time.Time `json:"shipped_at" db:"shipped_at"`
	DeliveredAt       *time.Time `json:"delivered_at" db:"delivered_at"`
	CancelledAt       *time.Time `json:"cancelled_at" db:"cancelled_at"`
}

// OrderItem represents individual items in an order
type OrderItem struct {
	ID               int64   `json:"id" db:"id"`
	OrderID          int64   `json:"order_id" db:"order_id"`
	ProductID        int64   `json:"product_id" db:"product_id"`
	ProductVariantID *int64  `json:"product_variant_id" db:"product_variant_id"`
	ProductName      string  `json:"product_name" db:"product_name"`
	ProductSKU       string  `json:"product_sku" db:"product_sku"`
	Quantity         int     `json:"quantity" db:"quantity"`
	UnitPrice        float64 `json:"unit_price" db:"unit_price"`
	TotalPrice       float64 `json:"total_price" db:"total_price"`
	TaxAmount        float64 `json:"tax_amount" db:"tax_amount"`
	DiscountAmount   float64 `json:"discount_amount" db:"discount_amount"`
	IsDigital        bool    `json:"is_digital" db:"is_digital"`
	RequiresShipping bool    `json:"requires_shipping" db:"requires_shipping"`
}

// OrderAddress represents shipping and billing addresses for an order
type OrderAddress struct {
	ID         int64  `json:"id" db:"id"`
	OrderID    int64  `json:"order_id" db:"order_id"`
	Type       string `json:"type" db:"type"` // shipping, billing
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
	Company    string `json:"company" db:"company"`
	Address1   string `json:"address1" db:"address1"`
	Address2   string `json:"address2" db:"address2"`
	City       string `json:"city" db:"city"`
	State      string `json:"state" db:"state"`
	Country    string `json:"country" db:"country"`
	PostalCode string `json:"postal_code" db:"postal_code"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
}

// OrderStatusHistory represents the history of order status changes
type OrderStatusHistory struct {
	ID             int64     `json:"id" db:"id"`
	OrderID        int64     `json:"order_id" db:"order_id"`
	Status         string    `json:"status" db:"status"`
	PreviousStatus string    `json:"previous_status" db:"previous_status"`
	Notes          string    `json:"notes" db:"notes"`
	CreatedBy      int64     `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// OrderFulfillment represents order fulfillment information
type OrderFulfillment struct {
	ID                int64      `json:"id" db:"id"`
	OrderID           int64      `json:"order_id" db:"order_id"`
	TrackingNumber    string     `json:"tracking_number" db:"tracking_number"`
	Carrier           string     `json:"carrier" db:"carrier"`
	Service           string     `json:"service" db:"service"`
	Status            string     `json:"status" db:"status"` // pending, shipped, delivered, failed
	ShippedAt         *time.Time `json:"shipped_at" db:"shipped_at"`
	DeliveredAt       *time.Time `json:"delivered_at" db:"delivered_at"`
	EstimatedDelivery *time.Time `json:"estimated_delivery" db:"estimated_delivery"`
	Notes             string     `json:"notes" db:"notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// OrderRefund represents refund information for orders
type OrderRefund struct {
	ID          int64      `json:"id" db:"id"`
	OrderID     int64      `json:"order_id" db:"order_id"`
	Amount      float64    `json:"amount" db:"amount"`
	Reason      string     `json:"reason" db:"reason"`
	Status      string     `json:"status" db:"status"` // pending, processed, failed
	ProcessedAt *time.Time `json:"processed_at" db:"processed_at"`
	CreatedBy   int64      `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// OrderFilter represents filters for order queries
type OrderFilter struct {
	UserID            *int64     `json:"user_id"`
	Status            string     `json:"status"`
	PaymentStatus     string     `json:"payment_status"`
	FulfillmentStatus string     `json:"fulfillment_status"`
	DateFrom          *time.Time `json:"date_from"`
	DateTo            *time.Time `json:"date_to"`
	MinAmount         *float64   `json:"min_amount"`
	MaxAmount         *float64   `json:"max_amount"`
	Search            string     `json:"search"`
}
