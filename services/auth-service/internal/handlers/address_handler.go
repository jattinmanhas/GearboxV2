package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

type IAddressHandler interface {
	// Address operations
	CreateAddress(w http.ResponseWriter, r *http.Request)
	GetAddresses(w http.ResponseWriter, r *http.Request)
	GetAddressByID(w http.ResponseWriter, r *http.Request)
	GetDefaultAddress(w http.ResponseWriter, r *http.Request)
	UpdateAddress(w http.ResponseWriter, r *http.Request)
	DeleteAddress(w http.ResponseWriter, r *http.Request)
	SetDefaultAddress(w http.ResponseWriter, r *http.Request)

	// Phone number operations
	CreatePhoneNumber(w http.ResponseWriter, r *http.Request)
	GetPhoneNumbers(w http.ResponseWriter, r *http.Request)
	GetPhoneNumberByID(w http.ResponseWriter, r *http.Request)
	GetPrimaryPhone(w http.ResponseWriter, r *http.Request)
	UpdatePhoneNumber(w http.ResponseWriter, r *http.Request)
	DeletePhoneNumber(w http.ResponseWriter, r *http.Request)
	SetPrimaryPhone(w http.ResponseWriter, r *http.Request)
}

type addressHandler struct {
	addressService services.IAddressService
}

func NewAddressHandler(addressService services.IAddressService) IAddressHandler {
	return &addressHandler{
		addressService: addressService,
	}
}

// Address operations

func (h *addressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	var req dto.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(&req); len(errs) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	address, err := h.addressService.CreateAddress(r.Context(), userID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to create address", err)
		return
	}

	httpx.Created(w, "address created successfully", address)
}

func (h *addressHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	addresses, err := h.addressService.GetAddressesByUserID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get addresses", err)
		return
	}

	httpx.OK(w, "addresses retrieved successfully", addresses)
}

func (h *addressHandler) GetAddressByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	// Get address ID from URL
	addressIDStr := chi.URLParam(r, "id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid address ID", nil)
		return
	}

	address, err := h.addressService.GetAddressByID(r.Context(), userID, uint(addressID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "address not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to get address", err)
		return
	}

	httpx.OK(w, "address retrieved successfully", address)
}

func (h *addressHandler) GetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	address, err := h.addressService.GetDefaultAddressByUserID(r.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "no default address found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to get default address", err)
		return
	}

	httpx.OK(w, "default address retrieved successfully", address)
}

func (h *addressHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	// Get address ID from URL
	addressIDStr := chi.URLParam(r, "id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid address ID", nil)
		return
	}

	var req dto.UpdateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(&req); len(errs) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	address, err := h.addressService.UpdateAddress(r.Context(), userID, uint(addressID), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "address not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to update address", err)
		return
	}

	httpx.OK(w, "address updated successfully", address)
}

func (h *addressHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	// Get address ID from URL
	addressIDStr := chi.URLParam(r, "id")
	addressID, err := strconv.ParseUint(addressIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid address ID", nil)
		return
	}

	err = h.addressService.DeleteAddress(r.Context(), userID, uint(addressID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "address not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to delete address", err)
		return
	}

	httpx.OK(w, "address deleted successfully", nil)
}

func (h *addressHandler) SetDefaultAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	var req dto.SetDefaultAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(&req); len(errs) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	err := h.addressService.SetDefaultAddress(r.Context(), userID, req.AddressID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "address not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to set default address", err)
		return
	}

	httpx.OK(w, "default address set successfully", nil)
}

// Phone number operations

func (h *addressHandler) CreatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	var req dto.CreatePhoneNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(&req); len(errs) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	phone, err := h.addressService.CreatePhoneNumber(r.Context(), userID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to create phone number", err)
		return
	}

	httpx.Created(w, "phone number created successfully", phone)
}

func (h *addressHandler) GetPhoneNumbers(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	phones, err := h.addressService.GetPhoneNumbersByUserID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get phone numbers", err)
		return
	}

	httpx.OK(w, "phone numbers retrieved successfully", phones)
}

func (h *addressHandler) GetPhoneNumberByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	// Get phone ID from URL
	phoneIDStr := chi.URLParam(r, "id")
	phoneID, err := strconv.ParseUint(phoneIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid phone ID", nil)
		return
	}

	phone, err := h.addressService.GetPhoneNumberByID(r.Context(), userID, uint(phoneID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "phone number not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to get phone number", err)
		return
	}

	httpx.OK(w, "phone number retrieved successfully", phone)
}

func (h *addressHandler) GetPrimaryPhone(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	phone, err := h.addressService.GetPrimaryPhoneByUserID(r.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "no primary phone found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to get primary phone", err)
		return
	}

	httpx.OK(w, "primary phone retrieved successfully", phone)
}

func (h *addressHandler) UpdatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	// Get phone ID from URL
	phoneIDStr := chi.URLParam(r, "id")
	phoneID, err := strconv.ParseUint(phoneIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid phone ID", nil)
		return
	}

	var req dto.UpdatePhoneNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(&req); len(errs) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	phone, err := h.addressService.UpdatePhoneNumber(r.Context(), userID, uint(phoneID), &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "phone number not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to update phone number", err)
		return
	}

	httpx.OK(w, "phone number updated successfully", phone)
}

func (h *addressHandler) DeletePhoneNumber(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	// Get phone ID from URL
	phoneIDStr := chi.URLParam(r, "id")
	phoneID, err := strconv.ParseUint(phoneIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid phone ID", nil)
		return
	}

	err = h.addressService.DeletePhoneNumber(r.Context(), userID, uint(phoneID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "phone number not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to delete phone number", err)
		return
	}

	httpx.OK(w, "phone number deleted successfully", nil)
}

func (h *addressHandler) SetPrimaryPhone(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	var req dto.SetPrimaryPhoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid JSON", nil)
		return
	}

	// Validate request
	if errs := validation.ValidateStruct(&req); len(errs) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	err := h.addressService.SetPrimaryPhone(r.Context(), userID, req.PhoneID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, "phone number not found", nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to set primary phone", err)
		return
	}

	httpx.OK(w, "primary phone set successfully", nil)
}
