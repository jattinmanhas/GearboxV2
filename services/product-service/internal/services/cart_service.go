package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type CartService interface {
	// Cart Management
	CreateCart(ctx context.Context, req *dto.CreateCartRequest) (*domain.Cart, error)
	GetCartByID(ctx context.Context, id int64) (*domain.Cart, error)
	GetCartByUserID(ctx context.Context, userID int64) (*domain.Cart, error)
	GetCartBySessionID(ctx context.Context, sessionID string) (*domain.Cart, error)
	UpdateCart(ctx context.Context, id int64, req *dto.UpdateCartRequest) (*domain.Cart, error)
	DeleteCart(ctx context.Context, id int64) error
	GetOrCreateCart(ctx context.Context, userID *int64, sessionID string, currency string) (*domain.Cart, error)

	// Cart Items
	AddItemToCart(ctx context.Context, cartID int64, req *dto.AddToCartRequest) (*domain.CartItem, error)
	GetCartItemByID(ctx context.Context, id int64) (*domain.CartItem, error)
	UpdateCartItem(ctx context.Context, id int64, req *dto.UpdateCartItemRequest) (*domain.CartItem, error)
	DeleteCartItem(ctx context.Context, id int64) error
	GetCartItems(ctx context.Context, cartID int64) ([]*domain.CartItem, error)
	ClearCartItems(ctx context.Context, cartID int64) error

	// Cart Summary & Calculations
	GetCartSummary(ctx context.Context, cartID int64) (*dto.CartSummaryResponse, error)
	CalculateCartTotal(ctx context.Context, cartID int64) (float64, error)
	GetCartItemCount(ctx context.Context, cartID int64) (int, error)

	// Cart Coupons
	ApplyCouponToCart(ctx context.Context, cartID int64, req *dto.ApplyCouponRequest) (*domain.CartCoupon, error)
	RemoveCouponFromCart(ctx context.Context, cartID int64, req *dto.RemoveCouponRequest) error
	GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.CartCoupon, error)

	// Cart Shipping
	SetCartShipping(ctx context.Context, cartID int64, req *dto.SetShippingRequest) (*domain.CartShipping, error)
	UpdateCartShipping(ctx context.Context, cartID int64, req *dto.UpdateShippingRequest) (*domain.CartShipping, error)
	GetCartShipping(ctx context.Context, cartID int64) (*domain.CartShipping, error)
	DeleteCartShipping(ctx context.Context, cartID int64) error

	// Cart Operations
	MergeCarts(ctx context.Context, sourceCartID, targetCartID int64) error
	ClearCart(ctx context.Context, cartID int64) error
	GetCartAnalytics(ctx context.Context) (*dto.CartAnalyticsResponse, error)

	// Wishlist Management
	CreateWishlist(ctx context.Context, userID int64, req *dto.CreateWishlistRequest) (*domain.Wishlist, error)
	GetWishlistByID(ctx context.Context, id int64) (*domain.Wishlist, error)
	GetWishlistsByUserID(ctx context.Context, userID int64, page, limit int) (*dto.ListWishlistsResponse, error)
	UpdateWishlist(ctx context.Context, id int64, req *dto.UpdateWishlistRequest) (*domain.Wishlist, error)
	DeleteWishlist(ctx context.Context, id int64) error

	// Wishlist Items
	AddItemToWishlist(ctx context.Context, wishlistID int64, req *dto.AddToWishlistRequest) (*domain.WishlistItem, error)
	GetWishlistItemByID(ctx context.Context, id int64) (*domain.WishlistItem, error)
	GetWishlistItems(ctx context.Context, wishlistID int64, page, limit int) (*dto.ListWishlistItemsResponse, error)
	UpdateWishlistItem(ctx context.Context, id int64, req *dto.UpdateWishlistItemRequest) (*domain.WishlistItem, error)
	DeleteWishlistItem(ctx context.Context, id int64) error
	MoveItemToCart(ctx context.Context, wishlistItemID, cartID int64) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// Cart Management

// CreateCart creates a new cart
func (s *cartService) CreateCart(ctx context.Context, req *dto.CreateCartRequest) (*domain.Cart, error) {
	// Set expiration time (30 days from now)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	cart := &domain.Cart{
		UserID:    req.UserID,
		SessionID: req.SessionID,
		Currency:  req.Currency,
		ExpiresAt: &expiresAt,
	}

	err := s.cartRepo.CreateCart(ctx, cart)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return cart, nil
}

// GetCartByID retrieves a cart by ID
func (s *cartService) GetCartByID(ctx context.Context, id int64) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCartByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart, nil
}

// GetCartByUserID retrieves a cart by user ID
func (s *cartService) GetCartByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart, nil
}

// GetCartBySessionID retrieves a cart by session ID
func (s *cartService) GetCartBySessionID(ctx context.Context, sessionID string) (*domain.Cart, error) {
	cart, err := s.cartRepo.GetCartBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	return cart, nil
}

// UpdateCart updates an existing cart
func (s *cartService) UpdateCart(ctx context.Context, id int64, req *dto.UpdateCartRequest) (*domain.Cart, error) {
	// Get existing cart
	existingCart, err := s.cartRepo.GetCartByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing cart: %w", err)
	}

	// Update fields that are provided
	updateCart := *existingCart

	if req.Currency != nil {
		updateCart.Currency = *req.Currency
	}

	updateCart.UpdatedAt = time.Now()

	// Update cart in repository
	err = s.cartRepo.UpdateCart(ctx, &updateCart)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return &updateCart, nil
}

// DeleteCart deletes a cart
func (s *cartService) DeleteCart(ctx context.Context, id int64) error {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Delete cart
	err = s.cartRepo.DeleteCart(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	return nil
}

// GetOrCreateCart gets an existing cart or creates a new one
func (s *cartService) GetOrCreateCart(ctx context.Context, userID *int64, sessionID string, currency string) (*domain.Cart, error) {
	// Check if cart already exists for this session/user combination
	existingCart, err := s.cartRepo.GetCartBySessionOrUser(ctx, sessionID, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing cart: %w", err)
	}

	// If cart exists and is not expired, return existing cart
	if existingCart != nil && (existingCart.ExpiresAt == nil || existingCart.ExpiresAt.After(time.Now())) {
		// Update currency if different
		if existingCart.Currency != currency {
			existingCart.Currency = currency
			existingCart.UpdatedAt = time.Now()
			if err := s.cartRepo.UpdateCart(ctx, existingCart); err != nil {
				return nil, fmt.Errorf("failed to update cart currency: %w", err)
			}
		}
		return existingCart, nil
	}

	// Create new cart
	req := &dto.CreateCartRequest{
		UserID:    userID,
		SessionID: sessionID,
		Currency:  currency,
	}

	cart, err := s.CreateCart(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return cart, nil
}

// Cart Items

// AddItemToCart adds an item to the cart
func (s *cartService) AddItemToCart(ctx context.Context, cartID int64, req *dto.AddToCartRequest) (*domain.CartItem, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get current product price
	unitPrice, err := s.getCurrentProductPrice(ctx, req.ProductID, req.ProductVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product price: %w", err)
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetCartItemByProduct(ctx, cartID, req.ProductID, req.ProductVariantID)
	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += req.Quantity
		existingItem.UnitPrice = unitPrice // Use current product price
		existingItem.TotalPrice = existingItem.UnitPrice * float64(existingItem.Quantity)
		existingItem.UpdatedAt = time.Now()

		err = s.cartRepo.UpdateCartItem(ctx, existingItem.ID, existingItem)
		if err != nil {
			return nil, fmt.Errorf("failed to update cart item: %w", err)
		}

		return existingItem, nil
	}

	// Create new cart item
	cartItem := &domain.CartItem{
		CartID:           cartID,
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		Quantity:         req.Quantity,
		UnitPrice:        unitPrice, // Use current product price
		TotalPrice:       unitPrice * float64(req.Quantity),
	}

	err = s.cartRepo.AddItemToCart(ctx, cartItem)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	return cartItem, nil
}

// GetCartItemByID retrieves a cart item by ID
func (s *cartService) GetCartItemByID(ctx context.Context, id int64) (*domain.CartItem, error) {
	item, err := s.cartRepo.GetCartItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	return item, nil
}

// getCurrentProductPrice gets the current price for a product and variant
func (s *cartService) getCurrentProductPrice(ctx context.Context, productID int64, variantID *int64) (float64, error) {
	if variantID != nil {
		variant, err := s.productRepo.GetProductVariantByID(ctx, *variantID)
		if err != nil {
			return 0, fmt.Errorf("failed to get product variant: %w", err)
		}
		return variant.Price, nil
	}

	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return 0, fmt.Errorf("failed to get product: %w", err)
	}
	return product.Price, nil
}

// UpdateCartItem updates an existing cart item
func (s *cartService) UpdateCartItem(ctx context.Context, id int64, req *dto.UpdateCartItemRequest) (*domain.CartItem, error) {
	// Get existing item
	existingItem, err := s.cartRepo.GetCartItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing cart item: %w", err)
	}

	// Get current product price
	unitPrice, err := s.getCurrentProductPrice(ctx, existingItem.ProductID, existingItem.ProductVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product price: %w", err)
	}

	// Update fields that are provided
	updateItem := *existingItem

	if req.Quantity != nil {
		updateItem.Quantity = *req.Quantity
	}

	// Always use current product price
	updateItem.UnitPrice = unitPrice
	updateItem.TotalPrice = updateItem.UnitPrice * float64(updateItem.Quantity)
	updateItem.UpdatedAt = time.Now()

	// Update item in repository
	err = s.cartRepo.UpdateCartItem(ctx, id, &updateItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	return &updateItem, nil
}

// DeleteCartItem deletes a cart item
func (s *cartService) DeleteCartItem(ctx context.Context, id int64) error {
	// Check if item exists
	_, err := s.cartRepo.GetCartItemByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get cart item: %w", err)
	}

	// Delete item
	err = s.cartRepo.DeleteCartItem(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	return nil
}

// GetCartItems retrieves all items in a cart
func (s *cartService) GetCartItems(ctx context.Context, cartID int64) ([]*domain.CartItem, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	items, err := s.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	return items, nil
}

// ClearCartItems removes all items from a cart
func (s *cartService) ClearCartItems(ctx context.Context, cartID int64) error {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	err = s.cartRepo.ClearCartItems(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to clear cart items: %w", err)
	}

	return nil
}

// Cart Summary & Calculations

// GetCartSummary retrieves a complete cart summary
func (s *cartService) GetCartSummary(ctx context.Context, cartID int64) (*dto.CartSummaryResponse, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get cart summary from repository
	summary, err := s.cartRepo.GetCartSummary(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart summary: %w", err)
	}

	// Convert to response DTO
	itemResponses := make([]dto.CartItemResponse, len(summary.Items))
	for i, item := range summary.Items {
		itemResponses[i] = dto.CartItemResponse{
			ID:               item.ID,
			CartID:           item.CartID,
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			TotalPrice:       item.TotalPrice,
			CreatedAt:        item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        item.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &dto.CartSummaryResponse{
		CartID:         summary.CartID,
		ItemCount:      summary.ItemCount,
		Subtotal:       summary.Subtotal,
		TaxAmount:      summary.TaxAmount,
		ShippingAmount: summary.ShippingAmount,
		DiscountAmount: summary.DiscountAmount,
		TotalAmount:    summary.TotalAmount,
		Currency:       summary.Currency,
		Items:          itemResponses,
	}, nil
}

// CalculateCartTotal calculates the total amount for a cart
func (s *cartService) CalculateCartTotal(ctx context.Context, cartID int64) (float64, error) {
	total, err := s.cartRepo.CalculateCartTotal(ctx, cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate cart total: %w", err)
	}

	return total, nil
}

// GetCartItemCount gets the total number of items in a cart
func (s *cartService) GetCartItemCount(ctx context.Context, cartID int64) (int, error) {
	count, err := s.cartRepo.GetCartItemCount(ctx, cartID)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart item count: %w", err)
	}

	return count, nil
}

// Cart Coupons

// ApplyCouponToCart applies a coupon to a cart
func (s *cartService) ApplyCouponToCart(ctx context.Context, cartID int64, req *dto.ApplyCouponRequest) (*domain.CartCoupon, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Check if coupon is already applied
	_, err = s.cartRepo.GetCartCouponByCode(ctx, cartID, req.CouponCode)
	if err == nil {
		return nil, fmt.Errorf("coupon %s is already applied to this cart", req.CouponCode)
	}

	// In a real application, you would validate the coupon here
	// For now, we'll create a simple discount
	discountAmount := 10.0 // $10 discount

	cartCoupon := &domain.CartCoupon{
		CartID:         cartID,
		CouponCode:     req.CouponCode,
		DiscountAmount: discountAmount,
	}

	err = s.cartRepo.ApplyCouponToCart(ctx, cartCoupon)
	if err != nil {
		return nil, fmt.Errorf("failed to apply coupon to cart: %w", err)
	}

	return cartCoupon, nil
}

// RemoveCouponFromCart removes a coupon from a cart
func (s *cartService) RemoveCouponFromCart(ctx context.Context, cartID int64, req *dto.RemoveCouponRequest) error {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	err = s.cartRepo.RemoveCouponFromCart(ctx, cartID, req.CouponCode)
	if err != nil {
		return fmt.Errorf("failed to remove coupon from cart: %w", err)
	}

	return nil
}

// GetCartCoupons retrieves all coupons applied to a cart
func (s *cartService) GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.CartCoupon, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	coupons, err := s.cartRepo.GetCartCoupons(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart coupons: %w", err)
	}

	return coupons, nil
}

// Cart Shipping

// SetCartShipping sets shipping information for a cart
func (s *cartService) SetCartShipping(ctx context.Context, cartID int64, req *dto.SetShippingRequest) (*domain.CartShipping, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	shipping := &domain.CartShipping{
		CartID:           cartID,
		ShippingMethodID: req.ShippingMethodID,
		ShippingMethod:   req.ShippingMethod,
		ShippingAmount:   req.ShippingAmount,
		EstimatedDays:    req.EstimatedDays,
	}

	err = s.cartRepo.SetCartShipping(ctx, shipping)
	if err != nil {
		return nil, fmt.Errorf("failed to set cart shipping: %w", err)
	}

	return shipping, nil
}

// UpdateCartShipping updates shipping information for a cart
func (s *cartService) UpdateCartShipping(ctx context.Context, cartID int64, req *dto.UpdateShippingRequest) (*domain.CartShipping, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get existing shipping
	existingShipping, err := s.cartRepo.GetCartShipping(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing shipping: %w", err)
	}

	if existingShipping == nil {
		return nil, fmt.Errorf("no shipping information found for cart")
	}

	// Update fields that are provided
	updateShipping := *existingShipping

	if req.ShippingMethodID != nil {
		updateShipping.ShippingMethodID = *req.ShippingMethodID
	}
	if req.ShippingMethod != nil {
		updateShipping.ShippingMethod = *req.ShippingMethod
	}
	if req.ShippingAmount != nil {
		updateShipping.ShippingAmount = *req.ShippingAmount
	}
	if req.EstimatedDays != nil {
		updateShipping.EstimatedDays = *req.EstimatedDays
	}

	err = s.cartRepo.UpdateCartShipping(ctx, cartID, &updateShipping)
	if err != nil {
		return nil, fmt.Errorf("failed to update cart shipping: %w", err)
	}

	return &updateShipping, nil
}

// GetCartShipping retrieves shipping information for a cart
func (s *cartService) GetCartShipping(ctx context.Context, cartID int64) (*domain.CartShipping, error) {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	shipping, err := s.cartRepo.GetCartShipping(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart shipping: %w", err)
	}

	return shipping, nil
}

// DeleteCartShipping removes shipping information from a cart
func (s *cartService) DeleteCartShipping(ctx context.Context, cartID int64) error {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	err = s.cartRepo.DeleteCartShipping(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to delete cart shipping: %w", err)
	}

	return nil
}

// Cart Operations

// MergeCarts merges items from source cart to target cart
func (s *cartService) MergeCarts(ctx context.Context, sourceCartID, targetCartID int64) error {
	// Check if both carts exist
	_, err := s.cartRepo.GetCartByID(ctx, sourceCartID)
	if err != nil {
		return fmt.Errorf("failed to get source cart: %w", err)
	}

	_, err = s.cartRepo.GetCartByID(ctx, targetCartID)
	if err != nil {
		return fmt.Errorf("failed to get target cart: %w", err)
	}

	err = s.cartRepo.MergeCarts(ctx, sourceCartID, targetCartID)
	if err != nil {
		return fmt.Errorf("failed to merge carts: %w", err)
	}

	return nil
}

// ClearCart clears all items from a cart
func (s *cartService) ClearCart(ctx context.Context, cartID int64) error {
	// Check if cart exists
	_, err := s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	err = s.cartRepo.ClearCartItems(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// GetCartAnalytics retrieves analytics data for carts
func (s *cartService) GetCartAnalytics(ctx context.Context) (*dto.CartAnalyticsResponse, error) {
	analytics, err := s.cartRepo.GetCartAnalytics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart analytics: %w", err)
	}

	return &dto.CartAnalyticsResponse{
		TotalCarts:          analytics.TotalCarts,
		ActiveCarts:         analytics.ActiveCarts,
		AbandonedCarts:      analytics.AbandonedCarts,
		AverageCartValue:    analytics.AverageCartValue,
		TotalCartValue:      analytics.TotalCartValue,
		ConversionRate:      analytics.ConversionRate,
		AverageItemsPerCart: analytics.AverageItemsPerCart,
	}, nil
}

// Wishlist Management

// CreateWishlist creates a new wishlist
func (s *cartService) CreateWishlist(ctx context.Context, userID int64, req *dto.CreateWishlistRequest) (*domain.Wishlist, error) {
	wishlist := &domain.Wishlist{
		UserID:   userID,
		Name:     req.Name,
		IsPublic: req.IsPublic,
	}

	err := s.cartRepo.CreateWishlist(ctx, wishlist)
	if err != nil {
		return nil, fmt.Errorf("failed to create wishlist: %w", err)
	}

	return wishlist, nil
}

// GetWishlistByID retrieves a wishlist by ID
func (s *cartService) GetWishlistByID(ctx context.Context, id int64) (*domain.Wishlist, error) {
	wishlist, err := s.cartRepo.GetWishlistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	return wishlist, nil
}

// GetWishlistsByUserID retrieves wishlists for a user
func (s *cartService) GetWishlistsByUserID(ctx context.Context, userID int64, page, limit int) (*dto.ListWishlistsResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	wishlists, total, err := s.cartRepo.GetWishlistsByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlists: %w", err)
	}

	// Convert to response DTOs
	wishlistResponses := make([]dto.WishlistResponse, len(wishlists))
	for i, wishlist := range wishlists {
		wishlistResponses[i] = dto.WishlistResponse{
			ID:        wishlist.ID,
			UserID:    wishlist.UserID,
			Name:      wishlist.Name,
			IsPublic:  wishlist.IsPublic,
			CreatedAt: wishlist.CreatedAt.Format(time.RFC3339),
			UpdatedAt: wishlist.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.ListWishlistsResponse{
		Wishlists:  wishlistResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateWishlist updates an existing wishlist
func (s *cartService) UpdateWishlist(ctx context.Context, id int64, req *dto.UpdateWishlistRequest) (*domain.Wishlist, error) {
	// Get existing wishlist
	existingWishlist, err := s.cartRepo.GetWishlistByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing wishlist: %w", err)
	}

	// Update fields that are provided
	updateWishlist := *existingWishlist

	if req.Name != nil {
		updateWishlist.Name = *req.Name
	}
	if req.IsPublic != nil {
		updateWishlist.IsPublic = *req.IsPublic
	}

	updateWishlist.UpdatedAt = time.Now()

	// Update wishlist in repository
	err = s.cartRepo.UpdateWishlist(ctx, id, &updateWishlist)
	if err != nil {
		return nil, fmt.Errorf("failed to update wishlist: %w", err)
	}

	return &updateWishlist, nil
}

// DeleteWishlist deletes a wishlist
func (s *cartService) DeleteWishlist(ctx context.Context, id int64) error {
	// Check if wishlist exists
	_, err := s.cartRepo.GetWishlistByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get wishlist: %w", err)
	}

	err = s.cartRepo.DeleteWishlist(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete wishlist: %w", err)
	}

	return nil
}

// Wishlist Items

// AddItemToWishlist adds an item to a wishlist
func (s *cartService) AddItemToWishlist(ctx context.Context, wishlistID int64, req *dto.AddToWishlistRequest) (*domain.WishlistItem, error) {
	// Check if wishlist exists
	_, err := s.cartRepo.GetWishlistByID(ctx, wishlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	// Check if product exists
	_, err = s.productRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Check if variant exists (if provided)
	if req.ProductVariantID != nil {
		_, err = s.productRepo.GetProductVariantByID(ctx, *req.ProductVariantID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product variant: %w", err)
		}
	}

	wishlistItem := &domain.WishlistItem{
		WishlistID:       wishlistID,
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		Notes:            req.Notes,
	}

	err = s.cartRepo.AddItemToWishlist(ctx, wishlistItem)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to wishlist: %w", err)
	}

	return wishlistItem, nil
}

// GetWishlistItemByID retrieves a wishlist item by ID
func (s *cartService) GetWishlistItemByID(ctx context.Context, id int64) (*domain.WishlistItem, error) {
	item, err := s.cartRepo.GetWishlistItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist item: %w", err)
	}

	return item, nil
}

// GetWishlistItems retrieves items in a wishlist
func (s *cartService) GetWishlistItems(ctx context.Context, wishlistID int64, page, limit int) (*dto.ListWishlistItemsResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	items, total, err := s.cartRepo.GetWishlistItems(ctx, wishlistID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist items: %w", err)
	}

	// Convert to response DTOs
	itemResponses := make([]dto.WishlistItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.WishlistItemResponse{
			ID:               item.ID,
			WishlistID:       item.WishlistID,
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Notes:            item.Notes,
			CreatedAt:        item.CreatedAt.Format(time.RFC3339),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.ListWishlistItemsResponse{
		Items:      itemResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateWishlistItem updates an existing wishlist item
func (s *cartService) UpdateWishlistItem(ctx context.Context, id int64, req *dto.UpdateWishlistItemRequest) (*domain.WishlistItem, error) {
	// Get existing item
	existingItem, err := s.cartRepo.GetWishlistItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing wishlist item: %w", err)
	}

	// Update fields that are provided
	updateItem := *existingItem

	if req.Notes != nil {
		updateItem.Notes = *req.Notes
	}

	// Update item in repository
	err = s.cartRepo.UpdateWishlistItem(ctx, id, &updateItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update wishlist item: %w", err)
	}

	return &updateItem, nil
}

// DeleteWishlistItem deletes a wishlist item
func (s *cartService) DeleteWishlistItem(ctx context.Context, id int64) error {
	// Check if item exists
	_, err := s.cartRepo.GetWishlistItemByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get wishlist item: %w", err)
	}

	err = s.cartRepo.DeleteWishlistItem(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete wishlist item: %w", err)
	}

	return nil
}

// MoveItemToCart moves an item from wishlist to cart
func (s *cartService) MoveItemToCart(ctx context.Context, wishlistItemID, cartID int64) error {
	// Check if wishlist item exists
	_, err := s.cartRepo.GetWishlistItemByID(ctx, wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to get wishlist item: %w", err)
	}

	// Check if cart exists
	_, err = s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	err = s.cartRepo.MoveItemToCart(ctx, wishlistItemID, cartID)
	if err != nil {
		return fmt.Errorf("failed to move item to cart: %w", err)
	}

	return nil
}
