package dto

import "time"

// CreateOrderRequest represents the request to create a new order
type CreateOrderRequest struct {
	UserID                int64                    `json:"user_id" validate:"required"`
	CartID                *int64                   `json:"cart_id"` // Optional: create from cart
	Items                 []CreateOrderItemRequest `json:"items" validate:"required,min=1"`
	ShippingAddress       OrderAddressRequest      `json:"shipping_address"`         // Optional if using user address
	BillingAddress        *OrderAddressRequest     `json:"billing_address"`          // Optional: defaults to shipping
	UserShippingAddressID *uint                    `json:"user_shipping_address_id"` // Optional: use saved user address
	UserBillingAddressID  *uint                    `json:"user_billing_address_id"`  // Optional: use saved user address
	PaymentMethodID       int64                    `json:"payment_method_id" validate:"required"`
	Currency              string                   `json:"currency" validate:"required,len=3"`
	Notes                 string                   `json:"notes"`
	InternalNotes         string                   `json:"internal_notes"`
	ApplyCoupons          []string                 `json:"apply_coupons"` // Coupon codes to apply
}

// CreateOrderItemRequest represents an item in order creation
type CreateOrderItemRequest struct {
	ProductID        int64   `json:"product_id" validate:"required"`
	ProductVariantID *int64  `json:"product_variant_id"`
	Quantity         int     `json:"quantity" validate:"required,min=1"`
	UnitPrice        float64 `json:"unit_price" validate:"required,min=0"`
}

// OrderAddressRequest represents address information
type OrderAddressRequest struct {
	FirstName  string `json:"first_name" validate:"required,min=1,max=100"`
	LastName   string `json:"last_name" validate:"required,min=1,max=100"`
	Company    string `json:"company" validate:"max=100"`
	Address1   string `json:"address1" validate:"required,min=1,max=255"`
	Address2   string `json:"address2" validate:"max=255"`
	City       string `json:"city" validate:"required,min=1,max=100"`
	State      string `json:"state" validate:"required,min=1,max=100"`
	Country    string `json:"country" validate:"required,min=1,max=100"`
	PostalCode string `json:"postal_code" validate:"required,min=1,max=20"`
	Phone      string `json:"phone" validate:"max=20"`
	Email      string `json:"email" validate:"required,email"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status            *string `json:"status" validate:"omitempty,oneof=pending confirmed processing shipped delivered cancelled refunded"`
	PaymentStatus     *string `json:"payment_status" validate:"omitempty,oneof=pending paid failed refunded partially_refunded"`
	FulfillmentStatus *string `json:"fulfillment_status" validate:"omitempty,oneof=unfulfilled partial fulfilled"`
	Notes             *string `json:"notes"`
	InternalNotes     *string `json:"internal_notes"`
}

// OrderResponse represents the response for order data
type OrderResponse struct {
	ID                int64                        `json:"id"`
	OrderNumber       string                       `json:"order_number"`
	UserID            int64                        `json:"user_id"`
	Status            string                       `json:"status"`
	PaymentStatus     string                       `json:"payment_status"`
	FulfillmentStatus string                       `json:"fulfillment_status"`
	Subtotal          float64                      `json:"subtotal"`
	TaxAmount         float64                      `json:"tax_amount"`
	ShippingAmount    float64                      `json:"shipping_amount"`
	DiscountAmount    float64                      `json:"discount_amount"`
	TotalAmount       float64                      `json:"total_amount"`
	Currency          string                       `json:"currency"`
	Notes             string                       `json:"notes"`
	InternalNotes     string                       `json:"internal_notes"`
	Items             []OrderItemResponse          `json:"items"`
	Addresses         []OrderAddressResponse       `json:"addresses"`
	StatusHistory     []OrderStatusHistoryResponse `json:"status_history"`
	Fulfillment       *OrderFulfillmentResponse    `json:"fulfillment"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
	ConfirmedAt       *time.Time                   `json:"confirmed_at"`
	ShippedAt         *time.Time                   `json:"shipped_at"`
	DeliveredAt       *time.Time                   `json:"delivered_at"`
	CancelledAt       *time.Time                   `json:"cancelled_at"`
}

// OrderItemResponse represents the response for order item data
type OrderItemResponse struct {
	ID               int64   `json:"id"`
	OrderID          int64   `json:"order_id"`
	ProductID        int64   `json:"product_id"`
	ProductVariantID *int64  `json:"product_variant_id"`
	ProductName      string  `json:"product_name"`
	ProductSKU       string  `json:"product_sku"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
	TotalPrice       float64 `json:"total_price"`
	TaxAmount        float64 `json:"tax_amount"`
	DiscountAmount   float64 `json:"discount_amount"`
	IsDigital        bool    `json:"is_digital"`
	RequiresShipping bool    `json:"requires_shipping"`
}

// OrderAddressResponse represents the response for order address data
type OrderAddressResponse struct {
	ID         int64  `json:"id"`
	OrderID    int64  `json:"order_id"`
	Type       string `json:"type"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Company    string `json:"company"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

// OrderStatusHistoryResponse represents the response for order status history
type OrderStatusHistoryResponse struct {
	ID             int64     `json:"id"`
	OrderID        int64     `json:"order_id"`
	Status         string    `json:"status"`
	PreviousStatus string    `json:"previous_status"`
	Notes          string    `json:"notes"`
	CreatedBy      int64     `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

// OrderFulfillmentResponse represents the response for order fulfillment
type OrderFulfillmentResponse struct {
	ID                int64      `json:"id"`
	OrderID           int64      `json:"order_id"`
	TrackingNumber    string     `json:"tracking_number"`
	Carrier           string     `json:"carrier"`
	Service           string     `json:"service"`
	Status            string     `json:"status"`
	ShippedAt         *time.Time `json:"shipped_at"`
	DeliveredAt       *time.Time `json:"delivered_at"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// OrderRefundResponse represents the response for order refund
type OrderRefundResponse struct {
	ID          int64      `json:"id"`
	OrderID     int64      `json:"order_id"`
	Amount      float64    `json:"amount"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedBy   int64      `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ListOrdersRequest represents the request to list orders with filters
type ListOrdersRequest struct {
	UserID            *int64     `json:"user_id"`
	Status            string     `json:"status"`
	PaymentStatus     string     `json:"payment_status"`
	FulfillmentStatus string     `json:"fulfillment_status"`
	DateFrom          *time.Time `json:"date_from"`
	DateTo            *time.Time `json:"date_to"`
	MinAmount         *float64   `json:"min_amount"`
	MaxAmount         *float64   `json:"max_amount"`
	Search            string     `json:"search"`
	Page              int        `json:"page"`
	Limit             int        `json:"limit"`
}

// ListOrdersResponse represents the response for listing orders
type ListOrdersResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// UpdateOrderStatusRequest represents the request to update order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed processing shipped delivered cancelled refunded"`
	Notes  string `json:"notes"`
}

// CreateOrderFulfillmentRequest represents the request to create order fulfillment
type CreateOrderFulfillmentRequest struct {
	TrackingNumber    string     `json:"tracking_number" validate:"required"`
	Carrier           string     `json:"carrier" validate:"required"`
	Service           string     `json:"service"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	Notes             string     `json:"notes"`
}

// UpdateOrderFulfillmentRequest represents the request to update order fulfillment
type UpdateOrderFulfillmentRequest struct {
	TrackingNumber    *string    `json:"tracking_number"`
	Carrier           *string    `json:"carrier"`
	Service           *string    `json:"service"`
	Status            *string    `json:"status" validate:"omitempty,oneof=pending shipped delivered failed"`
	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	Notes             *string    `json:"notes"`
}

// CreateOrderRefundRequest represents the request to create order refund
type CreateOrderRefundRequest struct {
	Amount    float64 `json:"amount" validate:"required,min=0.01"`
	Reason    string  `json:"reason" validate:"required,min=1"`
	CreatedBy int64   `json:"created_by" validate:"required"`
}

// OrderAnalyticsResponse represents analytics data for orders
type OrderAnalyticsResponse struct {
	TotalOrders        int64   `json:"total_orders"`
	PendingOrders      int64   `json:"pending_orders"`
	ConfirmedOrders    int64   `json:"confirmed_orders"`
	ProcessingOrders   int64   `json:"processing_orders"`
	ShippedOrders      int64   `json:"shipped_orders"`
	DeliveredOrders    int64   `json:"delivered_orders"`
	CancelledOrders    int64   `json:"cancelled_orders"`
	TotalRevenue       float64 `json:"total_revenue"`
	AverageOrderValue  float64 `json:"average_order_value"`
	ConversionRate     float64 `json:"conversion_rate"`
	NewOrdersToday     int64   `json:"new_orders_today"`
	NewOrdersThisWeek  int64   `json:"new_orders_this_week"`
	NewOrdersThisMonth int64   `json:"new_orders_this_month"`
}

// OrderSummaryResponse represents a summary of order data
type OrderSummaryResponse struct {
	ID          int64     `json:"id"`
	OrderNumber string    `json:"order_number"`
	UserID      int64     `json:"user_id"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `json:"currency"`
	ItemCount   int       `json:"item_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrderFromCartRequest represents the request to create an order from a cart
type CreateOrderFromCartRequest struct {
	ShippingAddress       OrderAddressRequest  `json:"shipping_address"` // Optional if using user address
	BillingAddress        *OrderAddressRequest `json:"billing_address"`
	UserShippingAddressID *uint                `json:"user_shipping_address_id"` // Optional: use saved user address
	UserBillingAddressID  *uint                `json:"user_billing_address_id"`  // Optional: use saved user address
	PaymentMethodID       int64                `json:"payment_method_id" validate:"required"`
	Currency              string               `json:"currency" validate:"required,len=3"`
	Notes                 string               `json:"notes"`
	ApplyCoupons          []string             `json:"apply_coupons"`
}

// ProductOrderStatsResponse represents the response for product order statistics
type ProductOrderStatsResponse struct {
	ProductID     int64   `json:"product_id"`
	ProductName   string  `json:"product_name"`
	SKU           string  `json:"sku"`
	TotalQuantity int64   `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"`
	OrderCount    int64   `json:"order_count"`
	AveragePrice  float64 `json:"average_price"`
}
