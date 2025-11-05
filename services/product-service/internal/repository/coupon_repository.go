package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CouponRepository interface {
	// Coupon CRUD
	CreateCoupon(ctx context.Context, coupon *domain.Coupon) error
	GetCouponByID(ctx context.Context, id int64) (*domain.Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error)
	UpdateCoupon(ctx context.Context, id int64, coupon *domain.Coupon) error
	DeleteCoupon(ctx context.Context, id int64) error
	ListCoupons(ctx context.Context, filter *domain.CouponFilter, offset, limit int) ([]*domain.Coupon, int64, error)

	// Coupon Usage
	RecordCouponUsage(ctx context.Context, usage *domain.CouponUsage) error
	GetCouponUsage(ctx context.Context, filter *domain.CouponUsageFilter, offset, limit int) ([]*domain.CouponUsage, int64, error)
	GetCouponUsageCount(ctx context.Context, couponID int64) (int, error)
	GetUserCouponUsageCount(ctx context.Context, couponID, userID int64) (int, error)

	// Coupon Validation
	IsCouponValid(ctx context.Context, code string, cartAmount float64, userID *int64) (*domain.Coupon, error)
	CalculateDiscount(ctx context.Context, coupon *domain.Coupon, cartAmount float64) (float64, error)
}

type couponRepository struct {
	db *sqlx.DB
}

func NewCouponRepository(db *sqlx.DB) CouponRepository {
	return &couponRepository{
		db: db,
	}
}

// Coupon CRUD

// CreateCoupon creates a new coupon
func (r *couponRepository) CreateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	query := `
		INSERT INTO coupons (code, name, description, type, value, minimum_amount, maximum_discount, 
			usage_limit, is_active, starts_at, expires_at, created_at, updated_at)
		VALUES (:code, :name, :description, :type, :value, :minimum_amount, :maximum_discount, 
			:usage_limit, :is_active, :starts_at, :expires_at, :created_at, :updated_at)
		RETURNING id`

	coupon.CreatedAt = time.Now()
	coupon.UpdatedAt = time.Now()

	rows, err := r.db.NamedQueryContext(ctx, query, coupon)
	if err != nil {
		return fmt.Errorf("failed to create coupon: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&coupon.ID); err != nil {
			return fmt.Errorf("failed to get coupon ID: %w", err)
		}
	}

	return nil
}

// GetCouponByID retrieves a coupon by ID
func (r *couponRepository) GetCouponByID(ctx context.Context, id int64) (*domain.Coupon, error) {
	query := `SELECT * FROM coupons WHERE id = $1`

	var coupon domain.Coupon
	err := r.db.GetContext(ctx, &coupon, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("coupon with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get coupon: %w", err)
	}

	return &coupon, nil
}

// GetCouponByCode retrieves a coupon by code
func (r *couponRepository) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	query := `SELECT * FROM coupons WHERE code = $1`

	var coupon domain.Coupon
	err := r.db.GetContext(ctx, &coupon, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("coupon with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to get coupon: %w", err)
	}

	return &coupon, nil
}

// UpdateCoupon updates an existing coupon
func (r *couponRepository) UpdateCoupon(ctx context.Context, id int64, coupon *domain.Coupon) error {
	query := `
		UPDATE coupons SET 
			name = :name, description = :description, type = :type, value = :value,
			minimum_amount = :minimum_amount, maximum_discount = :maximum_discount,
			usage_limit = :usage_limit, is_active = :is_active, starts_at = :starts_at,
			expires_at = :expires_at, updated_at = :updated_at
		WHERE id = :id`

	coupon.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, coupon)
	if err != nil {
		return fmt.Errorf("failed to update coupon: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("coupon with ID %d not found", id)
	}

	return nil
}

// DeleteCoupon deletes a coupon
func (r *couponRepository) DeleteCoupon(ctx context.Context, id int64) error {
	query := `DELETE FROM coupons WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete coupon: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("coupon with ID %d not found", id)
	}

	return nil
}

// ListCoupons retrieves coupons with filters
func (r *couponRepository) ListCoupons(ctx context.Context, filter *domain.CouponFilter, offset, limit int) ([]*domain.Coupon, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.Code != nil {
		whereClause += fmt.Sprintf(" AND code ILIKE $%d", argIndex)
		args = append(args, "%"+*filter.Code+"%")
		argIndex++
	}

	if filter.Type != nil {
		whereClause += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.IsActive != nil {
		whereClause += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d OR code ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM coupons %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count coupons: %w", err)
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT * FROM coupons %s 
		ORDER BY created_at DESC 
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var coupons []*domain.Coupon
	err = r.db.SelectContext(ctx, &coupons, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list coupons: %w", err)
	}

	return coupons, total, nil
}

// Coupon Usage

// RecordCouponUsage records coupon usage
func (r *couponRepository) RecordCouponUsage(ctx context.Context, usage *domain.CouponUsage) error {
	query := `
		INSERT INTO coupon_usage (coupon_id, order_id, user_id, cart_id, discount_amount, created_at)
		VALUES (:coupon_id, :order_id, :user_id, :cart_id, :discount_amount, :created_at)`

	usage.CreatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, usage)
	if err != nil {
		return fmt.Errorf("failed to record coupon usage: %w", err)
	}

	// Update used count
	updateQuery := `UPDATE coupons SET used_count = used_count + 1 WHERE id = $1`
	_, err = r.db.ExecContext(ctx, updateQuery, usage.CouponID)
	if err != nil {
		return fmt.Errorf("failed to update coupon used count: %w", err)
	}

	return nil
}

// GetCouponUsage retrieves coupon usage with filters
func (r *couponRepository) GetCouponUsage(ctx context.Context, filter *domain.CouponUsageFilter, offset, limit int) ([]*domain.CouponUsage, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.CouponID != nil {
		whereClause += fmt.Sprintf(" AND coupon_id = $%d", argIndex)
		args = append(args, *filter.CouponID)
		argIndex++
	}

	if filter.UserID != nil {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.OrderID != nil {
		whereClause += fmt.Sprintf(" AND order_id = $%d", argIndex)
		args = append(args, *filter.OrderID)
		argIndex++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM coupon_usage %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count coupon usage: %w", err)
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT * FROM coupon_usage %s 
		ORDER BY created_at DESC 
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var usage []*domain.CouponUsage
	err = r.db.SelectContext(ctx, &usage, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list coupon usage: %w", err)
	}

	return usage, total, nil
}

// GetCouponUsageCount gets total usage count for a coupon
func (r *couponRepository) GetCouponUsageCount(ctx context.Context, couponID int64) (int, error) {
	query := `SELECT COUNT(*) FROM coupon_usage WHERE coupon_id = $1`

	var count int
	err := r.db.GetContext(ctx, &count, query, couponID)
	if err != nil {
		return 0, fmt.Errorf("failed to get coupon usage count: %w", err)
	}

	return count, nil
}

// GetUserCouponUsageCount gets usage count for a specific user and coupon
func (r *couponRepository) GetUserCouponUsageCount(ctx context.Context, couponID, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM coupon_usage WHERE coupon_id = $1 AND user_id = $2`

	var count int
	err := r.db.GetContext(ctx, &count, query, couponID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user coupon usage count: %w", err)
	}

	return count, nil
}

// Coupon Validation

// IsCouponValid validates if a coupon can be used
func (r *couponRepository) IsCouponValid(ctx context.Context, code string, cartAmount float64, userID *int64) (*domain.Coupon, error) {
	// Get coupon
	coupon, err := r.GetCouponByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Check if coupon is active
	if !coupon.IsActive {
		return nil, fmt.Errorf("coupon is not active")
	}

	// Check if coupon has started
	if coupon.StartsAt.After(time.Now()) {
		return nil, fmt.Errorf("coupon has not started yet")
	}

	// Check if coupon has expired
	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("coupon has expired")
	}

	// Check minimum amount
	if cartAmount < coupon.MinimumAmount {
		return nil, fmt.Errorf("minimum cart amount not met")
	}

	// Check usage limit
	if coupon.UsageLimit != nil && coupon.UsedCount >= *coupon.UsageLimit {
		return nil, fmt.Errorf("coupon usage limit exceeded")
	}

	return coupon, nil
}

// CalculateDiscount calculates the discount amount for a coupon
func (r *couponRepository) CalculateDiscount(ctx context.Context, coupon *domain.Coupon, cartAmount float64) (float64, error) {
	var discountAmount float64

	switch coupon.Type {
	case "percentage":
		discountAmount = cartAmount * (coupon.Value / 100)
		// Apply maximum discount if set (only if > 0)
		if coupon.MaximumDiscount != nil && *coupon.MaximumDiscount > 0 && discountAmount > *coupon.MaximumDiscount {
			discountAmount = *coupon.MaximumDiscount
		}
	case "fixed_amount":
		discountAmount = coupon.Value
		// Don't exceed cart amount
		if discountAmount > cartAmount {
			discountAmount = cartAmount
		}
	case "free_shipping":
		// Free shipping doesn't reduce cart amount, it's handled separately
		discountAmount = 0
	default:
		return 0, fmt.Errorf("invalid coupon type: %s", coupon.Type)
	}

	return discountAmount, nil
}
