package dto

// CreateCartRequest represents the request to create a new cart
type CreateCartRequest struct {
	UserID    *int64 `json:"user_id" validate:"omitempty"`
	SessionID string `json:"session_id" validate:"omitempty"`
	Currency  string `json:"currency" validate:"required,len=3"`
}

// UpdateCartRequest represents the request to update an existing cart
type UpdateCartRequest struct {
	Currency *string `json:"currency" validate:"omitempty,len=3"`
}

// CartResponse represents the response for cart data
type CartResponse struct {
	ID        int64   `json:"id"`
	UserID    *int64  `json:"user_id"`
	SessionID string  `json:"session_id"`
	Currency  string  `json:"currency"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ExpiresAt *string `json:"expires_at"`
}

// AddToCartRequest represents the request to add an item to cart
type AddToCartRequest struct {
	ProductID        int64  `json:"product_id" validate:"required"`
	ProductVariantID *int64 `json:"product_variant_id" validate:"omitempty"`
	Quantity         int    `json:"quantity" validate:"required,min=1"`
}

// UpdateCartItemRequest represents the request to update a cart item
type UpdateCartItemRequest struct {
	Quantity *int `json:"quantity" validate:"omitempty,min=1"`
}

// CartItemResponse represents the response for cart item data
type CartItemResponse struct {
	ID               int64   `json:"id"`
	CartID           int64   `json:"cart_id"`
	ProductID        int64   `json:"product_id"`
	ProductVariantID *int64  `json:"product_variant_id"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unit_price"`
	TotalPrice       float64 `json:"total_price"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// CartSummaryResponse represents the response for cart summary
type CartSummaryResponse struct {
	CartID         int64              `json:"cart_id"`
	ItemCount      int                `json:"item_count"`
	Subtotal       float64            `json:"subtotal"`
	TaxAmount      float64            `json:"tax_amount"`
	ShippingAmount float64            `json:"shipping_amount"`
	DiscountAmount float64            `json:"discount_amount"`
	TotalAmount    float64            `json:"total_amount"`
	Currency       string             `json:"currency"`
	Items          []CartItemResponse `json:"items"`
}

// ApplyCouponRequest represents the request to apply a coupon to cart
type ApplyCouponRequest struct {
	CouponCode string `json:"coupon_code" validate:"required,min=1,max=50"`
}

// RemoveCouponRequest represents the request to remove a coupon from cart
type RemoveCouponRequest struct {
	CouponCode string `json:"coupon_code" validate:"required,min=1,max=50"`
}

// CartCouponResponse represents the response for cart coupon data
type CartCouponResponse struct {
	ID             int64   `json:"id"`
	CartID         int64   `json:"cart_id"`
	CouponCode     string  `json:"coupon_code"`
	DiscountAmount float64 `json:"discount_amount"`
	CreatedAt      string  `json:"created_at"`
}

// SetShippingRequest represents the request to set shipping for cart
type SetShippingRequest struct {
	ShippingMethodID int64   `json:"shipping_method_id" validate:"required"`
	ShippingMethod   string  `json:"shipping_method" validate:"required,min=1,max=100"`
	ShippingAmount   float64 `json:"shipping_amount" validate:"required,min=0"`
	EstimatedDays    int     `json:"estimated_days" validate:"required,min=1"`
}

// UpdateShippingRequest represents the request to update shipping for cart
type UpdateShippingRequest struct {
	ShippingMethodID *int64   `json:"shipping_method_id" validate:"omitempty"`
	ShippingMethod   *string  `json:"shipping_method" validate:"omitempty,min=1,max=100"`
	ShippingAmount   *float64 `json:"shipping_amount" validate:"omitempty,min=0"`
	EstimatedDays    *int     `json:"estimated_days" validate:"omitempty,min=1"`
}

// CartShippingResponse represents the response for cart shipping data
type CartShippingResponse struct {
	ID               int64   `json:"id"`
	CartID           int64   `json:"cart_id"`
	ShippingMethodID int64   `json:"shipping_method_id"`
	ShippingMethod   string  `json:"shipping_method"`
	ShippingAmount   float64 `json:"shipping_amount"`
	EstimatedDays    int     `json:"estimated_days"`
	CreatedAt        string  `json:"created_at"`
}

// Wishlist Management DTOs

// CreateWishlistRequest represents the request to create a new wishlist
type CreateWishlistRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	IsPublic bool   `json:"is_public"`
}

// UpdateWishlistRequest represents the request to update an existing wishlist
type UpdateWishlistRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=1,max=100"`
	IsPublic *bool   `json:"is_public"`
}

// WishlistResponse represents the response for wishlist data
type WishlistResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AddToWishlistRequest represents the request to add an item to wishlist
type AddToWishlistRequest struct {
	ProductID        int64  `json:"product_id" validate:"required"`
	ProductVariantID *int64 `json:"product_variant_id" validate:"omitempty"`
	Notes            string `json:"notes" validate:"omitempty,max=500"`
}

// UpdateWishlistItemRequest represents the request to update a wishlist item
type UpdateWishlistItemRequest struct {
	Notes *string `json:"notes" validate:"omitempty,max=500"`
}

// WishlistItemResponse represents the response for wishlist item data
type WishlistItemResponse struct {
	ID               int64  `json:"id"`
	WishlistID       int64  `json:"wishlist_id"`
	ProductID        int64  `json:"product_id"`
	ProductVariantID *int64 `json:"product_variant_id"`
	Notes            string `json:"notes"`
	CreatedAt        string `json:"created_at"`
}

// ListWishlistsResponse represents the response for listing wishlists
type ListWishlistsResponse struct {
	Wishlists  []WishlistResponse `json:"wishlists"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// ListWishlistItemsResponse represents the response for listing wishlist items
type ListWishlistItemsResponse struct {
	Items      []WishlistItemResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// Cart Operations DTOs

// MergeCartRequest represents the request to merge carts
type MergeCartRequest struct {
	SourceCartID int64 `json:"source_cart_id" validate:"required"`
}

// ClearCartRequest represents the request to clear cart
type ClearCartRequest struct {
	Confirm bool `json:"confirm" validate:"required"`
}

// CartAnalyticsResponse represents analytics data for cart
type CartAnalyticsResponse struct {
	TotalCarts          int64   `json:"total_carts"`
	ActiveCarts         int64   `json:"active_carts"`
	AbandonedCarts      int64   `json:"abandoned_carts"`
	AverageCartValue    float64 `json:"average_cart_value"`
	TotalCartValue      float64 `json:"total_cart_value"`
	ConversionRate      float64 `json:"conversion_rate"`
	AverageItemsPerCart float64 `json:"average_items_per_cart"`
}
