package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

type ICartHandler interface {
	// Cart Management
	CreateCart(w http.ResponseWriter, r *http.Request)
	GetCart(w http.ResponseWriter, r *http.Request)
	UpdateCart(w http.ResponseWriter, r *http.Request)
	DeleteCart(w http.ResponseWriter, r *http.Request)
	GetOrCreateCart(w http.ResponseWriter, r *http.Request)

	// Cart Items
	AddItemToCart(w http.ResponseWriter, r *http.Request)
	GetCartItem(w http.ResponseWriter, r *http.Request)
	UpdateCartItem(w http.ResponseWriter, r *http.Request)
	DeleteCartItem(w http.ResponseWriter, r *http.Request)
	GetCartItems(w http.ResponseWriter, r *http.Request)
	ClearCartItems(w http.ResponseWriter, r *http.Request)

	// Cart Summary & Calculations
	GetCartSummary(w http.ResponseWriter, r *http.Request)
	GetCartTotal(w http.ResponseWriter, r *http.Request)
	GetCartItemCount(w http.ResponseWriter, r *http.Request)

	// Cart Coupons
	ApplyCouponToCart(w http.ResponseWriter, r *http.Request)
	RemoveCouponFromCart(w http.ResponseWriter, r *http.Request)
	GetCartCoupons(w http.ResponseWriter, r *http.Request)

	// Cart Shipping
	SetCartShipping(w http.ResponseWriter, r *http.Request)
	UpdateCartShipping(w http.ResponseWriter, r *http.Request)
	GetCartShipping(w http.ResponseWriter, r *http.Request)
	DeleteCartShipping(w http.ResponseWriter, r *http.Request)

	// Cart Operations
	MergeCarts(w http.ResponseWriter, r *http.Request)
	ClearCart(w http.ResponseWriter, r *http.Request)
	GetCartAnalytics(w http.ResponseWriter, r *http.Request)

	// Wishlist Management
	CreateWishlist(w http.ResponseWriter, r *http.Request)
	GetWishlist(w http.ResponseWriter, r *http.Request)
	GetWishlists(w http.ResponseWriter, r *http.Request)
	UpdateWishlist(w http.ResponseWriter, r *http.Request)
	DeleteWishlist(w http.ResponseWriter, r *http.Request)

	// Wishlist Items
	AddItemToWishlist(w http.ResponseWriter, r *http.Request)
	GetWishlistItem(w http.ResponseWriter, r *http.Request)
	GetWishlistItems(w http.ResponseWriter, r *http.Request)
	UpdateWishlistItem(w http.ResponseWriter, r *http.Request)
	DeleteWishlistItem(w http.ResponseWriter, r *http.Request)
	MoveItemToCart(w http.ResponseWriter, r *http.Request)
}

type cartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) ICartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}

// Cart Management

func (h *cartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	cart, err := h.cartService.CreateCart(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create cart", err)
		return
	}

	response := dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		SessionID: cart.SessionID,
		Currency:  cart.Currency,
		CreatedAt: cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if cart.ExpiresAt != nil {
		expiresAt := cart.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
		response.ExpiresAt = &expiresAt
	}

	httpx.Created(w, "Cart created successfully", response)
}

func (h *cartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	cart, err := h.cartService.GetCartByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Cart not found", err)
		return
	}

	response := dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		SessionID: cart.SessionID,
		Currency:  cart.Currency,
		CreatedAt: cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if cart.ExpiresAt != nil {
		expiresAt := cart.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
		response.ExpiresAt = &expiresAt
	}

	httpx.OK(w, "Cart retrieved successfully", response)
}

func (h *cartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.UpdateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	cart, err := h.cartService.UpdateCart(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update cart", err)
		return
	}

	response := dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		SessionID: cart.SessionID,
		Currency:  cart.Currency,
		CreatedAt: cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if cart.ExpiresAt != nil {
		expiresAt := cart.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
		response.ExpiresAt = &expiresAt
	}

	httpx.OK(w, "Cart updated successfully", response)
}

func (h *cartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	err = h.cartService.DeleteCart(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete cart", err)
		return
	}

	httpx.OK(w, "Cart deleted successfully", nil)
}

func (h *cartHandler) GetOrCreateCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	sessionID := r.URL.Query().Get("session_id")
	currency := r.URL.Query().Get("currency")

	if sessionID == "" || currency == "" {
		httpx.Error(w, http.StatusBadRequest, "Missing required parameters", errors.New("session_id and currency are required"))
		return
	}

	var userID *int64
	if userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "Invalid user ID", err)
			return
		}
		userID = &id
	}

	cart, err := h.cartService.GetOrCreateCart(r.Context(), userID, sessionID, currency)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get or create cart", err)
		return
	}

	response := dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		SessionID: cart.SessionID,
		Currency:  cart.Currency,
		CreatedAt: cart.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if cart.ExpiresAt != nil {
		expiresAt := cart.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
		response.ExpiresAt = &expiresAt
	}

	httpx.OK(w, "Cart retrieved or created successfully", response)
}

// Cart Items

func (h *cartHandler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	item, err := h.cartService.AddItemToCart(r.Context(), cartID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to add item to cart", err)
		return
	}

	response := dto.CartItemResponse{
		ID:               item.ID,
		CartID:           item.CartID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		Quantity:         item.Quantity,
		UnitPrice:        item.UnitPrice,
		TotalPrice:       item.TotalPrice,
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Item added to cart successfully", response)
}

func (h *cartHandler) GetCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	item, err := h.cartService.GetCartItemByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Cart item not found", err)
		return
	}

	response := dto.CartItemResponse{
		ID:               item.ID,
		CartID:           item.CartID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		Quantity:         item.Quantity,
		UnitPrice:        item.UnitPrice,
		TotalPrice:       item.TotalPrice,
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Cart item retrieved successfully", response)
}

func (h *cartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	item, err := h.cartService.UpdateCartItem(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update cart item", err)
		return
	}

	response := dto.CartItemResponse{
		ID:               item.ID,
		CartID:           item.CartID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		Quantity:         item.Quantity,
		UnitPrice:        item.UnitPrice,
		TotalPrice:       item.TotalPrice,
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Cart item updated successfully", response)
}

func (h *cartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	err = h.cartService.DeleteCartItem(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete cart item", err)
		return
	}

	httpx.OK(w, "Cart item deleted successfully", nil)
}

func (h *cartHandler) GetCartItems(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	items, err := h.cartService.GetCartItems(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get cart items", err)
		return
	}

	responses := make([]dto.CartItemResponse, len(items))
	for i, item := range items {
		responses[i] = dto.CartItemResponse{
			ID:               item.ID,
			CartID:           item.CartID,
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			TotalPrice:       item.TotalPrice,
			CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	httpx.OK(w, "Cart items retrieved successfully", responses)
}

func (h *cartHandler) ClearCartItems(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	err = h.cartService.ClearCartItems(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to clear cart items", err)
		return
	}

	httpx.OK(w, "Cart items cleared successfully", nil)
}

// Cart Summary & Calculations

func (h *cartHandler) GetCartSummary(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	summary, err := h.cartService.GetCartSummary(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get cart summary", err)
		return
	}

	httpx.OK(w, "Cart summary retrieved successfully", summary)
}

func (h *cartHandler) GetCartTotal(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	total, err := h.cartService.CalculateCartTotal(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to calculate cart total", err)
		return
	}

	httpx.OK(w, "Cart total calculated successfully", map[string]float64{"total": total})
}

func (h *cartHandler) GetCartItemCount(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	count, err := h.cartService.GetCartItemCount(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get cart item count", err)
		return
	}

	httpx.OK(w, "Cart item count retrieved successfully", map[string]int{"count": count})
}

// Cart Coupons

func (h *cartHandler) ApplyCouponToCart(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.ApplyCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	coupon, err := h.cartService.ApplyCouponToCart(r.Context(), cartID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to apply coupon", err)
		return
	}

	response := dto.CartCouponResponse{
		ID:             coupon.ID,
		CartID:         coupon.CartID,
		CouponCode:     coupon.CouponCode,
		DiscountAmount: coupon.DiscountAmount,
		CreatedAt:      coupon.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Coupon applied successfully", response)
}

func (h *cartHandler) RemoveCouponFromCart(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.RemoveCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = h.cartService.RemoveCouponFromCart(r.Context(), cartID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to remove coupon", err)
		return
	}

	httpx.OK(w, "Coupon removed successfully", nil)
}

func (h *cartHandler) GetCartCoupons(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	coupons, err := h.cartService.GetCartCoupons(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get cart coupons", err)
		return
	}

	responses := make([]dto.CartCouponResponse, len(coupons))
	for i, coupon := range coupons {
		responses[i] = dto.CartCouponResponse{
			ID:             coupon.ID,
			CartID:         coupon.CartID,
			CouponCode:     coupon.CouponCode,
			DiscountAmount: coupon.DiscountAmount,
			CreatedAt:      coupon.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	httpx.OK(w, "Cart coupons retrieved successfully", responses)
}

// Cart Shipping

func (h *cartHandler) SetCartShipping(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.SetShippingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	shipping, err := h.cartService.SetCartShipping(r.Context(), cartID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to set cart shipping", err)
		return
	}

	response := dto.CartShippingResponse{
		ID:               shipping.ID,
		CartID:           shipping.CartID,
		ShippingMethodID: shipping.ShippingMethodID,
		ShippingMethod:   shipping.ShippingMethod,
		ShippingAmount:   shipping.ShippingAmount,
		EstimatedDays:    shipping.EstimatedDays,
		CreatedAt:        shipping.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Cart shipping set successfully", response)
}

func (h *cartHandler) UpdateCartShipping(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.UpdateShippingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	shipping, err := h.cartService.UpdateCartShipping(r.Context(), cartID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update cart shipping", err)
		return
	}

	response := dto.CartShippingResponse{
		ID:               shipping.ID,
		CartID:           shipping.CartID,
		ShippingMethodID: shipping.ShippingMethodID,
		ShippingMethod:   shipping.ShippingMethod,
		ShippingAmount:   shipping.ShippingAmount,
		EstimatedDays:    shipping.EstimatedDays,
		CreatedAt:        shipping.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Cart shipping updated successfully", response)
}

func (h *cartHandler) GetCartShipping(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	shipping, err := h.cartService.GetCartShipping(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get cart shipping", err)
		return
	}

	if shipping == nil {
		httpx.OK(w, "No shipping information found", nil)
		return
	}

	response := dto.CartShippingResponse{
		ID:               shipping.ID,
		CartID:           shipping.CartID,
		ShippingMethodID: shipping.ShippingMethodID,
		ShippingMethod:   shipping.ShippingMethod,
		ShippingAmount:   shipping.ShippingAmount,
		EstimatedDays:    shipping.EstimatedDays,
		CreatedAt:        shipping.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Cart shipping retrieved successfully", response)
}

func (h *cartHandler) DeleteCartShipping(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	err = h.cartService.DeleteCartShipping(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete cart shipping", err)
		return
	}

	httpx.OK(w, "Cart shipping deleted successfully", nil)
}

// Cart Operations

func (h *cartHandler) MergeCarts(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	targetCartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.MergeCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = h.cartService.MergeCarts(r.Context(), req.SourceCartID, targetCartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to merge carts", err)
		return
	}

	httpx.OK(w, "Carts merged successfully", nil)
}

func (h *cartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	cartIDStr := chi.URLParam(r, "id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	var req dto.ClearCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if !req.Confirm {
		httpx.Error(w, http.StatusBadRequest, "Confirmation required", errors.New("Please confirm cart clearing"))
		return
	}

	err = h.cartService.ClearCart(r.Context(), cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to clear cart", err)
		return
	}

	httpx.OK(w, "Cart cleared successfully", nil)
}

func (h *cartHandler) GetCartAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.cartService.GetCartAnalytics(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get cart analytics", err)
		return
	}

	httpx.OK(w, "Cart analytics retrieved successfully", analytics)
}

// Wishlist Management

func (h *cartHandler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.CreateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	wishlist, err := h.cartService.CreateWishlist(r.Context(), userID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create wishlist", err)
		return
	}

	response := dto.WishlistResponse{
		ID:        wishlist.ID,
		UserID:    wishlist.UserID,
		Name:      wishlist.Name,
		IsPublic:  wishlist.IsPublic,
		CreatedAt: wishlist.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: wishlist.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Wishlist created successfully", response)
}

func (h *cartHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}

	wishlist, err := h.cartService.GetWishlistByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Wishlist not found", err)
		return
	}

	response := dto.WishlistResponse{
		ID:        wishlist.ID,
		UserID:    wishlist.UserID,
		Name:      wishlist.Name,
		IsPublic:  wishlist.IsPublic,
		CreatedAt: wishlist.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: wishlist.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Wishlist retrieved successfully", response)
}

func (h *cartHandler) GetWishlists(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	response, err := h.cartService.GetWishlistsByUserID(r.Context(), userID, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get wishlists", err)
		return
	}

	httpx.OK(w, "Wishlists retrieved successfully", response)
}

func (h *cartHandler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}

	var req dto.UpdateWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	wishlist, err := h.cartService.UpdateWishlist(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update wishlist", err)
		return
	}

	response := dto.WishlistResponse{
		ID:        wishlist.ID,
		UserID:    wishlist.UserID,
		Name:      wishlist.Name,
		IsPublic:  wishlist.IsPublic,
		CreatedAt: wishlist.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: wishlist.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Wishlist updated successfully", response)
}

func (h *cartHandler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}

	err = h.cartService.DeleteWishlist(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete wishlist", err)
		return
	}

	httpx.OK(w, "Wishlist deleted successfully", nil)
}

// Wishlist Items

func (h *cartHandler) AddItemToWishlist(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := chi.URLParam(r, "id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}

	var req dto.AddToWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	item, err := h.cartService.AddItemToWishlist(r.Context(), wishlistID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to add item to wishlist", err)
		return
	}

	response := dto.WishlistItemResponse{
		ID:               item.ID,
		WishlistID:       item.WishlistID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		Notes:            item.Notes,
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Item added to wishlist successfully", response)
}

func (h *cartHandler) GetWishlistItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	item, err := h.cartService.GetWishlistItemByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Wishlist item not found", err)
		return
	}

	response := dto.WishlistItemResponse{
		ID:               item.ID,
		WishlistID:       item.WishlistID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		Notes:            item.Notes,
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Wishlist item retrieved successfully", response)
}

func (h *cartHandler) GetWishlistItems(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := chi.URLParam(r, "id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid wishlist ID", err)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	response, err := h.cartService.GetWishlistItems(r.Context(), wishlistID, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get wishlist items", err)
		return
	}

	httpx.OK(w, "Wishlist items retrieved successfully", response)
}

func (h *cartHandler) UpdateWishlistItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	var req dto.UpdateWishlistItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	item, err := h.cartService.UpdateWishlistItem(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update wishlist item", err)
		return
	}

	response := dto.WishlistItemResponse{
		ID:               item.ID,
		WishlistID:       item.WishlistID,
		ProductID:        item.ProductID,
		ProductVariantID: item.ProductVariantID,
		Notes:            item.Notes,
		CreatedAt:        item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Wishlist item updated successfully", response)
}

func (h *cartHandler) DeleteWishlistItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	err = h.cartService.DeleteWishlistItem(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete wishlist item", err)
		return
	}

	httpx.OK(w, "Wishlist item deleted successfully", nil)
}

func (h *cartHandler) MoveItemToCart(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid item ID", err)
		return
	}

	cartIDStr := r.URL.Query().Get("cart_id")
	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart ID", err)
		return
	}

	err = h.cartService.MoveItemToCart(r.Context(), itemID, cartID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to move item to cart", err)
		return
	}

	httpx.OK(w, "Item moved to cart successfully", nil)
}
