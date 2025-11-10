package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/repository"
	razorpay "github.com/razorpay/razorpay-go"
	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
)

// GatewayResponse represents a response from a payment gateway
type GatewayResponse struct {
	Success       bool   `json:"success"`
	Status        string `json:"status"`
	Response      string `json:"response"`
	ErrorMessage  string `json:"error_message,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
	RedirectURL   string `json:"redirect_url,omitempty"`
}

// PaymentGatewayService handles payment gateway integrations
type PaymentGatewayService struct {
	paymentRepo *repository.PaymentRepository
	gateways    map[string]GatewayProvider
	config      *config.Config
}

// GatewayProvider interface for different payment gateways
type GatewayProvider interface {
	ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error)
	ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error)
	ValidateWebhook(ctx context.Context, payload []byte, signature string) (bool, error)
	ProcessWebhook(ctx context.Context, payload []byte) (*GatewayResponse, error)
}

// NewPaymentGatewayService creates a new payment gateway service
func NewPaymentGatewayService(paymentRepo *repository.PaymentRepository, cfg *config.Config) *PaymentGatewayService {
	service := &PaymentGatewayService{
		paymentRepo: paymentRepo,
		gateways:    make(map[string]GatewayProvider),
		config:      cfg,
	}

	// Initialize gateway providers with configuration
	service.gateways[domain.GatewayStripe] = NewStripeGateway(cfg.PaymentGateways.Stripe)
	service.gateways[domain.GatewayPayPal] = NewPayPalGateway(cfg.PaymentGateways.PayPal)
	service.gateways[domain.GatewayRazorpay] = NewRazorpayGateway(cfg.PaymentGateways.Razorpay)

	return service
}

// ProcessPayment processes a payment through the specified gateway
func (s *PaymentGatewayService) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error) {
	gateway, exists := s.gateways[payment.GatewayID]
	if !exists {
		return nil, fmt.Errorf("unsupported payment gateway: %s", payment.GatewayID)
	}

	return gateway.ProcessPayment(ctx, payment, paymentData)
}

// ProcessRefund processes a refund through the specified gateway
func (s *PaymentGatewayService) ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error) {
	gateway, exists := s.gateways[payment.GatewayID]
	if !exists {
		return nil, fmt.Errorf("unsupported payment gateway: %s", payment.GatewayID)
	}

	return gateway.ProcessRefund(ctx, payment, amount, reason)
}

// ValidateWebhook validates a webhook from a payment gateway
func (s *PaymentGatewayService) ValidateWebhook(ctx context.Context, gatewayID string, payload []byte, signature string) (bool, error) {
	gateway, exists := s.gateways[gatewayID]
	if !exists {
		return false, fmt.Errorf("unsupported payment gateway: %s", gatewayID)
	}

	return gateway.ValidateWebhook(ctx, payload, signature)
}

// ProcessWebhook processes a webhook from a payment gateway
func (s *PaymentGatewayService) ProcessWebhook(ctx context.Context, gatewayID string, payload []byte) (*GatewayResponse, error) {
	gateway, exists := s.gateways[gatewayID]
	if !exists {
		return nil, fmt.Errorf("unsupported payment gateway: %s", gatewayID)
	}

	return gateway.ProcessWebhook(ctx, payload)
}

// Stripe Gateway Implementation

type StripeGateway struct {
	secretKey      string
	publishableKey string
	webhookSecret  string
}

func NewStripeGateway(cfg config.StripeConfig) *StripeGateway {
	// Set Stripe API key
	if cfg.SecretKey != "" {
		stripe.Key = cfg.SecretKey
	}

	return &StripeGateway{
		secretKey:      cfg.SecretKey,
		publishableKey: cfg.PublishableKey,
		webhookSecret:  cfg.WebhookSecret,
	}
}

func (g *StripeGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error) {
	if g.secretKey == "" {
		return nil, fmt.Errorf("stripe secret key not configured")
	}

	// Extract payment method ID from paymentData if available
	paymentMethodID, _ := paymentData["payment_method_id"].(string)

	// Create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)), // Convert to cents
		Currency: stripe.String(strings.ToLower(payment.Currency)),
		Metadata: map[string]string{
			"order_id":       fmt.Sprintf("%d", payment.OrderID),
			"transaction_id": payment.TransactionID,
		},
	}

	// Add payment method if provided
	if paymentMethodID != "" {
		params.PaymentMethod = stripe.String(paymentMethodID)
		params.ConfirmationMethod = stripe.String(string(stripe.PaymentIntentConfirmationMethodManual))
		params.Confirm = stripe.Bool(true)
	}

	// Add customer email if available
	if email, ok := paymentData["customer_email"].(string); ok && email != "" {
		params.ReceiptEmail = stripe.String(email)
	}

	// Add description
	if description, ok := paymentData["description"].(string); ok && description != "" {
		params.Description = stripe.String(description)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		// Check if it's a Stripe error
		if stripeErr, ok := err.(*stripe.Error); ok {
			return &GatewayResponse{
				Success:      false,
				Status:       string(pi.Status),
				ErrorMessage: stripeErr.Msg,
				Response:     fmt.Sprintf(`{"error": {"type": "%s", "code": "%s", "message": "%s"}}`, stripeErr.Type, stripeErr.Code, stripeErr.Msg),
			}, nil
		}
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Convert response to JSON
	responseJSON, _ := json.Marshal(pi)

	return &GatewayResponse{
		Success:       pi.Status == stripe.PaymentIntentStatusSucceeded,
		Status:        string(pi.Status),
		TransactionID: pi.ID,
		Response:      string(responseJSON),
	}, nil
}

func (g *StripeGateway) ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error) {
	if g.secretKey == "" {
		return nil, fmt.Errorf("stripe secret key not configured")
	}

	// Get the payment intent ID from gateway response or transaction ID
	paymentIntentID := payment.TransactionID
	if payment.GatewayResponse != "" {
		var piData map[string]interface{}
		if err := json.Unmarshal([]byte(payment.GatewayResponse), &piData); err == nil {
			if id, ok := piData["id"].(string); ok {
				paymentIntentID = id
			}
		}
	}

	// Create refund
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
		Amount:        stripe.Int64(int64(amount * 100)), // Convert to cents
		Reason:        stripe.String(reason),
	}

	ref, err := refund.New(params)
	if err != nil {
		// Check if it's a Stripe error
		if stripeErr, ok := err.(*stripe.Error); ok {
			return &GatewayResponse{
				Success:      false,
				Status:       "failed",
				ErrorMessage: stripeErr.Msg,
				Response:     fmt.Sprintf(`{"error": {"type": "%s", "code": "%s", "message": "%s"}}`, stripeErr.Type, stripeErr.Code, stripeErr.Msg),
			}, nil
		}
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// Convert response to JSON
	responseJSON, _ := json.Marshal(ref)

	return &GatewayResponse{
		Success:       ref.Status == stripe.RefundStatusSucceeded || ref.Status == stripe.RefundStatusPending,
		Status:        string(ref.Status),
		TransactionID: ref.ID,
		Response:      string(responseJSON),
	}, nil
}

func (g *StripeGateway) ValidateWebhook(ctx context.Context, payload []byte, signature string) (bool, error) {
	if g.webhookSecret == "" {
		return false, fmt.Errorf("stripe webhook secret not configured")
	}

	// Stripe webhook signature validation
	// Format: timestamp,signature1,signature2
	parts := strings.Split(signature, ",")
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid signature format")
	}

	// Extract timestamp and signatures
	var timestamp string
	var signatures []string
	for _, part := range parts {
		if strings.HasPrefix(part, "t=") {
			timestamp = strings.TrimPrefix(part, "t=")
		} else if strings.HasPrefix(part, "v1=") {
			signatures = append(signatures, strings.TrimPrefix(part, "v1="))
		}
	}

	if timestamp == "" || len(signatures) == 0 {
		return false, fmt.Errorf("invalid signature format")
	}

	// Create signed payload
	signedPayload := fmt.Sprintf("%s.%s", timestamp, string(payload))

	// Compute HMAC
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	for _, sig := range signatures {
		if hmac.Equal([]byte(sig), []byte(expectedSignature)) {
			return true, nil
		}
	}

	return false, nil
}

func (g *StripeGateway) ProcessWebhook(ctx context.Context, payload []byte) (*GatewayResponse, error) {
	var event stripe.Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Handle different Stripe events
	switch event.Type {
	case "payment_intent.succeeded":
		return &GatewayResponse{
			Success:  true,
			Status:   "succeeded",
			Response: string(payload),
		}, nil
	case "payment_intent.payment_failed":
		return &GatewayResponse{
			Success:  false,
			Status:   "failed",
			Response: string(payload),
		}, nil
	case "charge.refunded":
		return &GatewayResponse{
			Success:  true,
			Status:   "refunded",
			Response: string(payload),
		}, nil
	default:
		return &GatewayResponse{
			Success:  true,
			Status:   "processed",
			Response: string(payload),
		}, nil
	}
}

// PayPal Gateway Implementation

type PayPalGateway struct {
	clientID     string
	clientSecret string
	mode         string // sandbox or live
	baseURL      string
	accessToken  string
}

func NewPayPalGateway(cfg config.PayPalConfig) *PayPalGateway {
	baseURL := "https://api.sandbox.paypal.com"
	if cfg.Mode == "live" {
		baseURL = "https://api.paypal.com"
	}

	return &PayPalGateway{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		mode:         cfg.Mode,
		baseURL:      baseURL,
	}
}

// getAccessToken retrieves PayPal OAuth access token
func (g *PayPalGateway) getAccessToken(ctx context.Context) (string, error) {
	if g.accessToken != "" {
		return g.accessToken, nil
	}

	if g.clientID == "" || g.clientSecret == "" {
		return "", fmt.Errorf("paypal credentials not configured")
	}

	url := fmt.Sprintf("%s/v1/oauth2/token", g.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(g.clientID, g.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get access token: %s", string(body))
	}

	g.accessToken = tokenResp.AccessToken
	return g.accessToken, nil
}

func (g *PayPalGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error) {
	accessToken, err := g.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Create PayPal order
	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": payment.TransactionID,
				"amount": map[string]interface{}{
					"currency_code": payment.Currency,
					"value":         fmt.Sprintf("%.2f", payment.Amount),
				},
			},
		},
	}

	// Add description if available
	if description, ok := paymentData["description"].(string); ok && description != "" {
		orderData["purchase_units"].([]map[string]interface{})[0]["description"] = description
	}

	jsonData, _ := json.Marshal(orderData)

	url := fmt.Sprintf("%s/v2/checkout/orders", g.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.Unmarshal(body, &errorResp)
		return &GatewayResponse{
			Success:      false,
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("PayPal API error: %v", errorResp),
			Response:     string(body),
		}, nil
	}

	var orderResp map[string]interface{}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	// Extract approval URL for redirect
	var redirectURL string
	if links, ok := orderResp["links"].([]interface{}); ok {
		for _, link := range links {
			if linkMap, ok := link.(map[string]interface{}); ok {
				if rel, ok := linkMap["rel"].(string); ok && rel == "approve" {
					if href, ok := linkMap["href"].(string); ok {
						redirectURL = href
					}
				}
			}
		}
	}

	orderID, _ := orderResp["id"].(string)
	status, _ := orderResp["status"].(string)

	return &GatewayResponse{
		Success:       status == "CREATED" || status == "APPROVED",
		Status:        status,
		TransactionID: orderID,
		RedirectURL:   redirectURL,
		Response:      string(body),
	}, nil
}

func (g *PayPalGateway) ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error) {
	accessToken, err := g.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Get capture ID from payment gateway response
	var captureID string
	if payment.GatewayResponse != "" {
		var paymentData map[string]interface{}
		if err := json.Unmarshal([]byte(payment.GatewayResponse), &paymentData); err == nil {
			if purchaseUnits, ok := paymentData["purchase_units"].([]interface{}); ok && len(purchaseUnits) > 0 {
				if unit, ok := purchaseUnits[0].(map[string]interface{}); ok {
					if captures, ok := unit["payments"].(map[string]interface{})["captures"].([]interface{}); ok && len(captures) > 0 {
						if capture, ok := captures[0].(map[string]interface{}); ok {
							if id, ok := capture["id"].(string); ok {
								captureID = id
							}
						}
					}
				}
			}
		}
	}

	if captureID == "" {
		return nil, fmt.Errorf("capture ID not found in payment data")
	}

	// Create refund
	refundData := map[string]interface{}{
		"amount": map[string]interface{}{
			"value":         fmt.Sprintf("%.2f", amount),
			"currency_code": payment.Currency,
		},
		"note_to_payer": reason,
	}

	jsonData, _ := json.Marshal(refundData)

	url := fmt.Sprintf("%s/v2/payments/captures/%s/refund", g.baseURL, captureID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		json.Unmarshal(body, &errorResp)
		return &GatewayResponse{
			Success:      false,
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("PayPal API error: %v", errorResp),
			Response:     string(body),
		}, nil
	}

	var refundResp map[string]interface{}
	if err := json.Unmarshal(body, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to parse refund response: %w", err)
	}

	refundID, _ := refundResp["id"].(string)
	status, _ := refundResp["status"].(string)

	return &GatewayResponse{
		Success:       status == "COMPLETED",
		Status:        status,
		TransactionID: refundID,
		Response:      string(body),
	}, nil
}

func (g *PayPalGateway) ValidateWebhook(ctx context.Context, payload []byte, signature string) (bool, error) {
	// PayPal webhook validation requires verifying the signature
	// This is a simplified version - in production, you should use PayPal's webhook verification API
	if g.clientID == "" || g.clientSecret == "" {
		return false, fmt.Errorf("paypal credentials not configured")
	}

	// For now, we'll do basic validation
	// In production, you should call PayPal's webhook verification endpoint
	return signature != "", nil
}

func (g *PayPalGateway) ProcessWebhook(ctx context.Context, payload []byte) (*GatewayResponse, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	eventType, ok := event["event_type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid webhook event type")
	}

	// Handle different PayPal events
	switch eventType {
	case "PAYMENT.SALE.COMPLETED", "PAYMENT.CAPTURE.COMPLETED":
		return &GatewayResponse{
			Success:  true,
			Status:   "completed",
			Response: string(payload),
		}, nil
	case "PAYMENT.SALE.DENIED", "PAYMENT.CAPTURE.DENIED":
		return &GatewayResponse{
			Success:  false,
			Status:   "denied",
			Response: string(payload),
		}, nil
	case "PAYMENT.CAPTURE.REFUNDED":
		return &GatewayResponse{
			Success:  true,
			Status:   "refunded",
			Response: string(payload),
		}, nil
	default:
		return &GatewayResponse{
			Success:  true,
			Status:   "processed",
			Response: string(payload),
		}, nil
	}
}

// Razorpay Gateway Implementation

type RazorpayGateway struct {
	keyID         string
	keySecret     string
	webhookSecret string
	client        *razorpay.Client
}

func NewRazorpayGateway(cfg config.RazorpayConfig) *RazorpayGateway {
	var razorpayClient *razorpay.Client
	if cfg.KeyID != "" && cfg.KeySecret != "" {
		razorpayClient = razorpay.NewClient(cfg.KeyID, cfg.KeySecret)
	}

	return &RazorpayGateway{
		keyID:         cfg.KeyID,
		keySecret:     cfg.KeySecret,
		webhookSecret: cfg.WebhookSecret,
		client:        razorpayClient,
	}
}

func (g *RazorpayGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error) {
	if g.keyID == "" || g.keySecret == "" {
		return nil, fmt.Errorf("razorpay credentials not configured")
	}

	if g.client == nil {
		return nil, fmt.Errorf("razorpay client not initialized")
	}

	// Convert amount to smallest currency unit (paise for INR)
	amount := int64(payment.Amount * 100)

	// Create payment data
	paymentDataMap := map[string]interface{}{
		"amount":   amount,
		"currency": payment.Currency,
		"receipt":  payment.TransactionID,
		"notes": map[string]interface{}{
			"order_id":       payment.OrderID,
			"transaction_id": payment.TransactionID,
		},
	}

	// Add customer details if available
	if email, ok := paymentData["customer_email"].(string); ok && email != "" {
		paymentDataMap["email"] = email
	}
	if phone, ok := paymentData["customer_phone"].(string); ok && phone != "" {
		paymentDataMap["contact"] = phone
	}
	if name, ok := paymentData["customer_name"].(string); ok && name != "" {
		paymentDataMap["name"] = name
	}

	// For Razorpay, we create an order first
	// The actual payment capture happens after user completes payment on Razorpay's page
	orderData := map[string]interface{}{
		"amount":   amount,
		"currency": payment.Currency,
		"receipt":  payment.TransactionID,
		"notes": map[string]interface{}{
			"order_id":       payment.OrderID,
			"transaction_id": payment.TransactionID,
		},
	}

	razorpayOrder, err := g.client.Order.Create(orderData, nil)
	if err != nil {
		return &GatewayResponse{
			Success:      false,
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, nil
	}

	// For Razorpay, we return the order which will be used to initiate payment
	// The actual payment happens on Razorpay's checkout page
	razorpayPaymentResp := razorpayOrder

	// Convert response to JSON
	responseJSON, _ := json.Marshal(razorpayPaymentResp)

	status := "created"
	if razorpayPaymentResp["status"] != nil {
		status = fmt.Sprintf("%v", razorpayPaymentResp["status"])
	}

	paymentID := ""
	if razorpayPaymentResp["id"] != nil {
		paymentID = fmt.Sprintf("%v", razorpayPaymentResp["id"])
	}

	return &GatewayResponse{
		Success:       status == "authorized" || status == "captured",
		Status:        status,
		TransactionID: paymentID,
		Response:      string(responseJSON),
	}, nil
}

func (g *RazorpayGateway) ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error) {
	if g.keyID == "" || g.keySecret == "" {
		return nil, fmt.Errorf("razorpay credentials not configured")
	}

	if g.client == nil {
		return nil, fmt.Errorf("razorpay client not initialized")
	}

	// Get payment ID from transaction ID or gateway response
	paymentID := payment.TransactionID
	if payment.GatewayResponse != "" {
		var paymentData map[string]interface{}
		if err := json.Unmarshal([]byte(payment.GatewayResponse), &paymentData); err == nil {
			if id, ok := paymentData["id"].(string); ok {
				paymentID = id
			}
		}
	}

	// Convert amount to smallest currency unit
	refundAmount := int64(amount * 100)

	// Create refund data
	refundData := map[string]interface{}{
		"amount": refundAmount,
		"notes": map[string]interface{}{
			"reason": reason,
		},
	}

	razorpayRefund, err := g.client.Refund.Create(refundData, nil)
	if err != nil {
		// Try with payment_id in the refund data
		refundData["payment_id"] = paymentID
		razorpayRefund, err = g.client.Refund.Create(refundData, nil)
	}
	if err != nil {
		return &GatewayResponse{
			Success:      false,
			Status:       "failed",
			ErrorMessage: err.Error(),
		}, nil
	}

	// Convert response to JSON
	responseJSON, _ := json.Marshal(razorpayRefund)

	status := "processed"
	if razorpayRefund["status"] != nil {
		status = fmt.Sprintf("%v", razorpayRefund["status"])
	}

	refundID := ""
	if razorpayRefund["id"] != nil {
		refundID = fmt.Sprintf("%v", razorpayRefund["id"])
	}

	return &GatewayResponse{
		Success:       status == "processed" || status == "pending",
		Status:        status,
		TransactionID: refundID,
		Response:      string(responseJSON),
	}, nil
}

func (g *RazorpayGateway) ValidateWebhook(ctx context.Context, payload []byte, signature string) (bool, error) {
	if g.webhookSecret == "" {
		return false, fmt.Errorf("razorpay webhook secret not configured")
	}

	// Razorpay webhook signature validation
	// Signature is HMAC SHA256 of payload with webhook secret
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

func (g *RazorpayGateway) ProcessWebhook(ctx context.Context, payload []byte) (*GatewayResponse, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	eventType, ok := event["event"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid webhook event type")
	}

	// Handle different Razorpay events
	switch eventType {
	case "payment.captured", "payment.authorized":
		return &GatewayResponse{
			Success:  true,
			Status:   "captured",
			Response: string(payload),
		}, nil
	case "payment.failed":
		return &GatewayResponse{
			Success:  false,
			Status:   "failed",
			Response: string(payload),
		}, nil
	case "refund.created", "refund.processed":
		return &GatewayResponse{
			Success:  true,
			Status:   "refunded",
			Response: string(payload),
		}, nil
	default:
		return &GatewayResponse{
			Success:  true,
			Status:   "processed",
			Response: string(payload),
		}, nil
	}
}
