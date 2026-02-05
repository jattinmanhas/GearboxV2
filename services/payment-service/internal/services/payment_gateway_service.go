package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/repository"
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

// PaymentGatewayService handles payment gateway integrations (Stripe only)
type PaymentGatewayService struct {
	paymentRepo   *repository.PaymentRepository
	stripeGateway *StripeGateway
	config        *config.Config
}

// NewPaymentGatewayService creates a new payment gateway service
func NewPaymentGatewayService(paymentRepo *repository.PaymentRepository, cfg *config.Config) *PaymentGatewayService {
	return &PaymentGatewayService{
		paymentRepo:   paymentRepo,
		stripeGateway: NewStripeGateway(cfg.PaymentGateways.Stripe),
		config:        cfg,
	}
}

// ProcessPayment processes a payment through Stripe
func (s *PaymentGatewayService) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error) {
	if payment.GatewayID != domain.GatewayStripe {
		return nil, fmt.Errorf("unsupported payment gateway: %s (only stripe is supported)", payment.GatewayID)
	}
	return s.stripeGateway.ProcessPayment(ctx, payment, paymentData)
}

// ProcessRefund processes a refund through Stripe
func (s *PaymentGatewayService) ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error) {
	if payment.GatewayID != domain.GatewayStripe {
		return nil, fmt.Errorf("unsupported payment gateway: %s (only stripe is supported)", payment.GatewayID)
	}
	return s.stripeGateway.ProcessRefund(ctx, payment, amount, reason)
}

// ValidateWebhook validates a webhook from Stripe
func (s *PaymentGatewayService) ValidateWebhook(ctx context.Context, gatewayID string, payload []byte, signature string) (bool, error) {
	if gatewayID != domain.GatewayStripe {
		return false, fmt.Errorf("unsupported payment gateway: %s (only stripe is supported)", gatewayID)
	}
	return s.stripeGateway.ValidateWebhook(ctx, payload, signature)
}

// ProcessWebhook processes a webhook from Stripe
func (s *PaymentGatewayService) ProcessWebhook(ctx context.Context, gatewayID string, payload []byte) (*GatewayResponse, error) {
	if gatewayID != domain.GatewayStripe {
		return nil, fmt.Errorf("unsupported payment gateway: %s (only stripe is supported)", gatewayID)
	}
	return s.stripeGateway.ProcessWebhook(ctx, payload)
}

// StripeGateway handles Stripe payment processing
type StripeGateway struct {
	secretKey      string
	publishableKey string
	webhookSecret  string
}

// NewStripeGateway creates a new Stripe gateway instance
func NewStripeGateway(cfg config.StripeConfig) *StripeGateway {
	if cfg.SecretKey != "" {
		stripe.Key = cfg.SecretKey
	}

	return &StripeGateway{
		secretKey:      cfg.SecretKey,
		publishableKey: cfg.PublishableKey,
		webhookSecret:  cfg.WebhookSecret,
	}
}

// ProcessPayment processes a payment through Stripe
func (g *StripeGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentData map[string]interface{}) (*GatewayResponse, error) {
	if g.secretKey == "" {
		return nil, fmt.Errorf("stripe secret key not configured")
	}

	paymentMethodID, _ := paymentData["payment_method_id"].(string)

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)), // Convert to cents
		Currency: stripe.String(strings.ToLower(payment.Currency)),
		Metadata: map[string]string{
			"order_id":       fmt.Sprintf("%d", payment.OrderID),
			"transaction_id": payment.TransactionID,
		},
	}

	if paymentMethodID != "" {
		params.PaymentMethod = stripe.String(paymentMethodID)
		params.ConfirmationMethod = stripe.String(string(stripe.PaymentIntentConfirmationMethodManual))
		params.Confirm = stripe.Bool(true)
	}

	if email, ok := paymentData["customer_email"].(string); ok && email != "" {
		params.ReceiptEmail = stripe.String(email)
	}

	if description, ok := paymentData["description"].(string); ok && description != "" {
		params.Description = stripe.String(description)
	}

	pi, err := paymentintent.New(params)
	if err != nil {
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

	responseJSON, _ := json.Marshal(pi)

	return &GatewayResponse{
		Success:       pi.Status == stripe.PaymentIntentStatusSucceeded,
		Status:        string(pi.Status),
		TransactionID: pi.ID,
		Response:      string(responseJSON),
	}, nil
}

// ProcessRefund processes a refund through Stripe
func (g *StripeGateway) ProcessRefund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*GatewayResponse, error) {
	if g.secretKey == "" {
		return nil, fmt.Errorf("stripe secret key not configured")
	}

	paymentIntentID := payment.TransactionID
	if payment.GatewayResponse != "" {
		var piData map[string]interface{}
		if err := json.Unmarshal([]byte(payment.GatewayResponse), &piData); err == nil {
			if id, ok := piData["id"].(string); ok {
				paymentIntentID = id
			}
		}
	}

	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
		Amount:        stripe.Int64(int64(amount * 100)),
		Reason:        stripe.String(reason),
	}

	ref, err := refund.New(params)
	if err != nil {
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

	responseJSON, _ := json.Marshal(ref)

	return &GatewayResponse{
		Success:       ref.Status == stripe.RefundStatusSucceeded || ref.Status == stripe.RefundStatusPending,
		Status:        string(ref.Status),
		TransactionID: ref.ID,
		Response:      string(responseJSON),
	}, nil
}

// ValidateWebhook validates a Stripe webhook signature
func (g *StripeGateway) ValidateWebhook(ctx context.Context, payload []byte, signature string) (bool, error) {
	if g.webhookSecret == "" {
		return false, fmt.Errorf("stripe webhook secret not configured")
	}

	parts := strings.Split(signature, ",")
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid signature format")
	}

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

	signedPayload := fmt.Sprintf("%s.%s", timestamp, string(payload))
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write([]byte(signedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	for _, sig := range signatures {
		if hmac.Equal([]byte(sig), []byte(expectedSignature)) {
			return true, nil
		}
	}

	return false, nil
}

// ProcessWebhook processes a Stripe webhook event
func (g *StripeGateway) ProcessWebhook(ctx context.Context, payload []byte) (*GatewayResponse, error) {
	var event stripe.Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

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
