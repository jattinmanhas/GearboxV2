package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/httpx"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
)

type registerRequest struct {
	Username    string     `json:"username"`
	Password    string     `json:"password"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	MiddleName  string     `json:"middle_name"`
	LastName    string     `json:"last_name"`
	Avatar      string     `json:"avatar"`
	Gender      string     `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type IAuthHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	userService services.IUserService
}

func NewAuthHandler(userService services.IUserService) IAuthHandler {
	return &authHandler{userService: userService}
}

func (h *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if req.DateOfBirth == nil {
		httpx.Error(w, http.StatusBadRequest, "date_of_birth is required", nil)
		return
	}

	user := &domain.User{
		Username:    strings.TrimSpace(req.Username),
		Password:    req.Password,
		Email:       strings.TrimSpace(req.Email),
		FirstName:   strings.TrimSpace(req.FirstName),
		MiddleName:  strings.TrimSpace(req.MiddleName),
		LastName:    strings.TrimSpace(req.LastName),
		Avatar:      req.Avatar,
		Gender:      req.Gender,
		DateOfBirth: *req.DateOfBirth,
	}

	if err := h.userService.RegisterNewUser(r.Context(), user); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to register user", err)
		return
	}

	httpx.Created(w, "user registered", map[string]any{"id": user.ID, "username": user.Username, "email": user.Email})
}

func (h *authHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		httpx.Error(w, http.StatusBadRequest, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id", err)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get user", err)
		return
	}

	httpx.OK(w, "user found", user)
}

func (h *authHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid limit", err)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid offset", err)
		return
	}

	users, err := h.userService.GetAllUsers(r.Context(), limit, offset)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get users", err)
		return
	}

	httpx.OK(w, "fetched users", users)
}
