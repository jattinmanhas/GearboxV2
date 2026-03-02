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

// Payments

// CreatePayment creates a new payment
func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (order_id, payment_method, transaction_id, gateway_id, amount, currency, 
		                     status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at)
		VALUES (:order_id, :payment_method, :transaction_id, :gateway_id, :amount, :currency, 
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
		SELECT id, order_id, payment_method, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments WHERE id = $1`

	payment := &domain.Payment{}
	var metadataJSON []byte

	// Use QueryRowContext to manually scan fields, handling metadata as []byte
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.GatewayID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.GatewayStatus,
		&payment.GatewayResponse,
		&payment.FailureReason,
		&payment.ProcessedAt,
		&metadataJSON, // Scan metadata as []byte
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Unmarshal metadata if present
	if len(metadataJSON) > 0 && string(metadataJSON) != "null" {
		if err := json.Unmarshal(metadataJSON, &payment.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return payment, nil
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (r *PaymentRepository) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, payment_method, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments WHERE transaction_id = $1`

	payment := &domain.Payment{}
	var metadataJSON []byte

	// Use QueryRowContext to manually scan fields, handling metadata as []byte
	row := r.db.QueryRowContext(ctx, query, transactionID)
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.GatewayID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.GatewayStatus,
		&payment.GatewayResponse,
		&payment.FailureReason,
		&payment.ProcessedAt,
		&metadataJSON, // Scan metadata as []byte
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Unmarshal metadata if present
	if len(metadataJSON) > 0 && string(metadataJSON) != "null" {
		if err := json.Unmarshal(metadataJSON, &payment.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return payment, nil
}

// GetLatestPaymentByOrderID retrieves the latest payment by order ID
func (r *PaymentRepository) GetLatestPaymentByOrderID(ctx context.Context, orderID int64) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, payment_method, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments
		WHERE order_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	payment := &domain.Payment{}
	var metadataJSON []byte

	row := r.db.QueryRowContext(ctx, query, orderID)
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.GatewayID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.GatewayStatus,
		&payment.GatewayResponse,
		&payment.FailureReason,
		&payment.ProcessedAt,
		&metadataJSON,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by order ID: %w", err)
	}

	if len(metadataJSON) > 0 && string(metadataJSON) != "null" {
		if err := json.Unmarshal(metadataJSON, &payment.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return payment, nil
}

// ListPayments retrieves payments with filtering and pagination
func (r *PaymentRepository) ListPayments(ctx context.Context, filter *domain.PaymentFilter, offset, limit int) ([]*domain.Payment, error) {
	query := `
		SELECT id, order_id, payment_method, transaction_id, gateway_id, amount, currency,
		       status, gateway_status, gateway_response, failure_reason, processed_at, metadata, created_at, updated_at
		FROM payments WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.OrderID != nil {
		query += fmt.Sprintf(" AND order_id = $%d", argIndex)
		args = append(args, *filter.OrderID)
		argIndex++
	}

	if filter.PaymentMethod != nil {
		query += fmt.Sprintf(" AND payment_method = $%d", argIndex)
		args = append(args, *filter.PaymentMethod)
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

	// Use QueryContext to manually scan rows, handling metadata as []byte
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		payment := &domain.Payment{}
		var metadataJSON []byte

		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.PaymentMethod,
			&payment.TransactionID,
			&payment.GatewayID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.GatewayStatus,
			&payment.GatewayResponse,
			&payment.FailureReason,
			&payment.ProcessedAt,
			&metadataJSON, // Scan metadata as []byte
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		// Unmarshal metadata if present
		if len(metadataJSON) > 0 && string(metadataJSON) != "null" {
			if err := json.Unmarshal(metadataJSON, &payment.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate payments: %w", err)
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

	if filter.PaymentMethod != nil {
		query += fmt.Sprintf(" AND payment_method = $%d", argIndex)
		args = append(args, *filter.PaymentMethod)
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
		SET order_id = :order_id, payment_method = :payment_method, transaction_id = :transaction_id, gateway_id = :gateway_id,
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
