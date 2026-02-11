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
	"github.com/jattinmanhas/GearboxV2/services/shared/middleware"
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
	GetCartAnalyticsByDateRange(ctx context.Context, startDate, endDate time.Time) (*dto.CartAnalyticsResponse, error)
	GetTopProductsInCarts(ctx context.Context, limit int) ([]*dto.ProductCartStatsResponse, error)
	GetCartAbandonmentRate(ctx context.Context) (float64, error)
	GetCartConversionFunnel(ctx context.Context) (*dto.CartConversionFunnelResponse, error)

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
	MoveItemToCartWithQuantity(ctx context.Context, wishlistItemID, cartID int64, quantity int) error
}

type cartService struct {
	cartRepo         repository.CartRepository
	productRepo      repository.ProductRepository
	inventoryService InventoryService
	couponService    CouponService
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository, inventoryService InventoryService, couponService CouponService) CartService {
	return &cartService{
		cartRepo:         cartRepo,
		productRepo:      productRepo,
		inventoryService: inventoryService,
		couponService:    couponService,
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

	// Check stock availability
	available, err := s.inventoryService.CheckStockAvailability(ctx, req.ProductID, req.ProductVariantID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to check stock availability: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("insufficient stock for requested quantity")
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

	// Check stock availability if quantity is being updated
	if req.Quantity != nil {
		available, err := s.inventoryService.CheckStockAvailability(ctx, existingItem.ProductID, existingItem.ProductVariantID, *req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to check stock availability: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("insufficient stock for requested quantity")
		}
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

	// Recalculate discount based on current subtotal to ensure maximum discount limit is respected
	// Get applied coupons and recalculate discount
	var recalculatedDiscount float64
	cartCoupons, err := s.couponService.GetCartCoupons(ctx, cartID)
	if err == nil && len(cartCoupons) > 0 {
		// Get cart coupon records to update their discount amounts
		cartCouponRecords, err := s.cartRepo.GetCartCoupons(ctx, cartID)
		if err == nil {
			// Recalculate discount for each coupon based on current subtotal
			for i, coupon := range cartCoupons {
				discountAmount, err := s.couponService.CalculateDiscount(ctx, coupon, summary.Subtotal)
				if err == nil {
					// Ensure discount doesn't exceed subtotal for this coupon
					if discountAmount > summary.Subtotal {
						discountAmount = summary.Subtotal
					}
					recalculatedDiscount += discountAmount
					// Update the stored discount amount in cart_coupons table
					if i < len(cartCouponRecords) {
						err = s.cartRepo.UpdateCartCouponDiscount(ctx, cartID, cartCouponRecords[i].CouponCode, discountAmount)
						if err != nil {
							// Log error but don't fail the operation
							fmt.Printf("Warning: failed to update cart coupon discount: %v\n", err)
						}
					}
				}
			}
		}
		// Ensure discount doesn't exceed subtotal
		if recalculatedDiscount > summary.Subtotal {
			recalculatedDiscount = summary.Subtotal
		}
		// Recalculate total with new discount
		summary.TotalAmount = summary.Subtotal - summary.DiscountAmount
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
	// Get user ID from context (extracted from JWT token by middleware)
	// For guest users, this will be 0, which is fine
	userIDValue := middleware.GetUserIDFromContext(ctx)
	var userID *int64
	if userIDValue != 0 {
		userIDVal := int64(userIDValue)
		userID = &userIDVal
	}

	// Use coupon service to apply coupon
	coupon, discountAmount, err := s.couponService.ApplyCouponToCart(ctx, cartID, req.CouponCode, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to apply coupon: %w", err)
	}

	// Return cart coupon representation
	cartCoupon := &domain.CartCoupon{
		CartID:         cartID,
		CouponCode:     coupon.Code,
		DiscountAmount: discountAmount,
	}

	return cartCoupon, nil
}

// RemoveCouponFromCart removes a coupon from a cart
func (s *cartService) RemoveCouponFromCart(ctx context.Context, cartID int64, req *dto.RemoveCouponRequest) error {
	// Use coupon service to remove coupon
	err := s.couponService.RemoveCouponFromCart(ctx, cartID, req.CouponCode)
	if err != nil {
		return fmt.Errorf("failed to remove coupon from cart: %w", err)
	}

	return nil
}

// GetCartCoupons retrieves all coupons applied to a cart
func (s *cartService) GetCartCoupons(ctx context.Context, cartID int64) ([]*domain.CartCoupon, error) {
	// Use coupon service to get cart coupons
	coupons, err := s.couponService.GetCartCoupons(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart coupons: %w", err)
	}

	// Convert domain coupons to cart coupons
	cartCoupons := make([]*domain.CartCoupon, len(coupons))
	for i, coupon := range coupons {
		// Get discount amount from cart coupon record
		cartCoupon, err := s.cartRepo.GetCartCouponByCode(ctx, cartID, coupon.Code)
		if err != nil {
			// Skip if cart coupon not found
			continue
		}
		cartCoupons[i] = cartCoupon
	}

	return cartCoupons, nil
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

	// Get additional metrics
	abandonmentRate, err := s.cartRepo.GetCartAbandonmentRate(ctx)
	if err != nil {
		abandonmentRate = 0 // Don't fail the whole request for this
	}

	// Get time-based metrics
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get new carts today
	todayAnalytics, err := s.cartRepo.GetCartAnalyticsByDateRange(ctx, today, now)
	if err != nil {
		todayAnalytics = &domain.CartAnalytics{}
	}

	// Get new carts this week
	weekAnalytics, err := s.cartRepo.GetCartAnalyticsByDateRange(ctx, weekStart, now)
	if err != nil {
		weekAnalytics = &domain.CartAnalytics{}
	}

	// Get new carts this month
	monthAnalytics, err := s.cartRepo.GetCartAnalyticsByDateRange(ctx, monthStart, now)
	if err != nil {
		monthAnalytics = &domain.CartAnalytics{}
	}

	return &dto.CartAnalyticsResponse{
		TotalCarts:          analytics.TotalCarts,
		ActiveCarts:         analytics.ActiveCarts,
		AbandonedCarts:      analytics.AbandonedCarts,
		AverageCartValue:    analytics.AverageCartValue,
		TotalCartValue:      analytics.TotalCartValue,
		ConversionRate:      analytics.ConversionRate,
		AverageItemsPerCart: analytics.AverageItemsPerCart,
		AbandonmentRate:     abandonmentRate,
		NewCartsToday:       todayAnalytics.TotalCarts,
		NewCartsThisWeek:    weekAnalytics.TotalCarts,
		NewCartsThisMonth:   monthAnalytics.TotalCarts,
	}, nil
}

// GetCartAnalyticsByDateRange retrieves analytics data for carts within a date range
func (s *cartService) GetCartAnalyticsByDateRange(ctx context.Context, startDate, endDate time.Time) (*dto.CartAnalyticsResponse, error) {
	analytics, err := s.cartRepo.GetCartAnalyticsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart analytics by date range: %w", err)
	}

	return &dto.CartAnalyticsResponse{
		TotalCarts:          analytics.TotalCarts,
		ActiveCarts:         analytics.ActiveCarts,
		AbandonedCarts:      analytics.AbandonedCarts,
		AverageCartValue:    analytics.AverageCartValue,
		TotalCartValue:      analytics.TotalCartValue,
		ConversionRate:      analytics.ConversionRate,
		AverageItemsPerCart: analytics.AverageItemsPerCart,
		AbandonmentRate:     analytics.AbandonmentRate,
		NewCartsToday:       0, // Not applicable for date range
		NewCartsThisWeek:    0, // Not applicable for date range
		NewCartsThisMonth:   0, // Not applicable for date range
	}, nil
}

// GetTopProductsInCarts retrieves top products by quantity/value in carts
func (s *cartService) GetTopProductsInCarts(ctx context.Context, limit int) ([]*dto.ProductCartStatsResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	stats, err := s.cartRepo.GetTopProductsInCarts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top products in carts: %w", err)
	}

	responses := make([]*dto.ProductCartStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = &dto.ProductCartStatsResponse{
			ProductID:     stat.ProductID,
			ProductName:   stat.ProductName,
			SKU:           stat.SKU,
			TotalQuantity: stat.TotalQuantity,
			TotalValue:    stat.TotalValue,
			CartCount:     stat.CartCount,
			AveragePrice:  stat.AveragePrice,
		}
	}

	return responses, nil
}

// GetCartAbandonmentRate calculates the cart abandonment rate
func (s *cartService) GetCartAbandonmentRate(ctx context.Context) (float64, error) {
	rate, err := s.cartRepo.GetCartAbandonmentRate(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get cart abandonment rate: %w", err)
	}

	return rate, nil
}

// GetCartConversionFunnel retrieves conversion funnel data
func (s *cartService) GetCartConversionFunnel(ctx context.Context) (*dto.CartConversionFunnelResponse, error) {
	funnel, err := s.cartRepo.GetCartConversionFunnel(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart conversion funnel: %w", err)
	}

	return &dto.CartConversionFunnelResponse{
		Visitors:            funnel.Visitors,
		CartCreated:         funnel.CartCreated,
		ItemsAdded:          funnel.ItemsAdded,
		CheckoutStarted:     funnel.CheckoutStarted,
		OrderCompleted:      funnel.OrderCompleted,
		CartToCheckoutRate:  funnel.CartToCheckoutRate,
		CheckoutToOrderRate: funnel.CheckoutToOrderRate,
		OverallConversion:   funnel.OverallConversion,
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

// GetWishlistItems retrieves items in a wishlist with product details
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

	// Convert to response DTOs with product details
	itemResponses := make([]dto.WishlistItemResponse, len(items))
	for i, item := range items {
		// Get product details for each item
		product, err := s.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			// If product not found, create a minimal product response
			product = &domain.Product{
				ID:    item.ProductID,
				Name:  "Product Not Found",
				SKU:   "N/A",
				Price: 0,
			}
		}

		// Convert product to ProductResponse
		productResponse := dto.ProductResponse{
			ID:               product.ID,
			Name:             product.Name,
			Description:      product.Description,
			ShortDesc:        product.ShortDesc,
			SKU:              product.SKU,
			Price:            product.Price,
			ComparePrice:     product.ComparePrice,
			CostPrice:        product.CostPrice,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsDigital:        product.IsDigital,
			RequiresShipping: product.RequiresShipping,
			Taxable:          product.Taxable,
			TrackQuantity:    product.TrackQuantity,
			MinQuantity:      product.MinQuantity,
			MaxQuantity:      product.MaxQuantity,
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CategoryIDs:      product.CategoryIDs,
			CategoryNames:    product.CategoryNames,
			CreatedAt:        product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        product.UpdatedAt.Format(time.RFC3339),
		}

		itemResponses[i] = dto.WishlistItemResponse{
			ID:               item.ID,
			WishlistID:       item.WishlistID,
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Notes:            item.Notes,
			CreatedAt:        item.CreatedAt.Format(time.RFC3339),
			Product:          productResponse,
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
	// Get wishlist item
	wishlistItem, err := s.cartRepo.GetWishlistItemByID(ctx, wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to get wishlist item: %w", err)
	}

	// Check if cart exists
	_, err = s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Check stock availability (default quantity of 1 for wishlist items)
	available, err := s.inventoryService.CheckStockAvailability(ctx, wishlistItem.ProductID, wishlistItem.ProductVariantID, 1)
	if err != nil {
		return fmt.Errorf("failed to check stock availability: %w", err)
	}
	if !available {
		return fmt.Errorf("insufficient stock for this item")
	}

	// Get current product price
	unitPrice, err := s.getCurrentProductPrice(ctx, wishlistItem.ProductID, wishlistItem.ProductVariantID)
	if err != nil {
		return fmt.Errorf("failed to get current product price: %w", err)
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetCartItemByProduct(ctx, cartID, wishlistItem.ProductID, wishlistItem.ProductVariantID)
	if err == nil {
		// Item exists, update quantity and price
		existingItem.Quantity += 1
		existingItem.UnitPrice = unitPrice // Use current price
		existingItem.TotalPrice = existingItem.UnitPrice * float64(existingItem.Quantity)
		existingItem.UpdatedAt = time.Now()

		err = s.cartRepo.UpdateCartItem(ctx, existingItem.ID, existingItem)
		if err != nil {
			return fmt.Errorf("failed to update existing cart item: %w", err)
		}
	} else {
		// Create new cart item with current price
		cartItem := &domain.CartItem{
			CartID:           cartID,
			ProductID:        wishlistItem.ProductID,
			ProductVariantID: wishlistItem.ProductVariantID,
			Quantity:         1,
			UnitPrice:        unitPrice, // Use current product price
			TotalPrice:       unitPrice, // 1 * unitPrice
		}

		err = s.cartRepo.AddItemToCart(ctx, cartItem)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}
	}

	// Remove from wishlist
	err = s.cartRepo.DeleteWishlistItem(ctx, wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to remove item from wishlist: %w", err)
	}

	return nil
}

// MoveItemToCartWithQuantity moves an item from wishlist to cart with specified quantity
func (s *cartService) MoveItemToCartWithQuantity(ctx context.Context, wishlistItemID, cartID int64, quantity int) error {
	// Validate quantity
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}

	// Get wishlist item
	wishlistItem, err := s.cartRepo.GetWishlistItemByID(ctx, wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to get wishlist item: %w", err)
	}

	// Check if cart exists
	_, err = s.cartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Check stock availability
	available, err := s.inventoryService.CheckStockAvailability(ctx, wishlistItem.ProductID, wishlistItem.ProductVariantID, quantity)
	if err != nil {
		return fmt.Errorf("failed to check stock availability: %w", err)
	}
	if !available {
		return fmt.Errorf("insufficient stock for requested quantity")
	}

	// Get current product price
	unitPrice, err := s.getCurrentProductPrice(ctx, wishlistItem.ProductID, wishlistItem.ProductVariantID)
	if err != nil {
		return fmt.Errorf("failed to get current product price: %w", err)
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetCartItemByProduct(ctx, cartID, wishlistItem.ProductID, wishlistItem.ProductVariantID)
	if err == nil {
		// Item exists, update quantity and price
		existingItem.Quantity += quantity
		existingItem.UnitPrice = unitPrice // Use current price
		existingItem.TotalPrice = existingItem.UnitPrice * float64(existingItem.Quantity)
		existingItem.UpdatedAt = time.Now()

		err = s.cartRepo.UpdateCartItem(ctx, existingItem.ID, existingItem)
		if err != nil {
			return fmt.Errorf("failed to update existing cart item: %w", err)
		}
	} else {
		// Create new cart item with current price
		cartItem := &domain.CartItem{
			CartID:           cartID,
			ProductID:        wishlistItem.ProductID,
			ProductVariantID: wishlistItem.ProductVariantID,
			Quantity:         quantity,
			UnitPrice:        unitPrice, // Use current product price
			TotalPrice:       unitPrice * float64(quantity),
		}

		err = s.cartRepo.AddItemToCart(ctx, cartItem)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}
	}

	// Remove from wishlist
	err = s.cartRepo.DeleteWishlistItem(ctx, wishlistItemID)
	if err != nil {
		return fmt.Errorf("failed to remove item from wishlist: %w", err)
	}

	return nil
}
