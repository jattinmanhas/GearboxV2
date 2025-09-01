package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/httpx"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/validation"
)

type registerRequest struct {
	Username    string    `json:"username" validate:"required,username"`
	Password    string    `json:"password" validate:"required,password"`
	Email       string    `json:"email" validate:"required,email"`
	FirstName   string    `json:"first_name" validate:"required,min=1,max=50"`
	MiddleName  string    `json:"middle_name" validate:"omitempty,max=50"`
	LastName    string    `json:"last_name" validate:"omitempty,max=50"` // Made optional
	Avatar      string    `json:"avatar" validate:"omitempty,url"`
	Gender      string    `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"` // Made optional
	DateOfBirth time.Time `json:"date_of_birth" validate:"omitempty,date_of_birth"`                      // Made optional
}

type updateUserRequest struct {
	FirstName   string     `json:"first_name" validate:"omitempty,min=1,max=50"`
	MiddleName  string     `json:"middle_name" validate:"omitempty,max=50"`
	LastName    string     `json:"last_name" validate:"omitempty,min=1,max=50"`
	Avatar      string     `json:"avatar" validate:"omitempty,url"`
	Gender      string     `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"`
	DateOfBirth *time.Time `json:"date_of_birth" validate:"omitempty,date_of_birth"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password"`
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type IAuthHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	LogoutAll(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	userService services.IUserService
	authService services.IAuthService
	jwtService  *services.JWTService
}

func NewAuthHandler(userService services.IUserService, authService services.IAuthService, jwtService *services.JWTService) IAuthHandler {
	return &authHandler{
		userService: userService,
		authService: authService,
		jwtService:  jwtService,
	}
}

func (h *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate the request using our validation package
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
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
		DateOfBirth: req.DateOfBirth,
	}

	if err := h.userService.RegisterNewUser(r.Context(), user); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to register user", err)
		return
	}

	httpx.Created(w, "user registered", map[string]any{"id": user.ID, "username": user.Username, "email": user.Email})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate the request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	// Get user agent and IP address
	userAgent := r.UserAgent()
	ipAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ipAddress = r.RemoteAddr // fallback, but may include port
	}
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ipAddress = forwardedFor
	}

	// Authenticate user and generate tokens
	user, refreshToken, accessToken, err := h.authService.Login(r.Context(), req.Username, req.Password, userAgent, ipAddress)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}

	// Set HTTP-only cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to false in development
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.jwtService.GetAccessTokenExpiry().Seconds()),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to false in development
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.jwtService.GetRefreshTokenExpiry().Seconds()),
	})

	// Return success response
	httpx.OK(w, "login successful", map[string]any{
		"user": map[string]any{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
		},
		"message": "Login successful. Tokens stored in HTTP-only cookies.",
	})
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from cookie (more secure than request body)
	refreshToken := services.ExtractTokenFromCookie(r, "refresh_token")
	if refreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "refresh token required in cookie", nil)
		return
	}

	// Refresh tokens
	user, newRefreshToken, accessToken, err := h.authService.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	// Set new HTTP-only cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to false in development
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.jwtService.GetAccessTokenExpiry().Seconds()),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to false in development
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.jwtService.GetRefreshTokenExpiry().Seconds()),
	})

	// Return success response
	httpx.OK(w, "tokens refreshed successfully", map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"message": "Tokens refreshed successfully. New tokens stored in HTTP-only cookies.",
	})
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from cookie
	refreshToken := services.ExtractTokenFromCookie(r, "refresh_token")
	if refreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "refresh token not found", nil)
		return
	}

	// Revoke refresh token
	if err := h.authService.Logout(r.Context(), refreshToken); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to logout", err)
		return
	}

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	httpx.OK(w, "logout successful", map[string]any{
		"message": "Logout successful. All tokens cleared.",
	})
}

func (h *authHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	// Get user ID from access token
	accessToken := services.ExtractTokenFromCookie(r, "access_token")
	if accessToken == "" {
		httpx.Error(w, http.StatusUnauthorized, "access token not found", nil)
		return
	}

	// Validate access token and get user
	claims, err := h.authService.ValidateAccessToken(r.Context(), accessToken)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid access token", err)
		return
	}

	// Revoke all user tokens
	if err := h.authService.LogoutAll(r.Context(), claims.UserID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to logout from all devices", err)
		return
	}

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	httpx.OK(w, "logout from all devices successful", map[string]any{
		"message": "Logout from all devices successful. All tokens revoked.",
	})
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

func (h *authHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "id")
	if userIdStr == "" {
		httpx.Error(w, http.StatusBadRequest, "id is required", nil)
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid id", err)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	// Convert request to domain.User for the service
	updateData := &domain.User{
		FirstName:  strings.TrimSpace(req.FirstName),
		MiddleName: strings.TrimSpace(req.MiddleName),
		LastName:   strings.TrimSpace(req.LastName),
		Avatar:     req.Avatar,
		Gender:     req.Gender,
	}

	// Handle DateOfBirth separately since it's a pointer in the request
	if req.DateOfBirth != nil {
		updateData.DateOfBirth = *req.DateOfBirth
	}

	// Call the service to update the user (service handles merging with existing data)
	updatedUser, err := h.userService.UpdateUser(r.Context(), userId, updateData)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	httpx.OK(w, "user updated successfully", map[string]any{
		"id":        updatedUser.ID,
		"username":  updatedUser.Username,
		"firstName": updatedUser.FirstName,
		"lastName":  updatedUser.LastName,
		"email":     updatedUser.Email,
		"updatedAt": updatedUser.UpdatedAt,
	})
}

func (h *authHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	if err := h.userService.ChangePassword(r.Context(), id, req.CurrentPassword, req.NewPassword); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to change password", err)
		return
	}

	httpx.OK(w, "password changed successfully", nil)
}

func (h *authHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to delete user", err)
		return
	}

	httpx.OK(w, "user deleted successfully", nil)
}
