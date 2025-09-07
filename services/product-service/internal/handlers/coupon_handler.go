package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

type CouponHandler interface {
	// Coupon Management
	CreateCoupon(w http.ResponseWriter, r *http.Request)
	GetCoupon(w http.ResponseWriter, r *http.Request)
	UpdateCoupon(w http.ResponseWriter, r *http.Request)
	DeleteCoupon(w http.ResponseWriter, r *http.Request)
	ListCoupons(w http.ResponseWriter, r *http.Request)

	// Coupon Validation
	ValidateCoupon(w http.ResponseWriter, r *http.Request)

	// Coupon Usage
	GetCouponUsage(w http.ResponseWriter, r *http.Request)
}

type couponHandler struct {
	couponService services.CouponService
}

func NewCouponHandler(couponService services.CouponService) CouponHandler {
	return &couponHandler{
		couponService: couponService,
	}
}

// Coupon Management

// CreateCoupon handles POST /api/v1/coupons
func (h *couponHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	coupon, err := h.couponService.CreateCoupon(r.Context(), &req)
	if err != nil {
		if err.Error() == "coupon with code already exists" {
			httpx.Error(w, http.StatusConflict, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to create coupon", err)
		}
		return
	}

	response := dto.CouponResponse{
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

	httpx.Created(w, "Coupon created successfully", response)
}

// GetCoupon handles GET /api/v1/coupons/{id}
func (h *couponHandler) GetCoupon(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid coupon ID", err)
		return
	}

	coupon, err := h.couponService.GetCouponByID(r.Context(), id)
	if err != nil {
		if err.Error() == "coupon with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to get coupon", err)
		}
		return
	}

	response := dto.CouponResponse{
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

	httpx.OK(w, "Coupon retrieved successfully", response)
}

// UpdateCoupon handles PUT /api/v1/coupons/{id}
func (h *couponHandler) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid coupon ID", err)
		return
	}

	var req dto.UpdateCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	coupon, err := h.couponService.UpdateCoupon(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "coupon with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to update coupon", err)
		}
		return
	}

	response := dto.CouponResponse{
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

	httpx.OK(w, "Coupon updated successfully", response)
}

// DeleteCoupon handles DELETE /api/v1/coupons/{id}
func (h *couponHandler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid coupon ID", err)
		return
	}

	err = h.couponService.DeleteCoupon(r.Context(), id)
	if err != nil {
		if err.Error() == "coupon with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to delete coupon", err)
		}
		return
	}

	httpx.OK(w, "Coupon deleted successfully", nil)
}

// ListCoupons handles GET /api/v1/coupons
func (h *couponHandler) ListCoupons(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := &dto.ListCouponsRequest{
		Code:     getStringParam(r, "code"),
		Type:     getStringParam(r, "type"),
		IsActive: getBoolParam(r, "is_active"),
		Search:   r.URL.Query().Get("search"),
		Page:     getIntParam(r, "page", 1),
		Limit:    getIntParam(r, "limit", 10),
	}

	response, err := h.couponService.ListCoupons(r.Context(), req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to list coupons", err)
		return
	}

	httpx.OK(w, "Coupons retrieved successfully", response)
}

// Coupon Validation

// ValidateCoupon handles POST /api/v1/coupons/validate
func (h *couponHandler) ValidateCoupon(w http.ResponseWriter, r *http.Request) {
	var req dto.ValidateCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	response, err := h.couponService.ValidateCoupon(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to validate coupon", err)
		return
	}

	httpx.OK(w, "Coupon validation completed", response)
}

// Coupon Usage

// GetCouponUsage handles GET /api/v1/coupons/usage
func (h *couponHandler) GetCouponUsage(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := &dto.ListCouponUsageRequest{
		CouponID: getInt64Param(r, "coupon_id"),
		UserID:   getInt64Param(r, "user_id"),
		OrderID:  getInt64Param(r, "order_id"),
		Page:     getIntParam(r, "page", 1),
		Limit:    getIntParam(r, "limit", 10),
	}

	response, err := h.couponService.GetCouponUsage(r.Context(), req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get coupon usage", err)
		return
	}

	httpx.OK(w, "Coupon usage retrieved successfully", response)
}

// Helper functions

func getStringParam(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func getBoolParam(r *http.Request, key string) *bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	if value == "true" {
		return &[]bool{true}[0]
	}
	return &[]bool{false}[0]
}

func getIntParam(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

func getInt64Param(r *http.Request, key string) *int64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return &intValue
	}
	return nil
}
