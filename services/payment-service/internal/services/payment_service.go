package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/repository"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo    *repository.PaymentRepository
	gatewayService *PaymentGatewayService
}

// NewPaymentService creates a new payment service
func NewPaymentService(paymentRepo *repository.PaymentRepository, gatewayService *PaymentGatewayService) *PaymentService {
	return &PaymentService{
		paymentRepo:    paymentRepo,
		gatewayService: gatewayService,
	}
}

// Payments

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Validate gateway (currently only Stripe is supported)
	if req.GatewayID != "stripe" {
		return nil, fmt.Errorf("unsupported payment gateway: %s", req.GatewayID)
	}

	// Generate unique transaction ID
	transactionID, err := s.generateTransactionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ID: %w", err)
	}

	payment := &domain.Payment{
		OrderID:       req.OrderID,
		PaymentMethod: req.PaymentMethod,
		TransactionID: transactionID,
		GatewayID:     req.GatewayID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        domain.PaymentStatusPending,
		Metadata:      req.Metadata,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return s.paymentToResponse(payment), nil
}

// ProcessPayment processes a payment through the gateway
func (s *PaymentService) ProcessPayment(ctx context.Context, req *dto.ProcessPaymentRequest) (*dto.PaymentResponse, error) {
	// Get payment
	payment, err := s.paymentRepo.GetPaymentByID(ctx, req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	// Check if payment can be processed
	if payment.Status != domain.PaymentStatusPending {
		return nil, fmt.Errorf("payment cannot be processed: current status is %s", payment.Status)
	}

	// Update status to processing
	payment.Status = domain.PaymentStatusProcessing
	payment.UpdatedAt = time.Now()
	if err := s.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// Process payment through gateway
	gatewayResponse, err := s.gatewayService.ProcessPayment(ctx, payment, req.PaymentData)
	if err != nil {
		// Update payment status to failed
		payment.Status = domain.PaymentStatusFailed
		payment.FailureReason = err.Error()
		payment.UpdatedAt = time.Now()
		s.paymentRepo.UpdatePayment(ctx, payment)
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Update payment with gateway response
	payment.GatewayStatus = gatewayResponse.Status
	payment.GatewayResponse = gatewayResponse.Response
	payment.UpdatedAt = time.Now()

	if gatewayResponse.Success {
		payment.Status = domain.PaymentStatusCompleted
		now := time.Now()
		payment.ProcessedAt = &now
	} else {
		payment.Status = domain.PaymentStatusFailed
		payment.FailureReason = gatewayResponse.ErrorMessage
	}

	if err := s.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return s.paymentToResponse(payment), nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, id int64) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	return s.paymentToResponse(payment), nil
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (s *PaymentService) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetPaymentByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	return s.paymentToResponse(payment), nil
}

// GetPaymentByOrderID retrieves the latest payment by order ID
func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID int64) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetLatestPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	return s.paymentToResponse(payment), nil
}

// ListPayments retrieves payments with filtering and pagination
func (s *PaymentService) ListPayments(ctx context.Context, filter *domain.PaymentFilter, page, limit int) (*dto.PaymentListResponse, error) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	payments, err := s.paymentRepo.ListPayments(ctx, filter, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	total, err := s.paymentRepo.CountPayments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count payments: %w", err)
	}

	var responses []dto.PaymentResponse
	for _, payment := range payments {
		responses = append(responses, *s.paymentToResponse(payment))
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))
	if pages == 0 {
		pages = 1
	}

	return &dto.PaymentListResponse{
		Payments: responses,
		Total:    total,
		Page:     page,
		Limit:    limit,
		Pages:    pages,
	}, nil
}

// UpdatePaymentStatus updates payment status
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, id int64, req *dto.UpdatePaymentStatusRequest) (*dto.PaymentResponse, error) {
	payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	var processedAt *time.Time
	if req.Status == domain.PaymentStatusCompleted {
		now := time.Now()
		processedAt = &now
	}

	if err := s.paymentRepo.UpdatePaymentStatus(ctx, id, req.Status, req.GatewayStatus, req.GatewayResponse, req.FailureReason, processedAt); err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// Get updated payment
	payment, err = s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated payment: %w", err)
	}

	return s.paymentToResponse(payment), nil
}

// RefundPayment refunds a payment
func (s *PaymentService) RefundPayment(ctx context.Context, req *dto.RefundPaymentRequest) (*dto.PaymentRefundResponse, error) {
	// Get payment
	payment, err := s.paymentRepo.GetPaymentByID(ctx, req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	// Check if payment can be refunded
	if payment.Status != domain.PaymentStatusCompleted {
		return nil, fmt.Errorf("payment cannot be refunded: current status is %s", payment.Status)
	}

	// Check if refund amount is valid
	if req.Amount <= 0 || req.Amount > payment.Amount {
		return nil, fmt.Errorf("invalid refund amount")
	}

	// Generate unique refund ID
	refundID, err := s.generateRefundID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refund ID: %w", err)
	}

	// Process refund through gateway
	gatewayResponse, err := s.gatewayService.ProcessRefund(ctx, payment, req.Amount, req.Reason)
	if err != nil {
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}

	// Create refund record
	refund := &domain.PaymentRefund{
		PaymentID:       payment.ID,
		RefundID:        refundID,
		Amount:          req.Amount,
		Reason:          req.Reason,
		Status:          domain.RefundStatusProcessed,
		GatewayResponse: gatewayResponse.Response,
		ProcessedAt:     &time.Time{},
		CreatedBy:       0, // TODO: Get from context
		CreatedAt:       time.Now(),
	}

	if gatewayResponse.Success {
		now := time.Now()
		refund.ProcessedAt = &now
	} else {
		refund.Status = domain.RefundStatusFailed
	}

	if err := s.paymentRepo.CreatePaymentRefund(ctx, refund); err != nil {
		return nil, fmt.Errorf("failed to create payment refund: %w", err)
	}

	// Update payment status if fully refunded
	if req.Amount >= payment.Amount {
		payment.Status = domain.PaymentStatusRefunded
		payment.UpdatedAt = time.Now()
		s.paymentRepo.UpdatePayment(ctx, payment)
	}

	return s.paymentRefundToResponse(refund), nil
}

// GetPaymentSummary retrieves payment summary statistics
func (s *PaymentService) GetPaymentSummary(ctx context.Context, dateFrom, dateTo *time.Time) (*dto.PaymentSummaryResponse, error) {
	summary, err := s.paymentRepo.GetPaymentSummary(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	return &dto.PaymentSummaryResponse{
		TotalPayments:      summary.TotalPayments,
		SuccessfulPayments: summary.SuccessfulPayments,
		FailedPayments:     summary.FailedPayments,
		PendingPayments:    summary.PendingPayments,
		TotalAmount:        summary.TotalAmount,
		RefundedAmount:     summary.RefundedAmount,
		NetAmount:          summary.NetAmount,
	}, nil
}

// Helper methods

func (s *PaymentService) generateTransactionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "txn_" + hex.EncodeToString(bytes), nil
}

func (s *PaymentService) generateRefundID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "ref_" + hex.EncodeToString(bytes), nil
}

// Response conversion methods

func (s *PaymentService) paymentToResponse(payment *domain.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		PaymentMethod:   payment.PaymentMethod,
		TransactionID:   payment.TransactionID,
		GatewayID:       payment.GatewayID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Status:          payment.Status,
		GatewayStatus:   payment.GatewayStatus,
		GatewayResponse: payment.GatewayResponse,
		FailureReason:   payment.FailureReason,
		ProcessedAt:     payment.ProcessedAt,
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
		Metadata:        payment.Metadata,
	}
}

func (s *PaymentService) paymentRefundToResponse(refund *domain.PaymentRefund) *dto.PaymentRefundResponse {
	return &dto.PaymentRefundResponse{
		ID:              refund.ID,
		PaymentID:       refund.PaymentID,
		RefundID:        refund.RefundID,
		Amount:          refund.Amount,
		Reason:          refund.Reason,
		Status:          refund.Status,
		GatewayResponse: refund.GatewayResponse,
		ProcessedAt:     refund.ProcessedAt,
		CreatedBy:       refund.CreatedBy,
		CreatedAt:       refund.CreatedAt,
	}
}
