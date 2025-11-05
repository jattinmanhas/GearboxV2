package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type CouponService interface {
	// Coupon Management
	CreateCoupon(ctx context.Context, req *dto.CreateCouponRequest) (*domain.Coupon, error)
	GetCouponByID(ctx context.Context, id int64) (*domain.Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error)
	UpdateCoupon(ctx context.Context, id int64, req *dto.UpdateCouponRequest) (*domain.Coupon, error)
	DeleteCoupon(ctx context.Context, id int64) error
	ListCoupons(ctx context.Context, req *dto.ListCouponsRequest) (*dto.ListCouponsResponse, error)

	// Coupon Validation & Usage
	ValidateCoupon(ctx context.Context, req *dto.ValidateCouponRequest) (*dto.ValidateCouponResponse, error)
	RecordCouponUsage(ctx context.Context, couponID int64, cartID *int64, orderID *int64, userID *int64, discountAmount float64) error
	GetCouponUsage(ctx context.Context, req *dto.ListCouponUsageRequest) (*dto.ListCouponUsageResponse, error)

	// Cart Integration
	ApplyCouponToCart(ctx context.Context, cartID int64, couponCode string, userID *int64) (*domain.Coupon, float64, error)
	RemoveCouponFromCart(ctx context.Context, cartID int64, couponCode string) error
	GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.Coupon, error)

	// Discount Calculation
	CalculateDiscount(ctx context.Context, coupon *domain.Coupon, cartAmount float64) (float64, error)
}

type couponService struct {
	couponRepo repository.CouponRepository
	cartRepo   repository.CartRepository
}

func NewCouponService(couponRepo repository.CouponRepository, cartRepo repository.CartRepository) CouponService {
	return &couponService{
		couponRepo: couponRepo,
		cartRepo:   cartRepo,
	}
}

// Coupon Management

// CreateCoupon creates a new coupon
func (s *couponService) CreateCoupon(ctx context.Context, req *dto.CreateCouponRequest) (*domain.Coupon, error) {
	// Check if coupon code already exists
	_, err := s.couponRepo.GetCouponByCode(ctx, req.Code)
	if err == nil {
		return nil, fmt.Errorf("coupon with code %s already exists", req.Code)
	}

	// Set default values
	startsAt := time.Now()
	if req.StartsAt != nil {
		startsAt = *req.StartsAt
	}

	coupon := &domain.Coupon{
		Code:            req.Code,
		Name:            req.Name,
		Description:     req.Description,
		Type:            req.Type,
		Value:           req.Value,
		MinimumAmount:   getFloat64Value(req.MinimumAmount),
		MaximumDiscount: req.MaximumDiscount,
		UsageLimit:      req.UsageLimit,
		IsActive:        req.IsActive,
		StartsAt:        startsAt,
		ExpiresAt:       req.ExpiresAt,
	}

	err = s.couponRepo.CreateCoupon(ctx, coupon)
	if err != nil {
		return nil, fmt.Errorf("failed to create coupon: %w", err)
	}

	return coupon, nil
}

// GetCouponByID retrieves a coupon by ID
func (s *couponService) GetCouponByID(ctx context.Context, id int64) (*domain.Coupon, error) {
	coupon, err := s.couponRepo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get coupon: %w", err)
	}

	return coupon, nil
}

// GetCouponByCode retrieves a coupon by code
func (s *couponService) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	coupon, err := s.couponRepo.GetCouponByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get coupon: %w", err)
	}

	return coupon, nil
}

// UpdateCoupon updates an existing coupon
func (s *couponService) UpdateCoupon(ctx context.Context, id int64, req *dto.UpdateCouponRequest) (*domain.Coupon, error) {
	// Get existing coupon
	existingCoupon, err := s.couponRepo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing coupon: %w", err)
	}

	// Update fields that are provided
	updateCoupon := *existingCoupon

	if req.Name != nil {
		updateCoupon.Name = *req.Name
	}
	if req.Description != nil {
		updateCoupon.Description = *req.Description
	}
	if req.Type != nil {
		updateCoupon.Type = *req.Type
	}
	if req.Value != nil {
		updateCoupon.Value = *req.Value
	}
	if req.MinimumAmount != nil {
		updateCoupon.MinimumAmount = *req.MinimumAmount
	}
	if req.MaximumDiscount != nil {
		updateCoupon.MaximumDiscount = req.MaximumDiscount
	}
	if req.UsageLimit != nil {
		updateCoupon.UsageLimit = req.UsageLimit
	}
	if req.IsActive != nil {
		updateCoupon.IsActive = *req.IsActive
	}
	if req.StartsAt != nil {
		updateCoupon.StartsAt = *req.StartsAt
	}
	if req.ExpiresAt != nil {
		updateCoupon.ExpiresAt = req.ExpiresAt
	}

	updateCoupon.UpdatedAt = time.Now()

	// Update coupon in repository
	err = s.couponRepo.UpdateCoupon(ctx, id, &updateCoupon)
	if err != nil {
		return nil, fmt.Errorf("failed to update coupon: %w", err)
	}

	return &updateCoupon, nil
}

// DeleteCoupon deletes a coupon
func (s *couponService) DeleteCoupon(ctx context.Context, id int64) error {
	// Check if coupon exists
	_, err := s.couponRepo.GetCouponByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get coupon: %w", err)
	}

	// Delete coupon
	err = s.couponRepo.DeleteCoupon(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete coupon: %w", err)
	}

	return nil
}

// ListCoupons retrieves coupons with filters
func (s *couponService) ListCoupons(ctx context.Context, req *dto.ListCouponsRequest) (*dto.ListCouponsResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// Build filter
	filter := &domain.CouponFilter{
		Code:     req.Code,
		Type:     req.Type,
		IsActive: req.IsActive,
		Search:   req.Search,
	}

	// Get coupons from repository
	coupons, total, err := s.couponRepo.ListCoupons(ctx, filter, offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list coupons: %w", err)
	}

	// Convert to response DTOs
	couponResponses := make([]dto.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		couponResponses[i] = dto.CouponResponse{
			ID:              coupon.ID,
			Code:            coupon.Code,
			Name:            coupon.Name,
			Description:     coupon.Description,
			Type:            coupon.Type,
			Value:           coupon.Value,
			MinimumAmount:   coupon.MinimumAmount,
			MaximumDiscount: coupon.MaximumDiscount,
			UsageLimit:      coupon.UsageLimit,
			UsedCount:       coupon.UsedCount,
			IsActive:        coupon.IsActive,
			StartsAt:        coupon.StartsAt,
			ExpiresAt:       coupon.ExpiresAt,
			CreatedAt:       coupon.CreatedAt,
			UpdatedAt:       coupon.UpdatedAt,
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.ListCouponsResponse{
		Coupons:    couponResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Coupon Validation & Usage

// ValidateCoupon validates a coupon for use
func (s *couponService) ValidateCoupon(ctx context.Context, req *dto.ValidateCouponRequest) (*dto.ValidateCouponResponse, error) {
	// Validate coupon
	coupon, err := s.couponRepo.IsCouponValid(ctx, req.CouponCode, req.CartAmount, req.UserID)
	if err != nil {
		return &dto.ValidateCouponResponse{
			Valid:          false,
			DiscountAmount: 0,
			Message:        err.Error(),
		}, nil
	}

	// Calculate discount amount
	discountAmount, err := s.couponRepo.CalculateDiscount(ctx, coupon, req.CartAmount)
	if err != nil {
		return &dto.ValidateCouponResponse{
			Valid:          false,
			DiscountAmount: 0,
			Message:        err.Error(),
		}, nil
	}

	// Convert coupon to response
	couponResponse := &dto.CouponResponse{
		ID:              coupon.ID,
		Code:            coupon.Code,
		Name:            coupon.Name,
		Description:     coupon.Description,
		Type:            coupon.Type,
		Value:           coupon.Value,
		MinimumAmount:   coupon.MinimumAmount,
		MaximumDiscount: coupon.MaximumDiscount,
		UsageLimit:      coupon.UsageLimit,
		UsedCount:       coupon.UsedCount,
		IsActive:        coupon.IsActive,
		StartsAt:        coupon.StartsAt,
		ExpiresAt:       coupon.ExpiresAt,
		CreatedAt:       coupon.CreatedAt,
		UpdatedAt:       coupon.UpdatedAt,
	}

	return &dto.ValidateCouponResponse{
		Valid:          true,
		DiscountAmount: discountAmount,
		Message:        "Coupon is valid",
		Coupon:         couponResponse,
	}, nil
}

// RecordCouponUsage records coupon usage
func (s *couponService) RecordCouponUsage(ctx context.Context, couponID int64, cartID *int64, orderID *int64, userID *int64, discountAmount float64) error {
	usage := &domain.CouponUsage{
		CouponID:       couponID,
		OrderID:        orderID,
		UserID:         userID,
		CartID:         cartID,
		DiscountAmount: discountAmount,
	}

	err := s.couponRepo.RecordCouponUsage(ctx, usage)
	if err != nil {
		return fmt.Errorf("failed to record coupon usage: %w", err)
	}

	return nil
}

// GetCouponUsage retrieves coupon usage with filters
func (s *couponService) GetCouponUsage(ctx context.Context, req *dto.ListCouponUsageRequest) (*dto.ListCouponUsageResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// Build filter
	filter := &domain.CouponUsageFilter{
		CouponID: req.CouponID,
		UserID:   req.UserID,
		OrderID:  req.OrderID,
	}

	// Get usage from repository
	usage, total, err := s.couponRepo.GetCouponUsage(ctx, filter, offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list coupon usage: %w", err)
	}

	// Convert to response DTOs
	usageResponses := make([]dto.CouponUsageResponse, len(usage))
	for i, u := range usage {
		usageResponses[i] = dto.CouponUsageResponse{
			ID:             u.ID,
			CouponID:       u.CouponID,
			OrderID:        u.OrderID,
			UserID:         u.UserID,
			CartID:         u.CartID,
			DiscountAmount: u.DiscountAmount,
			CreatedAt:      u.CreatedAt,
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.ListCouponUsageResponse{
		Usage:      usageResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Cart Integration

// ApplyCouponToCart applies a coupon to a cart
func (s *couponService) ApplyCouponToCart(ctx context.Context, cartID int64, couponCode string, userID *int64) (*domain.Coupon, float64, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get cart total
	cartTotal, err := s.cartRepo.CalculateCartTotal(ctx, cartID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to calculate cart total: %w", err)
	}

	// Validate coupon
	coupon, err := s.couponRepo.IsCouponValid(ctx, couponCode, cartTotal, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("coupon validation failed: %w", err)
	}

	// Calculate discount amount
	discountAmount, err := s.couponRepo.CalculateDiscount(ctx, coupon, cartTotal)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to calculate discount: %w", err)
	}

	// Check if coupon is already applied to cart
	_, err = s.cartRepo.GetCartCouponByCode(ctx, cartID, couponCode)
	if err == nil {
		return nil, 0, fmt.Errorf("coupon %s is already applied to this cart", couponCode)
	}

	// Check if any coupon is already applied to cart (enforce single coupon)
	existingCoupons, err := s.cartRepo.GetCartCoupons(ctx, cartID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check existing coupons: %w", err)
	}
	if len(existingCoupons) > 0 {
		return nil, 0, fmt.Errorf("only one coupon can be applied per cart. Please remove the existing coupon first")
	}

	// Apply coupon to cart
	cartCoupon := &domain.CartCoupon{
		CartID:         cartID,
		CouponCode:     couponCode,
		DiscountAmount: discountAmount,
	}

	err = s.cartRepo.ApplyCouponToCart(ctx, cartCoupon)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to apply coupon to cart: %w", err)
	}

	// Record usage
	err = s.RecordCouponUsage(ctx, coupon.ID, &cartID, nil, userID, discountAmount)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to record coupon usage: %v\n", err)
	}

	return coupon, discountAmount, nil
}

// RemoveCouponFromCart removes a coupon from a cart
func (s *couponService) RemoveCouponFromCart(ctx context.Context, cartID int64, couponCode string) error {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Remove coupon from cart
	err = s.cartRepo.RemoveCouponFromCart(ctx, cartID, couponCode)
	if err != nil {
		return fmt.Errorf("failed to remove coupon from cart: %w", err)
	}

	return nil
}

// GetCartCoupons retrieves all coupons applied to a cart
func (s *couponService) GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.Coupon, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get cart coupons
	cartCoupons, err := s.cartRepo.GetCartCoupons(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart coupons: %w", err)
	}

	// Convert to domain coupons
	coupons := make([]*domain.Coupon, len(cartCoupons))
	for i, cartCoupon := range cartCoupons {
		coupon, err := s.couponRepo.GetCouponByCode(ctx, cartCoupon.CouponCode)
		if err != nil {
			// Skip invalid coupons
			continue
		}
		coupons[i] = coupon
	}

	return coupons, nil
}

// CalculateDiscount calculates the discount amount for a coupon based on cart amount
func (s *couponService) CalculateDiscount(ctx context.Context, coupon *domain.Coupon, cartAmount float64) (float64, error) {
	return s.couponRepo.CalculateDiscount(ctx, coupon, cartAmount)
}

// Helper functions

func getFloat64Value(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
