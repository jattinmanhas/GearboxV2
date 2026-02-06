package dto

// CreatePaymentRequest represents a request to create a payment for an order
type CreatePaymentRequest struct {
	PaymentMethod string                 `json:"payment_method" validate:"required"`
	GatewayID     string                 `json:"gateway_id" validate:"required"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessPaymentRequest represents a request to process a payment
type ProcessPaymentRequest struct {
	PaymentData map[string]interface{} `json:"payment_data" validate:"required"`
	ReturnURL   string                 `json:"return_url,omitempty"`
	CancelURL   string                 `json:"cancel_url,omitempty"`
}

// PaymentResponse represents a payment response from the payment service
type PaymentResponse struct {
	ID              int64                  `json:"id"`
	OrderID         int64                  `json:"order_id"`
	PaymentMethod   string                 `json:"payment_method"`
	TransactionID   string                 `json:"transaction_id"`
	GatewayID       string                 `json:"gateway_id"`
	Amount          float64                `json:"amount"`
	Currency        string                 `json:"currency"`
	Status          string                 `json:"status"`
	GatewayStatus   string                 `json:"gateway_status"`
	GatewayResponse string                 `json:"gateway_response"`
	FailureReason   string                 `json:"failure_reason"`
	ProcessedAt     *string                `json:"processed_at"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}
