package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	// Cart Management
	CreateCart(ctx context.Context, cart *domain.Cart) error
	GetCartByID(ctx context.Context, id int64) (*domain.Cart, error)
	GetCartByUserID(ctx context.Context, userID int64) (*domain.Cart, error)
	GetCartBySessionID(ctx context.Context, sessionID string) (*domain.Cart, error)
	GetCartBySessionOrUser(ctx context.Context, sessionID string, userID *int64) (*domain.Cart, error)
	UpdateCart(ctx context.Context, cart *domain.Cart) error
	DeleteCart(ctx context.Context, id int64) error
	GetOrCreateCart(ctx context.Context, userID *int64, sessionID string, currency string) (*domain.Cart, error)

	// Cart Items
	AddItemToCart(ctx context.Context, item *domain.CartItem) error
	GetCartItemByID(ctx context.Context, id int64) (*domain.CartItem, error)
	GetCartItemByProduct(ctx context.Context, cartID, productID int64, variantID *int64) (*domain.CartItem, error)
	UpdateCartItem(ctx context.Context, id int64, item *domain.CartItem) error
	DeleteCartItem(ctx context.Context, id int64) error
	GetCartItems(ctx context.Context, cartID int64) ([]*domain.CartItem, error)
	ClearCartItems(ctx context.Context, cartID int64) error

	// Cart Summary & Calculations
	GetCartSummary(ctx context.Context, cartID int64) (*domain.CartSummary, error)
	CalculateCartTotal(ctx context.Context, cartID int64) (float64, error)
	GetCartItemCount(ctx context.Context, cartID int64) (int, error)

	// Cart Coupons
	ApplyCouponToCart(ctx context.Context, cartCoupon *domain.CartCoupon) error
	RemoveCouponFromCart(ctx context.Context, cartID int64, couponCode string) error
	GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.CartCoupon, error)
	GetCartCouponByCode(ctx context.Context, cartID int64, couponCode string) (*domain.CartCoupon, error)

	// Cart Shipping
	SetCartShipping(ctx context.Context, shipping *domain.CartShipping) error
	UpdateCartShipping(ctx context.Context, cartID int64, shipping *domain.CartShipping) error
	GetCartShipping(ctx context.Context, cartID int64) (*domain.CartShipping, error)
	DeleteCartShipping(ctx context.Context, cartID int64) error

	// Cart Analytics & Management
	GetExpiredCarts(ctx context.Context, before time.Time) ([]*domain.Cart, error)
	DeleteExpiredCarts(ctx context.Context, before time.Time) error
	GetCartAnalytics(ctx context.Context) (*domain.CartAnalytics, error)
	MergeCarts(ctx context.Context, sourceCartID, targetCartID int64) error

	// Wishlist Management
	CreateWishlist(ctx context.Context, wishlist *domain.Wishlist) error
	GetWishlistByID(ctx context.Context, id int64) (*domain.Wishlist, error)
	GetWishlistsByUserID(ctx context.Context, userID int64, offset, limit int) ([]*domain.Wishlist, int64, error)
	UpdateWishlist(ctx context.Context, id int64, wishlist *domain.Wishlist) error
	DeleteWishlist(ctx context.Context, id int64) error

	// Wishlist Items
	AddItemToWishlist(ctx context.Context, item *domain.WishlistItem) error
	GetWishlistItemByID(ctx context.Context, id int64) (*domain.WishlistItem, error)
	GetWishlistItems(ctx context.Context, wishlistID int64, offset, limit int) ([]*domain.WishlistItem, int64, error)
	UpdateWishlistItem(ctx context.Context, id int64, item *domain.WishlistItem) error
	DeleteWishlistItem(ctx context.Context, id int64) error
	MoveItemToCart(ctx context.Context, wishlistItemID, cartID int64) error
}

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &cartRepository{
		db: db,
	}
}

// Cart Management

// CreateCart creates a new cart
func (r *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	query := `
		INSERT INTO carts (user_id, session_id, currency, created_at, updated_at, expires_at)
		VALUES (:user_id, :session_id, :currency, :created_at, :updated_at, :expires_at)
		RETURNING id`

	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	result, err := r.db.NamedQueryContext(ctx, query, cart)
	if err != nil {
		return fmt.Errorf("failed to create cart: %w", err)
	}
	defer result.Close()

	if result.Next() {
		if err := result.Scan(&cart.ID); err != nil {
			return fmt.Errorf("failed to scan cart ID: %w", err)
		}
	}

	return nil
}

// GetCartByID retrieves a cart by ID
func (r *cartRepository) GetCartByID(ctx context.Context, id int64) (*domain.Cart, error) {
	query := `SELECT * FROM carts WHERE id = $1`

	var cart domain.Cart
	err := r.db.GetContext(ctx, &cart, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart, nil
}

// GetCartByUserID retrieves a cart by user ID
func (r *cartRepository) GetCartByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	query := `SELECT * FROM carts WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`

	var cart domain.Cart
	err := r.db.GetContext(ctx, &cart, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart for user ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart, nil
}

// GetCartBySessionID retrieves a cart by session ID
func (r *cartRepository) GetCartBySessionID(ctx context.Context, sessionID string) (*domain.Cart, error) {
	query := `SELECT * FROM carts WHERE session_id = $1 ORDER BY created_at DESC LIMIT 1`

	var cart domain.Cart
	err := r.db.GetContext(ctx, &cart, query, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart for session ID %s not found", sessionID)
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart, nil
}

// UpdateCart updates an existing cart
func (r *cartRepository) UpdateCart(ctx context.Context, cart *domain.Cart) error {
	query := `
		UPDATE carts SET
			user_id = :user_id, session_id = :session_id, currency = :currency,
			updated_at = :updated_at, expires_at = :expires_at
		WHERE id = :id`

	cart.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, cart)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart with ID %d not found", cart.ID)
	}

	return nil
}

// GetCartBySessionOrUser retrieves a cart by session ID or user ID
func (r *cartRepository) GetCartBySessionOrUser(ctx context.Context, sessionID string, userID *int64) (*domain.Cart, error) {
	var cart domain.Cart
	var query string
	var args []interface{}

	if userID != nil {
		// If user is logged in, prioritize user-based cart
		query = `SELECT id, user_id, session_id, currency, created_at, updated_at, expires_at 
				 FROM carts 
				 WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
				 ORDER BY created_at DESC 
				 LIMIT 1`
		args = []interface{}{*userID}
	} else {
		// If guest user, use session ID
		query = `SELECT id, user_id, session_id, currency, created_at, updated_at, expires_at 
				 FROM carts 
				 WHERE session_id = $1 AND (expires_at IS NULL OR expires_at > NOW())
				 ORDER BY created_at DESC 
				 LIMIT 1`
		args = []interface{}{sessionID}
	}

	err := r.db.GetContext(ctx, &cart, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return &cart, nil
}

// DeleteCart deletes a cart
func (r *cartRepository) DeleteCart(ctx context.Context, id int64) error {
	// Start transaction to delete cart and all related data
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete cart items
	_, err = tx.ExecContext(ctx, "DELETE FROM cart_items WHERE cart_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cart items: %w", err)
	}

	// Delete cart coupons
	_, err = tx.ExecContext(ctx, "DELETE FROM cart_coupons WHERE cart_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cart coupons: %w", err)
	}

	// Delete cart shipping
	_, err = tx.ExecContext(ctx, "DELETE FROM cart_shipping WHERE cart_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cart shipping: %w", err)
	}

	// Delete cart
	result, err := tx.ExecContext(ctx, "DELETE FROM carts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart with ID %d not found", id)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetOrCreateCart gets an existing cart or creates a new one
func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID *int64, sessionID string, currency string) (*domain.Cart, error) {
	// Try to get existing cart
	var cart *domain.Cart
	var err error

	if userID != nil {
		cart, err = r.GetCartByUserID(ctx, *userID)
		if err == nil {
			return cart, nil
		}
	}

	if sessionID != "" {
		cart, err = r.GetCartBySessionID(ctx, sessionID)
		if err == nil {
			return cart, nil
		}
	}

	// Create new cart
	newCart := &domain.Cart{
		UserID:    userID,
		SessionID: sessionID,
		Currency:  currency,
	}

	err = r.CreateCart(ctx, newCart)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return newCart, nil
}

// Cart Items

// AddItemToCart adds an item to the cart
func (r *cartRepository) AddItemToCart(ctx context.Context, item *domain.CartItem) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, product_variant_id, quantity, unit_price, total_price, created_at, updated_at)
		VALUES (:cart_id, :product_id, :product_variant_id, :quantity, :unit_price, :total_price, :created_at, :updated_at)`

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.TotalPrice = item.UnitPrice * float64(item.Quantity)

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get item ID: %w", err)
	}

	item.ID = id
	return nil
}

// GetCartItemByID retrieves a cart item by ID
func (r *cartRepository) GetCartItemByID(ctx context.Context, id int64) (*domain.CartItem, error) {
	query := `SELECT * FROM cart_items WHERE id = $1`

	var item domain.CartItem
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart item with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	return &item, nil
}

// GetCartItemByProduct retrieves a cart item by product and variant
func (r *cartRepository) GetCartItemByProduct(ctx context.Context, cartID, productID int64, variantID *int64) (*domain.CartItem, error) {
	var query string
	var args []interface{}

	if variantID != nil {
		query = `SELECT * FROM cart_items WHERE cart_id = $1 AND product_id = $2 AND product_variant_id = $3`
		args = []interface{}{cartID, productID, *variantID}
	} else {
		query = `SELECT * FROM cart_items WHERE cart_id = $1 AND product_id = $2 AND product_variant_id IS NULL`
		args = []interface{}{cartID, productID}
	}

	var item domain.CartItem
	err := r.db.GetContext(ctx, &item, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart item not found")
		}
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	return &item, nil
}

// UpdateCartItem updates an existing cart item
func (r *cartRepository) UpdateCartItem(ctx context.Context, id int64, item *domain.CartItem) error {
	query := `
		UPDATE cart_items SET
			quantity = :quantity, unit_price = :unit_price, total_price = :total_price,
			updated_at = :updated_at
		WHERE id = :id`

	item.UpdatedAt = time.Now()
	item.TotalPrice = item.UnitPrice * float64(item.Quantity)
	item.ID = id

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart item with ID %d not found", id)
	}

	return nil
}

// DeleteCartItem deletes a cart item
func (r *cartRepository) DeleteCartItem(ctx context.Context, id int64) error {
	query := `DELETE FROM cart_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart item with ID %d not found", id)
	}

	return nil
}

// GetCartItems retrieves all items in a cart
func (r *cartRepository) GetCartItems(ctx context.Context, cartID int64) ([]*domain.CartItem, error) {
	query := `SELECT * FROM cart_items WHERE cart_id = $1 ORDER BY created_at ASC`

	var items []*domain.CartItem
	err := r.db.SelectContext(ctx, &items, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	return items, nil
}

// ClearCartItems removes all items from a cart
func (r *cartRepository) ClearCartItems(ctx context.Context, cartID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`

	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to clear cart items: %w", err)
	}

	return nil
}

// Cart Summary & Calculations

// GetCartSummary retrieves a complete cart summary
func (r *cartRepository) GetCartSummary(ctx context.Context, cartID int64) (*domain.CartSummary, error) {
	// Get cart
	cart, err := r.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get cart items
	items, err := r.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	// Calculate totals
	var subtotal float64
	var itemCount int

	for _, item := range items {
		subtotal += item.TotalPrice
		itemCount += item.Quantity
	}

	// Get applied coupons
	coupons, err := r.GetCartCoupons(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart coupons: %w", err)
	}

	var discountAmount float64
	for _, coupon := range coupons {
		discountAmount += coupon.DiscountAmount
	}

	// Get shipping
	shipping, err := r.GetCartShipping(ctx, cartID)
	var shippingAmount float64
	if err == nil && shipping != nil {
		shippingAmount = shipping.ShippingAmount
	}

	// Calculate tax (simplified - in real app, this would be more complex)
	taxAmount := subtotal * 0.1 // 10% tax rate

	// Calculate total
	totalAmount := subtotal + taxAmount + shippingAmount - discountAmount

	// Convert []*domain.CartItem to []domain.CartItem
	cartItems := make([]domain.CartItem, len(items))
	for i, item := range items {
		cartItems[i] = *item
	}

	summary := &domain.CartSummary{
		CartID:         cartID,
		ItemCount:      itemCount,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		ShippingAmount: shippingAmount,
		DiscountAmount: discountAmount,
		TotalAmount:    totalAmount,
		Currency:       cart.Currency,
		Items:          cartItems,
	}

	return summary, nil
}

// CalculateCartTotal calculates the total amount for a cart
func (r *cartRepository) CalculateCartTotal(ctx context.Context, cartID int64) (float64, error) {
	summary, err := r.GetCartSummary(ctx, cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart summary: %w", err)
	}

	return summary.TotalAmount, nil
}

// GetCartItemCount gets the total number of items in a cart
func (r *cartRepository) GetCartItemCount(ctx context.Context, cartID int64) (int, error) {
	query := `SELECT COALESCE(SUM(quantity), 0) FROM cart_items WHERE cart_id = $1`

	var count int
	err := r.db.GetContext(ctx, &count, query, cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart item count: %w", err)
	}

	return count, nil
}

// Cart Coupons

// ApplyCouponToCart applies a coupon to a cart
func (r *cartRepository) ApplyCouponToCart(ctx context.Context, cartCoupon *domain.CartCoupon) error {
	query := `
		INSERT INTO cart_coupons (cart_id, coupon_code, discount_amount, created_at)
		VALUES (:cart_id, :coupon_code, :discount_amount, :created_at)`

	cartCoupon.CreatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, cartCoupon)
	if err != nil {
		return fmt.Errorf("failed to apply coupon to cart: %w", err)
	}

	return nil
}

// RemoveCouponFromCart removes a coupon from a cart
func (r *cartRepository) RemoveCouponFromCart(ctx context.Context, cartID int64, couponCode string) error {
	query := `DELETE FROM cart_coupons WHERE cart_id = $1 AND coupon_code = $2`

	result, err := r.db.ExecContext(ctx, query, cartID, couponCode)
	if err != nil {
		return fmt.Errorf("failed to remove coupon from cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("coupon %s not found in cart", couponCode)
	}

	return nil
}

// GetCartCoupons retrieves all coupons applied to a cart
func (r *cartRepository) GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.CartCoupon, error) {
	query := `SELECT * FROM cart_coupons WHERE cart_id = $1 ORDER BY created_at ASC`

	var coupons []*domain.CartCoupon
	err := r.db.SelectContext(ctx, &coupons, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart coupons: %w", err)
	}

	return coupons, nil
}

// GetCartCouponByCode retrieves a specific coupon from a cart
func (r *cartRepository) GetCartCouponByCode(ctx context.Context, cartID int64, couponCode string) (*domain.CartCoupon, error) {
	query := `SELECT * FROM cart_coupons WHERE cart_id = $1 AND coupon_code = $2`

	var coupon domain.CartCoupon
	err := r.db.GetContext(ctx, &coupon, query, cartID, couponCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("coupon %s not found in cart", couponCode)
		}
		return nil, fmt.Errorf("failed to get cart coupon: %w", err)
	}

	return &coupon, nil
}

// Cart Shipping

// SetCartShipping sets shipping information for a cart
func (r *cartRepository) SetCartShipping(ctx context.Context, shipping *domain.CartShipping) error {
	query := `
		INSERT INTO cart_shipping (cart_id, shipping_method_id, shipping_method, shipping_amount, estimated_days, created_at)
		VALUES (:cart_id, :shipping_method_id, :shipping_method, :shipping_amount, :estimated_days, :created_at)`

	shipping.CreatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, shipping)
	if err != nil {
		return fmt.Errorf("failed to set cart shipping: %w", err)
	}

	return nil
}

// UpdateCartShipping updates shipping information for a cart
func (r *cartRepository) UpdateCartShipping(ctx context.Context, cartID int64, shipping *domain.CartShipping) error {
	query := `
		UPDATE cart_shipping SET
			shipping_method_id = :shipping_method_id, shipping_method = :shipping_method,
			shipping_amount = :shipping_amount, estimated_days = :estimated_days
		WHERE cart_id = :cart_id`

	shipping.CartID = cartID

	result, err := r.db.NamedExecContext(ctx, query, shipping)
	if err != nil {
		return fmt.Errorf("failed to update cart shipping: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("shipping not found for cart ID %d", cartID)
	}

	return nil
}

// GetCartShipping retrieves shipping information for a cart
func (r *cartRepository) GetCartShipping(ctx context.Context, cartID int64) (*domain.CartShipping, error) {
	query := `SELECT * FROM cart_shipping WHERE cart_id = $1`

	var shipping domain.CartShipping
	err := r.db.GetContext(ctx, &shipping, query, cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No shipping set
		}
		return nil, fmt.Errorf("failed to get cart shipping: %w", err)
	}

	return &shipping, nil
}

// DeleteCartShipping removes shipping information from a cart
func (r *cartRepository) DeleteCartShipping(ctx context.Context, cartID int64) error {
	query := `DELETE FROM cart_shipping WHERE cart_id = $1`

	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to delete cart shipping: %w", err)
	}

	return nil
}

// Cart Analytics & Management

// GetExpiredCarts retrieves carts that have expired
func (r *cartRepository) GetExpiredCarts(ctx context.Context, before time.Time) ([]*domain.Cart, error) {
	query := `SELECT * FROM carts WHERE expires_at < $1`

	var carts []*domain.Cart
	err := r.db.SelectContext(ctx, &carts, query, before)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired carts: %w", err)
	}

	return carts, nil
}

// DeleteExpiredCarts deletes expired carts
func (r *cartRepository) DeleteExpiredCarts(ctx context.Context, before time.Time) error {
	// Get expired cart IDs first
	expiredCarts, err := r.GetExpiredCarts(ctx, before)
	if err != nil {
		return fmt.Errorf("failed to get expired carts: %w", err)
	}

	// Delete each expired cart (this will cascade delete related data)
	for _, cart := range expiredCarts {
		err = r.DeleteCart(ctx, cart.ID)
		if err != nil {
			// Log error but continue with other carts
			fmt.Printf("Warning: failed to delete expired cart %d: %v\n", cart.ID, err)
		}
	}

	return nil
}

// GetCartAnalytics retrieves analytics data for carts
func (r *cartRepository) GetCartAnalytics(ctx context.Context) (*domain.CartAnalytics, error) {
	// This would be implemented with more complex queries in a real application
	// For now, we'll return a basic structure

	analytics := &domain.CartAnalytics{
		TotalCarts:          0,
		ActiveCarts:         0,
		AbandonedCarts:      0,
		AverageCartValue:    0.0,
		TotalCartValue:      0.0,
		ConversionRate:      0.0,
		AverageItemsPerCart: 0.0,
	}

	return analytics, nil
}

// MergeCarts merges items from source cart to target cart
func (r *cartRepository) MergeCarts(ctx context.Context, sourceCartID, targetCartID int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get source cart items
	sourceItems, err := r.GetCartItems(ctx, sourceCartID)
	if err != nil {
		return fmt.Errorf("failed to get source cart items: %w", err)
	}

	// Move items to target cart
	for _, item := range sourceItems {
		// Check if item already exists in target cart
		existingItem, err := r.GetCartItemByProduct(ctx, targetCartID, item.ProductID, item.ProductVariantID)
		if err == nil {
			// Item exists, update quantity
			existingItem.Quantity += item.Quantity
			existingItem.TotalPrice = existingItem.UnitPrice * float64(existingItem.Quantity)
			existingItem.UpdatedAt = time.Now()

			_, err = tx.NamedExecContext(ctx, `
				UPDATE cart_items SET quantity = :quantity, total_price = :total_price, updated_at = :updated_at
				WHERE id = :id`, existingItem)
			if err != nil {
				return fmt.Errorf("failed to update existing cart item: %w", err)
			}
		} else {
			// Item doesn't exist, add it
			item.CartID = targetCartID
			item.ID = 0 // Reset ID for new insert
			item.CreatedAt = time.Now()
			item.UpdatedAt = time.Now()

			_, err = tx.NamedExecContext(ctx, `
				INSERT INTO cart_items (cart_id, product_id, product_variant_id, quantity, unit_price, total_price, created_at, updated_at)
				VALUES (:cart_id, :product_id, :product_variant_id, :quantity, :unit_price, :total_price, :created_at, :updated_at)`, item)
			if err != nil {
				return fmt.Errorf("failed to add item to target cart: %w", err)
			}
		}
	}

	// Delete source cart
	_, err = tx.ExecContext(ctx, "DELETE FROM carts WHERE id = $1", sourceCartID)
	if err != nil {
		return fmt.Errorf("failed to delete source cart: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Wishlist Management

// CreateWishlist creates a new wishlist
func (r *cartRepository) CreateWishlist(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `
		INSERT INTO wishlists (user_id, name, is_public, created_at, updated_at)
		VALUES (:user_id, :name, :is_public, :created_at, :updated_at)`

	wishlist.CreatedAt = time.Now()
	wishlist.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, wishlist)
	if err != nil {
		return fmt.Errorf("failed to create wishlist: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get wishlist ID: %w", err)
	}

	wishlist.ID = id
	return nil
}

// GetWishlistByID retrieves a wishlist by ID
func (r *cartRepository) GetWishlistByID(ctx context.Context, id int64) (*domain.Wishlist, error) {
	query := `SELECT * FROM wishlists WHERE id = $1`

	var wishlist domain.Wishlist
	err := r.db.GetContext(ctx, &wishlist, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wishlist with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	return &wishlist, nil
}

// GetWishlistsByUserID retrieves wishlists for a user
func (r *cartRepository) GetWishlistsByUserID(ctx context.Context, userID int64, offset, limit int) ([]*domain.Wishlist, int64, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM wishlists WHERE user_id = $1`
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count wishlists: %w", err)
	}

	// List query
	query := `SELECT * FROM wishlists WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var wishlists []*domain.Wishlist
	err = r.db.SelectContext(ctx, &wishlists, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wishlists: %w", err)
	}

	return wishlists, total, nil
}

// UpdateWishlist updates an existing wishlist
func (r *cartRepository) UpdateWishlist(ctx context.Context, id int64, wishlist *domain.Wishlist) error {
	query := `
		UPDATE wishlists SET
			name = :name, is_public = :is_public, updated_at = :updated_at
		WHERE id = :id`

	wishlist.UpdatedAt = time.Now()
	wishlist.ID = id

	result, err := r.db.NamedExecContext(ctx, query, wishlist)
	if err != nil {
		return fmt.Errorf("failed to update wishlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("wishlist with ID %d not found", id)
	}

	return nil
}

// DeleteWishlist deletes a wishlist
func (r *cartRepository) DeleteWishlist(ctx context.Context, id int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete wishlist items
	_, err = tx.ExecContext(ctx, "DELETE FROM wishlist_items WHERE wishlist_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete wishlist items: %w", err)
	}

	// Delete wishlist
	result, err := tx.ExecContext(ctx, "DELETE FROM wishlists WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("wishlist with ID %d not found", id)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Wishlist Items

// AddItemToWishlist adds an item to a wishlist
func (r *cartRepository) AddItemToWishlist(ctx context.Context, item *domain.WishlistItem) error {
	query := `
		INSERT INTO wishlist_items (wishlist_id, product_id, product_variant_id, notes, created_at)
		VALUES (:wishlist_id, :product_id, :product_variant_id, :notes, :created_at)`

	item.CreatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to add item to wishlist: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get wishlist item ID: %w", err)
	}

	item.ID = id
	return nil
}

// GetWishlistItemByID retrieves a wishlist item by ID
func (r *cartRepository) GetWishlistItemByID(ctx context.Context, id int64) (*domain.WishlistItem, error) {
	query := `SELECT * FROM wishlist_items WHERE id = $1`

	var item domain.WishlistItem
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("wishlist item with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get wishlist item: %w", err)
	}

	return &item, nil
}

// GetWishlistItems retrieves items in a wishlist
func (r *cartRepository) GetWishlistItems(ctx context.Context, wishlistID int64, offset, limit int) ([]*domain.WishlistItem, int64, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM wishlist_items WHERE wishlist_id = $1`
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, wishlistID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count wishlist items: %w", err)
	}

	// List query
	query := `SELECT * FROM wishlist_items WHERE wishlist_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var items []*domain.WishlistItem
	err = r.db.SelectContext(ctx, &items, query, wishlistID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wishlist items: %w", err)
	}

	return items, total, nil
}

// UpdateWishlistItem updates an existing wishlist item
func (r *cartRepository) UpdateWishlistItem(ctx context.Context, id int64, item *domain.WishlistItem) error {
	query := `UPDATE wishlist_items SET notes = :notes WHERE id = :id`

	item.ID = id

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update wishlist item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("wishlist item with ID %d not found", id)
	}

	return nil
}

// DeleteWishlistItem deletes a wishlist item
func (r *cartRepository) DeleteWishlistItem(ctx context.Context, id int64) error {
	query := `DELETE FROM wishlist_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete wishlist item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("wishlist item with ID %d not found", id)
	}

	return nil
}

// MoveItemToCart moves an item from wishlist to cart
func (r *cartRepository) MoveItemToCart(ctx context.Context, wishlistItemID, cartID int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get wishlist item
	wishlistItem, err := r.GetWishlistItemByID(ctx, wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to get wishlist item: %w", err)
	}

	// Create cart item (you would need to get the current price from product)
	cartItem := &domain.CartItem{
		CartID:           cartID,
		ProductID:        wishlistItem.ProductID,
		ProductVariantID: wishlistItem.ProductVariantID,
		Quantity:         1, // Default quantity
		UnitPrice:        0, // This should be fetched from product
		TotalPrice:       0,
	}

	// Add to cart
	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO cart_items (cart_id, product_id, product_variant_id, quantity, unit_price, total_price, created_at, updated_at)
		VALUES (:cart_id, :product_id, :product_variant_id, :quantity, :unit_price, :total_price, :created_at, :updated_at)`, cartItem)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	// Remove from wishlist
	_, err = tx.ExecContext(ctx, "DELETE FROM wishlist_items WHERE id = $1", wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to remove item from wishlist: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
