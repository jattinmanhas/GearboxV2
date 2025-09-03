package domain

import (
	"time"
)

// Cart represents a shopping cart
type Cart struct {
	ID        int64      `json:"id" db:"id"`
	UserID    *int64     `json:"user_id" db:"user_id"`       // null for guest carts
	SessionID string     `json:"session_id" db:"session_id"` // for guest carts
	Currency  string     `json:"currency" db:"currency"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
}

// CartItem represents items in a shopping cart
type CartItem struct {
	ID               int64     `json:"id" db:"id"`
	CartID           int64     `json:"cart_id" db:"cart_id"`
	ProductID        int64     `json:"product_id" db:"product_id"`
	ProductVariantID *int64    `json:"product_variant_id" db:"product_variant_id"`
	Quantity         int       `json:"quantity" db:"quantity"`
	UnitPrice        float64   `json:"unit_price" db:"unit_price"`
	TotalPrice       float64   `json:"total_price" db:"total_price"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// CartSummary represents a summary of cart contents
type CartSummary struct {
	CartID         int64      `json:"cart_id"`
	ItemCount      int        `json:"item_count"`
	Subtotal       float64    `json:"subtotal"`
	TaxAmount      float64    `json:"tax_amount"`
	ShippingAmount float64    `json:"shipping_amount"`
	DiscountAmount float64    `json:"discount_amount"`
	TotalAmount    float64    `json:"total_amount"`
	Currency       string     `json:"currency"`
	Items          []CartItem `json:"items"`
}

// Wishlist represents a user's wishlist
type Wishlist struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	IsPublic  bool      `json:"is_public" db:"is_public"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// WishlistItem represents items in a wishlist
type WishlistItem struct {
	ID               int64     `json:"id" db:"id"`
	WishlistID       int64     `json:"wishlist_id" db:"wishlist_id"`
	ProductID        int64     `json:"product_id" db:"product_id"`
	ProductVariantID *int64    `json:"product_variant_id" db:"product_variant_id"`
	Notes            string    `json:"notes" db:"notes"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// CartCoupon represents applied coupons to a cart
type CartCoupon struct {
	ID             int64     `json:"id" db:"id"`
	CartID         int64     `json:"cart_id" db:"cart_id"`
	CouponCode     string    `json:"coupon_code" db:"coupon_code"`
	DiscountAmount float64   `json:"discount_amount" db:"discount_amount"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CartShipping represents shipping information for a cart
type CartShipping struct {
	ID               int64     `json:"id" db:"id"`
	CartID           int64     `json:"cart_id" db:"cart_id"`
	ShippingMethodID int64     `json:"shipping_method_id" db:"shipping_method_id"`
	ShippingMethod   string    `json:"shipping_method" db:"shipping_method"`
	ShippingAmount   float64   `json:"shipping_amount" db:"shipping_amount"`
	EstimatedDays    int       `json:"estimated_days" db:"estimated_days"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// CartAnalytics represents analytics data for carts
type CartAnalytics struct {
	TotalCarts          int64   `json:"total_carts"`
	ActiveCarts         int64   `json:"active_carts"`
	AbandonedCarts      int64   `json:"abandoned_carts"`
	AverageCartValue    float64 `json:"average_cart_value"`
	TotalCartValue      float64 `json:"total_cart_value"`
	ConversionRate      float64 `json:"conversion_rate"`
	AverageItemsPerCart float64 `json:"average_items_per_cart"`
}
