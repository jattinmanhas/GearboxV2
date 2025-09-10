package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

// PaymentRepository handles payment-related database operations
type PaymentRepository struct {
	db *sqlx.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Payment Methods

// CreatePaymentMethod creates a new payment method
func (r *PaymentRepository) CreatePaymentMethod(ctx context.Context, pm *domain.PaymentMethod) error {
	query := `
		INSERT INTO payment_methods (name, code, type, is_active, is_default, sort_order, description, icon, created_at, updated_at)
		VALUES (:name, :code, :type, :is_active, :is_default, :sort_order, :description, :icon, :created_at, :updated_at)
		RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, pm)
	if err != nil {
		return fmt.Errorf("failed to create payment method: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&pm.ID); err != nil {
			return fmt.Errorf("failed to scan payment method ID: %w", err)
		}
	}

	return nil
}

// GetPaymentMethodByID retrieves a payment method by ID
func (r *PaymentRepository) GetPaymentMethodByID(ctx context.Context, id int64) (*domain.PaymentMethod, error) {
	query := `
		SELECT id, name, code, type, is_active, is_default, sort_order, description, icon, created_at, updated_at
		FROM payment_methods WHERE id = $1`

	pm := &domain.PaymentMethod{}
	err := r.db.GetContext(ctx, pm, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return pm, nil
}

// GetPaymentMethodByCode retrieves a payment method by code
func (r *PaymentRepository) GetPaymentMethodByCode(ctx context.Context, code string) (*domain.PaymentMethod, error) {
	query := `
		SELECT id, name, code, type, is_active, is_default, sort_order, description, icon, created_at, updated_at
		FROM payment_methods WHERE code = $1`

	pm := &domain.PaymentMethod{}
	err := r.db.GetContext(ctx, pm, query, code)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return pm, nil
}

// ListPaymentMethods retrieves payment methods with filtering and pagination
func (r *PaymentRepository) ListPaymentMethods(ctx context.Context, filter *domain.PaymentMethodFilter, offset, limit int) ([]*domain.PaymentMethod, error) {
	query := `
		SELECT id, name, code, type, is_active, is_default, sort_order, description, icon, created_at, updated_at
		FROM payment_methods WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	query += " ORDER BY sort_order ASC, name ASC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var methods []*domain.PaymentMethod
	err := r.db.SelectContext(ctx, &methods, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}

	return methods, nil
}

// UpdatePaymentMethod updates a payment method
func (r *PaymentRepository) UpdatePaymentMethod(ctx context.Context, pm *domain.PaymentMethod) error {
	query := `
		UPDATE payment_methods 
		SET name = :name, code = :code, type = :type, is_active = :is_active, is_default = :is_default, 
		    sort_order = :sort_order, description = :description, icon = :icon, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, pm)
	if err != nil {
		return fmt.Errorf("failed to update payment method: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment method not found")
	}

	return nil
}

// DeletePaymentMethod deletes a payment method
func (r *PaymentRepository) DeletePaymentMethod(ctx context.Context, id int64) error {
	query := `DELETE FROM payment_methods WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment method not found")
	}

	return nil
}

// Payment Gateways

// CreatePaymentGateway creates a new payment gateway
func (r *PaymentRepository) CreatePaymentGateway(ctx context.Context, pg *domain.PaymentGateway) error {
	query := `
		INSERT INTO payment_gateways (name, code, is_active, is_test_mode, config, webhook_url, sort_order, created_at, updated_at)
		VALUES (:name, :code, :is_active, :is_test_mode, :config, :webhook_url, :sort_order, :created_at, :updated_at)
		RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, pg)
	if err != nil {
		return fmt.Errorf("failed to create payment gateway: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&pg.ID); err != nil {
			return fmt.Errorf("failed to scan payment gateway ID: %w", err)
		}
	}

	return nil
}

// GetPaymentGatewayByID retrieves a payment gateway by ID
func (r *PaymentRepository) GetPaymentGatewayByID(ctx context.Context, id int64) (*domain.PaymentGateway, error) {
	query := `
		SELECT id, name, code, is_active, is_test_mode, config, webhook_url, sort_order, created_at, updated_at
		FROM payment_gateways WHERE id = $1`

	pg := &domain.PaymentGateway{}
	err := r.db.GetContext(ctx, pg, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment gateway: %w", err)
	}

	return pg, nil
}

// GetPaymentGatewayByCode retrieves a payment gateway by code
func (r *PaymentRepository) GetPaymentGatewayByCode(ctx context.Context, code string) (*domain.PaymentGateway, error) {
	query := `
		SELECT id, name, code, is_active, is_test_mode, config, webhook_url, sort_order, created_at, updated_at
		FROM payment_gateways WHERE code = $1`

	pg := &domain.PaymentGateway{}
	err := r.db.GetContext(ctx, pg, query, code)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment gateway: %w", err)
	}

	return pg, nil
}

// ListPaymentGateways retrieves payment gateways with filtering
func (r *PaymentRepository) ListPaymentGateways(ctx context.Context, filter *domain.PaymentGatewayFilter) ([]*domain.PaymentGateway, error) {
	query := `
		SELECT id, name, code, is_active, is_test_mode, config, webhook_url, sort_order, created_at, updated_at
		FROM payment_gateways WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.Code != nil {
		query += fmt.Sprintf(" AND code = $%d", argIndex)
		args = append(args, *filter.Code)
		argIndex++
	}

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.IsTestMode != nil {
		query += fmt.Sprintf(" AND is_test_mode = $%d", argIndex)
		args = append(args, *filter.IsTestMode)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR code ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	query += " ORDER BY sort_order ASC, name ASC"

	var gateways []*domain.PaymentGateway
	err := r.db.SelectContext(ctx, &gateways, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment gateways: %w", err)
	}

	return gateways, nil
}

// UpdatePaymentGateway updates a payment gateway
func (r *PaymentRepository) UpdatePaymentGateway(ctx context.Context, pg *domain.PaymentGateway) error {
	query := `
		UPDATE payment_gateways 
		SET name = :name, code = :code, is_active = :is_active, is_test_mode = :is_test_mode, 
		    config = :config, webhook_url = :webhook_url, sort_order = :sort_order, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, pg)
	if err != nil {
		return fmt.Errorf("failed to update payment gateway: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment gateway not found")
	}

	return nil
}

// DeletePaymentGateway deletes a payment gateway
func (r *PaymentRepository) DeletePaymentGateway(ctx context.Context, id int64) error {
	query := `DELETE FROM payment_gateways WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment gateway: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment gateway not found")
	}

	return nil
}

// Payments

// CreatePayment creates a new payment
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (order_id, payment_method_id, transaction_id, gateway_id, amount, currency, 
		                     status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at)
		VALUES (:order_id, :payment_method_id, :transaction_id, :gateway_id, :amount, :currency, 
		        :status, :gateway_status, :gateway_response, :failure_reason, :processed_at, :metadata, :created_at, :updated_at)
		RETURNING id`

	// Convert metadata to JSON for storage
	paymentData := struct {
		*domain.Payment
		MetadataJSON []byte `db:"metadata"`
	}{
		Payment: payment,
	}

	if payment.Metadata != nil {
		var err error
		paymentData.MetadataJSON, err = json.Marshal(payment.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	rows, err := r.db.NamedQueryContext(ctx, query, paymentData)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&payment.ID); err != nil {
			return fmt.Errorf("failed to scan payment ID: %w", err)
		}
	}

	return nil
}

// GetPaymentByID retrieves a payment by ID
func (r *PaymentRepository) GetPaymentByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, payment_method_id, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments WHERE id = $1`

	payment := &domain.Payment{}
	var metadataJSON []byte

	err := r.db.GetContext(ctx, payment, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Handle metadata separately
	metadataQuery := `SELECT metadata FROM payments WHERE id = $1`
	err = r.db.GetContext(ctx, &metadataJSON, metadataQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment metadata: %w", err)
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &payment.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return payment, nil
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (r *PaymentRepository) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, payment_method_id, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments WHERE transaction_id = $1`

	payment := &domain.Payment{}
	var metadataJSON []byte

	err := r.db.GetContext(ctx, payment, query, transactionID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Handle metadata separately
	metadataQuery := `SELECT metadata FROM payments WHERE transaction_id = $1`
	err = r.db.GetContext(ctx, &metadataJSON, metadataQuery, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment metadata: %w", err)
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &payment.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return payment, nil
}

// ListPayments retrieves payments with filtering and pagination
func (r *PaymentRepository) ListPayments(ctx context.Context, filter *domain.PaymentFilter, offset, limit int) ([]*domain.Payment, error) {
	query := `
		SELECT id, order_id, payment_method_id, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.OrderID != nil {
		query += fmt.Sprintf(" AND order_id = $%d", argIndex)
		args = append(args, *filter.OrderID)
		argIndex++
	}

	if filter.PaymentMethodID != nil {
		query += fmt.Sprintf(" AND payment_method_id = $%d", argIndex)
		args = append(args, *filter.PaymentMethodID)
		argIndex++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.GatewayID != nil {
		query += fmt.Sprintf(" AND gateway_id = $%d", argIndex)
		args = append(args, *filter.GatewayID)
		argIndex++
	}

	if filter.Currency != nil {
		query += fmt.Sprintf(" AND currency = $%d", argIndex)
		args = append(args, *filter.Currency)
		argIndex++
	}

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.DateTo)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (transaction_id ILIKE $%d OR gateway_id ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	// Add sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder != "" {
		sortOrder = filter.SortOrder
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var payments []*domain.Payment
	err := r.db.SelectContext(ctx, &payments, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	// Handle metadata for each payment
	for _, payment := range payments {
		if payment.ID != 0 {
			var metadataJSON []byte
			metadataQuery := `SELECT metadata FROM payments WHERE id = $1`
			err := r.db.GetContext(ctx, &metadataJSON, metadataQuery, payment.ID)
			if err == nil && len(metadataJSON) > 0 {
				json.Unmarshal(metadataJSON, &payment.Metadata)
			}
		}
	}

	return payments, nil
}

// CountPayments counts payments with filtering
func (r *PaymentRepository) CountPayments(ctx context.Context, filter *domain.PaymentFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM payments WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.OrderID != nil {
		query += fmt.Sprintf(" AND order_id = $%d", argIndex)
		args = append(args, *filter.OrderID)
		argIndex++
	}

	if filter.PaymentMethodID != nil {
		query += fmt.Sprintf(" AND payment_method_id = $%d", argIndex)
		args = append(args, *filter.PaymentMethodID)
		argIndex++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.GatewayID != nil {
		query += fmt.Sprintf(" AND gateway_id = $%d", argIndex)
		args = append(args, *filter.GatewayID)
		argIndex++
	}

	if filter.Currency != nil {
		query += fmt.Sprintf(" AND currency = $%d", argIndex)
		args = append(args, *filter.Currency)
		argIndex++
	}

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.DateTo)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (transaction_id ILIKE $%d OR gateway_id ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count payments: %w", err)
	}

	return count, nil
}

// UpdatePayment updates a payment
func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments 
		SET order_id = :order_id, payment_method_id = :payment_method_id, transaction_id = :transaction_id, gateway_id = :gateway_id,
		    amount = :amount, currency = :currency, status = :status, gateway_status = :gateway_status,
		    gateway_response = :gateway_response, failure_reason = :failure_reason, processed_at = :processed_at, 
		    metadata = :metadata, updated_at = :updated_at
		WHERE id = :id`

	// Convert metadata to JSON for storage
	paymentData := struct {
		*domain.Payment
		MetadataJSON []byte `db:"metadata"`
	}{
		Payment: payment,
	}

	if payment.Metadata != nil {
		var err error
		paymentData.MetadataJSON, err = json.Marshal(payment.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	result, err := r.db.NamedExecContext(ctx, query, paymentData)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// UpdatePaymentStatus updates payment status
func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, id int64, status, gatewayStatus, gatewayResponse, failureReason string, processedAt *time.Time) error {
	query := `
		UPDATE payments 
		SET status = $2, gateway_status = $3, gateway_response = $4, failure_reason = $5, processed_at = $6, updated_at = $7
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status, gatewayStatus, gatewayResponse, failureReason, processedAt, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

// GetPaymentSummary retrieves payment summary statistics
func (r *PaymentRepository) GetPaymentSummary(ctx context.Context, dateFrom, dateTo *time.Time) (*domain.PaymentSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total_payments,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as successful_payments,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_payments,
			COUNT(CASE WHEN status IN ('pending', 'processing') THEN 1 END) as pending_payments,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN amount ELSE 0 END), 0) as total_amount,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN amount ELSE 0 END), 0) as refunded_amount
		FROM payments WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if dateFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *dateFrom)
		argIndex++
	}

	if dateTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *dateTo)
		argIndex++
	}

	summary := &domain.PaymentSummary{}
	err := r.db.GetContext(ctx, summary, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment summary: %w", err)
	}

	summary.NetAmount = summary.TotalAmount - summary.RefundedAmount

	if summary.TotalPayments > 0 {
		summary.AverageAmount = summary.TotalAmount / float64(summary.TotalPayments)
		summary.SuccessRate = float64(summary.SuccessfulPayments) / float64(summary.TotalPayments) * 100
	}

	return summary, nil
}

// Payment Refunds

// CreatePaymentRefund creates a new payment refund
func (r *PaymentRepository) CreatePaymentRefund(ctx context.Context, refund *domain.PaymentRefund) error {
	query := `
		INSERT INTO payment_refunds (payment_id, refund_id, amount, reason, status, gateway_response, processed_at, created_by, created_at)
		VALUES (:payment_id, :refund_id, :amount, :reason, :status, :gateway_response, :processed_at, :created_by, :created_at)
		RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, refund)
	if err != nil {
		return fmt.Errorf("failed to create payment refund: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&refund.ID); err != nil {
			return fmt.Errorf("failed to scan payment refund ID: %w", err)
		}
	}

	return nil
}

// GetPaymentRefundByID retrieves a payment refund by ID
func (r *PaymentRepository) GetPaymentRefundByID(ctx context.Context, id int64) (*domain.PaymentRefund, error) {
	query := `
		SELECT id, payment_id, refund_id, amount, reason, status, gateway_response, processed_at, created_by, created_at
		FROM payment_refunds WHERE id = $1`

	refund := &domain.PaymentRefund{}
	err := r.db.GetContext(ctx, refund, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment refund: %w", err)
	}

	return refund, nil
}

// ListPaymentRefunds retrieves payment refunds for a payment
func (r *PaymentRepository) ListPaymentRefunds(ctx context.Context, paymentID int64) ([]*domain.PaymentRefund, error) {
	query := `
		SELECT id, payment_id, refund_id, amount, reason, status, gateway_response, processed_at, created_by, created_at
		FROM payment_refunds WHERE payment_id = $1 ORDER BY created_at DESC`

	var refunds []*domain.PaymentRefund
	err := r.db.SelectContext(ctx, &refunds, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment refunds: %w", err)
	}

	return refunds, nil
}

// UpdatePaymentRefund updates a payment refund
func (r *PaymentRepository) UpdatePaymentRefund(ctx context.Context, refund *domain.PaymentRefund) error {
	query := `
		UPDATE payment_refunds 
		SET amount = :amount, reason = :reason, status = :status, gateway_response = :gateway_response, processed_at = :processed_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, refund)
	if err != nil {
		return fmt.Errorf("failed to update payment refund: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment refund not found")
	}

	return nil
}

// Payment Webhooks

// CreatePaymentWebhook creates a new payment webhook
func (r *PaymentRepository) CreatePaymentWebhook(ctx context.Context, webhook *domain.PaymentWebhook) error {
	query := `
		INSERT INTO payment_webhooks (gateway_id, event_type, event_id, payload, signature, is_processed, processed_at, created_at)
		VALUES (:gateway_id, :event_type, :event_id, :payload, :signature, :is_processed, :processed_at, :created_at)
		RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, webhook)
	if err != nil {
		return fmt.Errorf("failed to create payment webhook: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&webhook.ID); err != nil {
			return fmt.Errorf("failed to scan payment webhook ID: %w", err)
		}
	}

	return nil
}

// GetPaymentWebhookByID retrieves a payment webhook by ID
func (r *PaymentRepository) GetPaymentWebhookByID(ctx context.Context, id int64) (*domain.PaymentWebhook, error) {
	query := `
		SELECT id, gateway_id, event_type, event_id, payload, signature, is_processed, processed_at, created_at
		FROM payment_webhooks WHERE id = $1`

	webhook := &domain.PaymentWebhook{}
	err := r.db.GetContext(ctx, webhook, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment webhook: %w", err)
	}

	return webhook, nil
}

// ListUnprocessedWebhooks retrieves unprocessed webhooks
func (r *PaymentRepository) ListUnprocessedWebhooks(ctx context.Context, limit int) ([]*domain.PaymentWebhook, error) {
	query := `
		SELECT id, gateway_id, event_type, event_id, payload, signature, is_processed, processed_at, created_at
		FROM payment_webhooks WHERE is_processed = false ORDER BY created_at ASC LIMIT $1`

	var webhooks []*domain.PaymentWebhook
	err := r.db.SelectContext(ctx, &webhooks, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list unprocessed webhooks: %w", err)
	}

	return webhooks, nil
}

// MarkWebhookAsProcessed marks a webhook as processed
func (r *PaymentRepository) MarkWebhookAsProcessed(ctx context.Context, id int64) error {
	query := `UPDATE payment_webhooks SET is_processed = true, processed_at = $2 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to mark webhook as processed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook not found")
	}

	return nil
}
