package dto

import "time"

// CreateCouponRequest represents the request to create a new coupon
type CreateCouponRequest struct {
	Code            string     `json:"code" validate:"required,min=1,max=50"`
	Name            string     `json:"name" validate:"required,min=1,max=255"`
	Description     string     `json:"description" validate:"omitempty,max=1000"`
	Type            string     `json:"type" validate:"required,oneof=percentage fixed_amount free_shipping"`
	Value           float64    `json:"value" validate:"required,min=0"`
	MinimumAmount   *float64   `json:"minimum_amount" validate:"omitempty,min=0"`
	MaximumDiscount *float64   `json:"maximum_discount" validate:"omitempty,min=0"`
	UsageLimit      *int       `json:"usage_limit" validate:"omitempty,min=1"`
	IsActive        bool       `json:"is_active"`
	StartsAt        *time.Time `json:"starts_at"`
	ExpiresAt       *time.Time `json:"expires_at" validate:"omitempty,gtfield=StartsAt"`
}

// UpdateCouponRequest represents the request to update an existing coupon
type UpdateCouponRequest struct {
	Name            *string    `json:"name" validate:"omitempty,min=1,max=255"`
	Description     *string    `json:"description" validate:"omitempty,max=1000"`
	Type            *string    `json:"type" validate:"omitempty,oneof=percentage fixed_amount free_shipping"`
	Value           *float64   `json:"value" validate:"omitempty,min=0"`
	MinimumAmount   *float64   `json:"minimum_amount" validate:"omitempty,min=0"`
	MaximumDiscount *float64   `json:"maximum_discount" validate:"omitempty,min=0"`
	UsageLimit      *int       `json:"usage_limit" validate:"omitempty,min=1"`
	IsActive        *bool      `json:"is_active"`
	StartsAt        *time.Time `json:"starts_at"`
	ExpiresAt       *time.Time `json:"expires_at" validate:"omitempty,gtfield=StartsAt"`
}

// CouponResponse represents the response for coupon data
type CouponResponse struct {
	ID              int64      `json:"id"`
	Code            string     `json:"code"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Type            string     `json:"type"`
	Value           float64    `json:"value"`
	MinimumAmount   float64    `json:"minimum_amount"`
	MaximumDiscount *float64   `json:"maximum_discount"`
	UsageLimit      *int       `json:"usage_limit"`
	UsedCount       int        `json:"used_count"`
	IsActive        bool       `json:"is_active"`
	StartsAt        time.Time  `json:"starts_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ListCouponsRequest represents the request to list coupons with filters
type ListCouponsRequest struct {
	Code     *string `json:"code"`
	Type     *string `json:"type"`
	IsActive *bool   `json:"is_active"`
	Search   string  `json:"search"`
	Page     int     `json:"page"`
	Limit    int     `json:"limit"`
}

// ListCouponsResponse represents the response for listing coupons
type ListCouponsResponse struct {
	Coupons    []CouponResponse `json:"coupons"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// ValidateCouponRequest represents the request to validate a coupon
type ValidateCouponRequest struct {
	CouponCode string  `json:"coupon_code" validate:"required,min=1,max=50"`
	CartAmount float64 `json:"cart_amount" validate:"min=0"`
	UserID     *int64  `json:"user_id"`
}

// ValidateCouponResponse represents the response for coupon validation
type ValidateCouponResponse struct {
	Valid          bool            `json:"valid"`
	DiscountAmount float64         `json:"discount_amount"`
	Message        string          `json:"message"`
	Coupon         *CouponResponse `json:"coupon,omitempty"`
}

// CouponUsageResponse represents the response for coupon usage data
type CouponUsageResponse struct {
	ID             int64     `json:"id"`
	CouponID       int64     `json:"coupon_id"`
	OrderID        *int64    `json:"order_id"`
	UserID         *int64    `json:"user_id"`
	CartID         *int64    `json:"cart_id"`
	DiscountAmount float64   `json:"discount_amount"`
	CreatedAt      time.Time `json:"created_at"`
}

// ListCouponUsageRequest represents the request to list coupon usage
type ListCouponUsageRequest struct {
	CouponID *int64 `json:"coupon_id"`
	UserID   *int64 `json:"user_id"`
	OrderID  *int64 `json:"order_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// ListCouponUsageResponse represents the response for listing coupon usage
type ListCouponUsageResponse struct {
	Usage      []CouponUsageResponse `json:"usage"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}
