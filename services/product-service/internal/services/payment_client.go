package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PaymentServiceClient handles communication with the payment service
type PaymentServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// PaymentServiceRequest represents a request to create a payment
type PaymentServiceRequest struct {
	OrderID       int64                  `json:"order_id"`
	PaymentMethod string                 `json:"payment_method"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	GatewayID     string                 `json:"gateway_id"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentServiceResponse represents a response from payment service
type PaymentServiceResponse struct {
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
	ProcessedAt     *time.Time             `json:"processed_at"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessPaymentRequest represents a request to process a payment
type ProcessPaymentRequest struct {
	PaymentID   int64                  `json:"payment_id"`
	PaymentData map[string]interface{} `json:"payment_data"`
	ReturnURL   string                 `json:"return_url,omitempty"`
	CancelURL   string                 `json:"cancel_url,omitempty"`
}

// UpdatePaymentStatusRequest represents a request to update payment status
type UpdatePaymentStatusRequest struct {
	Status          string `json:"status"`
	GatewayStatus   string `json:"gateway_status,omitempty"`
	GatewayResponse string `json:"gateway_response,omitempty"`
	FailureReason   string `json:"failure_reason,omitempty"`
}

// NewPaymentServiceClient creates a new payment service client
func NewPaymentServiceClient(baseURL string) *PaymentServiceClient {
	return &PaymentServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreatePayment creates a new payment in the payment service
func (c *PaymentServiceClient) CreatePayment(ctx context.Context, req *PaymentServiceRequest, authToken string) (*PaymentServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/protected/payments", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, fmt.Errorf("payment service error: %v", errorResp)
	}

	var response struct {
		Data PaymentServiceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// ProcessPayment processes a payment through the gateway
func (c *PaymentServiceClient) ProcessPayment(ctx context.Context, req *ProcessPaymentRequest, authToken string) (*PaymentServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/protected/payments/process", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, fmt.Errorf("payment service error: %v", errorResp)
	}

	var response struct {
		Data PaymentServiceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// GetPayment retrieves a payment by ID
func (c *PaymentServiceClient) GetPayment(ctx context.Context, paymentID int64, authToken string) (*PaymentServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/protected/payments/%d", c.baseURL, paymentID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, fmt.Errorf("payment service error: %v", errorResp)
	}

	var response struct {
		Data PaymentServiceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// UpdatePaymentStatus updates payment status
func (c *PaymentServiceClient) UpdatePaymentStatus(ctx context.Context, paymentID int64, req *UpdatePaymentStatusRequest, authToken string) (*PaymentServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/protected/payments/%d/status", c.baseURL, paymentID)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, fmt.Errorf("payment service error: %v", errorResp)
	}

	var response struct {
		Data PaymentServiceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (c *PaymentServiceClient) GetPaymentByTransactionID(ctx context.Context, transactionID, authToken string) (*PaymentServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/protected/payments/transaction/%s", c.baseURL, transactionID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, fmt.Errorf("payment service error: %v", errorResp)
	}

	var response struct {
		Data PaymentServiceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}

// GetOrderPayment retrieves payment information for an order
func (c *PaymentServiceClient) GetOrderPayment(ctx context.Context, orderID int64, authToken string) (*PaymentServiceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/protected/payments/order/%d", c.baseURL, orderID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No payment found for this order
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, fmt.Errorf("payment service error: %v", errorResp)
	}

	var response struct {
		Data PaymentServiceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Data, nil
}
