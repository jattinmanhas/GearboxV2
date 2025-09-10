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

// Payment Methods

// CreatePaymentMethod creates a new payment method
func (s *PaymentService) CreatePaymentMethod(ctx context.Context, req *dto.PaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	// Check if payment method with same code already exists
	existing, err := s.paymentRepo.GetPaymentMethodByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing payment method: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("payment method with code %s already exists", req.Code)
	}

	// If this is set as default, unset other defaults
	if req.IsDefault {
		if err := s.unsetDefaultPaymentMethod(ctx); err != nil {
			return nil, fmt.Errorf("failed to unset default payment method: %w", err)
		}
	}

	pm := &domain.PaymentMethod{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		IsActive:    req.IsActive,
		IsDefault:   req.IsDefault,
		SortOrder:   req.SortOrder,
		Description: req.Description,
		Icon:        req.Icon,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.paymentRepo.CreatePaymentMethod(ctx, pm); err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	return s.paymentMethodToResponse(pm), nil
}

// GetPaymentMethod retrieves a payment method by ID
func (s *PaymentService) GetPaymentMethod(ctx context.Context, id int64) (*dto.PaymentMethodResponse, error) {
	pm, err := s.paymentRepo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	if pm == nil {
		return nil, fmt.Errorf("payment method not found")
	}

	return s.paymentMethodToResponse(pm), nil
}

// ListPaymentMethods retrieves payment methods with filtering
func (s *PaymentService) ListPaymentMethods(ctx context.Context, filter *domain.PaymentMethodFilter, page, limit int) ([]*dto.PaymentMethodResponse, error) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	methods, err := s.paymentRepo.ListPaymentMethods(ctx, filter, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}

	var responses []*dto.PaymentMethodResponse
	for _, method := range methods {
		responses = append(responses, s.paymentMethodToResponse(method))
	}

	return responses, nil
}

// UpdatePaymentMethod updates a payment method
func (s *PaymentService) UpdatePaymentMethod(ctx context.Context, id int64, req *dto.PaymentMethodRequest) (*dto.PaymentMethodResponse, error) {
	pm, err := s.paymentRepo.GetPaymentMethodByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	if pm == nil {
		return nil, fmt.Errorf("payment method not found")
	}

	// Check if code is being changed and if new code already exists
	if pm.Code != req.Code {
		existing, err := s.paymentRepo.GetPaymentMethodByCode(ctx, req.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing payment method: %w", err)
		}
		if existing != nil {
			return nil, fmt.Errorf("payment method with code %s already exists", req.Code)
		}
	}

	// If this is set as default, unset other defaults
	if req.IsDefault && !pm.IsDefault {
		if err := s.unsetDefaultPaymentMethod(ctx); err != nil {
			return nil, fmt.Errorf("failed to unset default payment method: %w", err)
		}
	}

	pm.Name = req.Name
	pm.Code = req.Code
	pm.Type = req.Type
	pm.IsActive = req.IsActive
	pm.IsDefault = req.IsDefault
	pm.SortOrder = req.SortOrder
	pm.Description = req.Description
	pm.Icon = req.Icon
	pm.UpdatedAt = time.Now()

	if err := s.paymentRepo.UpdatePaymentMethod(ctx, pm); err != nil {
		return nil, fmt.Errorf("failed to update payment method: %w", err)
	}

	return s.paymentMethodToResponse(pm), nil
}

// DeletePaymentMethod deletes a payment method
func (s *PaymentService) DeletePaymentMethod(ctx context.Context, id int64) error {
	// Check if payment method is being used in any payments
	payments, err := s.paymentRepo.ListPayments(ctx, &domain.PaymentFilter{PaymentMethodID: &id}, 0, 1)
	if err != nil {
		return fmt.Errorf("failed to check payment method usage: %w", err)
	}
	if len(payments) > 0 {
		return fmt.Errorf("cannot delete payment method: it is being used in existing payments")
	}

	if err := s.paymentRepo.DeletePaymentMethod(ctx, id); err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	return nil
}

// Payment Gateways

// CreatePaymentGateway creates a new payment gateway
func (s *PaymentService) CreatePaymentGateway(ctx context.Context, req *dto.PaymentGatewayRequest) (*dto.PaymentGatewayResponse, error) {
	// Check if payment gateway with same code already exists
	existing, err := s.paymentRepo.GetPaymentGatewayByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing payment gateway: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("payment gateway with code %s already exists", req.Code)
	}

	pg := &domain.PaymentGateway{
		Name:       req.Name,
		Code:       req.Code,
		IsActive:   req.IsActive,
		IsTestMode: req.IsTestMode,
		Config:     req.Config,
		WebhookURL: req.WebhookURL,
		SortOrder:  req.SortOrder,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.paymentRepo.CreatePaymentGateway(ctx, pg); err != nil {
		return nil, fmt.Errorf("failed to create payment gateway: %w", err)
	}

	return s.paymentGatewayToResponse(pg), nil
}

// GetPaymentGateway retrieves a payment gateway by ID
func (s *PaymentService) GetPaymentGateway(ctx context.Context, id int64) (*dto.PaymentGatewayResponse, error) {
	pg, err := s.paymentRepo.GetPaymentGatewayByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment gateway: %w", err)
	}
	if pg == nil {
		return nil, fmt.Errorf("payment gateway not found")
	}

	return s.paymentGatewayToResponse(pg), nil
}

// ListPaymentGateways retrieves payment gateways with filtering
func (s *PaymentService) ListPaymentGateways(ctx context.Context, filter *domain.PaymentGatewayFilter) ([]*dto.PaymentGatewayResponse, error) {
	gateways, err := s.paymentRepo.ListPaymentGateways(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment gateways: %w", err)
	}

	var responses []*dto.PaymentGatewayResponse
	for _, gateway := range gateways {
		responses = append(responses, s.paymentGatewayToResponse(gateway))
	}

	return responses, nil
}

// UpdatePaymentGateway updates a payment gateway
func (s *PaymentService) UpdatePaymentGateway(ctx context.Context, id int64, req *dto.PaymentGatewayRequest) (*dto.PaymentGatewayResponse, error) {
	pg, err := s.paymentRepo.GetPaymentGatewayByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment gateway: %w", err)
	}
	if pg == nil {
		return nil, fmt.Errorf("payment gateway not found")
	}

	// Check if code is being changed and if new code already exists
	if pg.Code != req.Code {
		existing, err := s.paymentRepo.GetPaymentGatewayByCode(ctx, req.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing payment gateway: %w", err)
		}
		if existing != nil {
			return nil, fmt.Errorf("payment gateway with code %s already exists", req.Code)
		}
	}

	pg.Name = req.Name
	pg.Code = req.Code
	pg.IsActive = req.IsActive
	pg.IsTestMode = req.IsTestMode
	pg.Config = req.Config
	pg.WebhookURL = req.WebhookURL
	pg.SortOrder = req.SortOrder
	pg.UpdatedAt = time.Now()

	if err := s.paymentRepo.UpdatePaymentGateway(ctx, pg); err != nil {
		return nil, fmt.Errorf("failed to update payment gateway: %w", err)
	}

	return s.paymentGatewayToResponse(pg), nil
}

// DeletePaymentGateway deletes a payment gateway
func (s *PaymentService) DeletePaymentGateway(ctx context.Context, id int64) error {
	// Check if payment gateway is being used in any payments
	// Get gateway first to get the code
	gateway, err := s.paymentRepo.GetPaymentGatewayByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get payment gateway: %w", err)
	}
	if gateway == nil {
		return fmt.Errorf("payment gateway not found")
	}

	payments, err := s.paymentRepo.ListPayments(ctx, &domain.PaymentFilter{GatewayID: &gateway.Code}, 0, 1)
	if err != nil {
		return fmt.Errorf("failed to check payment gateway usage: %w", err)
	}
	if len(payments) > 0 {
		return fmt.Errorf("cannot delete payment gateway: it is being used in existing payments")
	}

	if err := s.paymentRepo.DeletePaymentGateway(ctx, id); err != nil {
		return fmt.Errorf("failed to delete payment gateway: %w", err)
	}

	return nil
}

// Payments

// CreatePayment creates a new payment
func (s *PaymentService) CreatePayment(ctx context.Context, req *dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	// Validate payment method exists and is active
	pm, err := s.paymentRepo.GetPaymentMethodByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}
	if pm == nil {
		return nil, fmt.Errorf("payment method not found")
	}
	if !pm.IsActive {
		return nil, fmt.Errorf("payment method is not active")
	}

	// Validate payment gateway exists and is active
	pg, err := s.paymentRepo.GetPaymentGatewayByCode(ctx, req.GatewayID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment gateway: %w", err)
	}
	if pg == nil {
		return nil, fmt.Errorf("payment gateway not found")
	}
	if !pg.IsActive {
		return nil, fmt.Errorf("payment gateway is not active")
	}

	// Generate unique transaction ID
	transactionID, err := s.generateTransactionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction ID: %w", err)
	}

	payment := &domain.Payment{
		OrderID:         req.OrderID,
		PaymentMethodID: req.PaymentMethodID,
		TransactionID:   transactionID,
		GatewayID:       req.GatewayID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Status:          domain.PaymentStatusPending,
		Metadata:        req.Metadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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

func (s *PaymentService) unsetDefaultPaymentMethod(ctx context.Context) error {
	// This would require a custom repository method to unset all defaults
	// For now, we'll implement it in the repository
	return nil
}

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

func (s *PaymentService) paymentMethodToResponse(pm *domain.PaymentMethod) *dto.PaymentMethodResponse {
	return &dto.PaymentMethodResponse{
		ID:          pm.ID,
		Name:        pm.Name,
		Code:        pm.Code,
		Type:        pm.Type,
		IsActive:    pm.IsActive,
		IsDefault:   pm.IsDefault,
		SortOrder:   pm.SortOrder,
		Description: pm.Description,
		Icon:        pm.Icon,
	}
}

func (s *PaymentService) paymentGatewayToResponse(pg *domain.PaymentGateway) *dto.PaymentGatewayResponse {
	return &dto.PaymentGatewayResponse{
		ID:         pg.ID,
		Name:       pg.Name,
		Code:       pg.Code,
		IsActive:   pg.IsActive,
		IsTestMode: pg.IsTestMode,
		Config:     pg.Config,
		WebhookURL: pg.WebhookURL,
		SortOrder:  pg.SortOrder,
	}
}

func (s *PaymentService) paymentToResponse(payment *domain.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		PaymentMethodID: payment.PaymentMethodID,
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
