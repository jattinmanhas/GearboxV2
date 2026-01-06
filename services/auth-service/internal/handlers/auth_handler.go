package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/helpers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	"github.com/jattinmanhas/GearboxV2/services/shared/middleware"
)

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
	CleanupExpiredTokens(w http.ResponseWriter, r *http.Request)
	// Profile methods
	GetProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	// User Analytics
	GetUserAnalytics(w http.ResponseWriter, r *http.Request)
	// Password reset methods
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	userService  services.IUserService
	authService  services.IAuthService
	emailService services.IEmailService
	jwtService   *jwt.JWTService
	environment  string
}

func NewAuthHandler(userService services.IUserService, authService services.IAuthService, emailService services.IEmailService, jwtService *jwt.JWTService, environment string) IAuthHandler {
	return &authHandler{
		userService:  userService,
		authService:  authService,
		emailService: emailService,
		jwtService:   jwtService,
		environment:  environment,
	}
}

func (h *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
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
		MiddleName:  domain.NewNullString(strings.TrimSpace(req.MiddleName)),
		LastName:    domain.NewNullString(strings.TrimSpace(req.LastName)),
		Avatar:      domain.NewNullString(req.Avatar),
		Gender:      domain.NewNullString(req.Gender),
		DateOfBirth: domain.NewNullTime(req.DateOfBirth),
	}

	if err := h.userService.RegisterNewUser(r.Context(), user); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to register user", err)
		return
	}

	httpx.Created(w, "user registered", map[string]any{"id": user.ID, "username": user.Username, "email": user.Email})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
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
	ipAddress := helpers.ExtractClientIP(r)

	// Authenticate user and generate tokens
	user, refreshToken, accessToken, err := h.authService.Login(r.Context(), req.Username, req.Password, userAgent, ipAddress)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}

	// Set only refresh token in HTTP-only cookie
	// Access token will be sent in response body for client to store in memory
	helpers.SetRefreshTokenCookie(w, refreshToken.RefreshToken, h.jwtService.GetRefreshTokenExpiry(), h.environment)

	// Return success response with access token in body
	httpx.OK(w, "login successful", map[string]any{
		"user": map[string]any{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName.String,
			"avatar":    user.Avatar.String,
			"role":      user.Role,
		},
		"access_token": accessToken, // NEW: Access token in response body
		"message":      "Login successful. Access token provided in response, refresh token in HTTP-only cookie.",
	})
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from cookie (more secure than request body)
	refreshToken := services.ExtractTokenFromCookie(r, "refresh_token")
	if refreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "refresh token required in cookie", nil)
		return
	}

	// Validate refresh token and generate new access token
	// Note: Refresh token is NOT rotated (same token remains valid)
	user, accessToken, err := h.authService.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	// Return success response with new access token in body
	// Refresh token cookie remains unchanged (no rotation)
	httpx.OK(w, "token refreshed successfully", map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"access_token": accessToken, // NEW: Access token in response body
		"message":      "Access token refreshed successfully. Refresh token remains valid.",
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
	helpers.ClearAuthCookies(w, h.environment)

	httpx.OK(w, "logout successful", map[string]any{
		"message": "Logout successful. All tokens cleared.",
	})
}

func (h *authHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	// Revoke all user tokens
	if err := h.authService.LogoutAll(r.Context(), claims.UserID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to logout from all devices", err)
		return
	}

	// Clear cookies
	helpers.ClearAuthCookies(w, h.environment)

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
	// Parse pagination parameters (frontend uses page-based pagination)
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		httpx.Error(w, http.StatusBadRequest, "invalid page", err)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		httpx.Error(w, http.StatusBadRequest, "invalid limit", err)
		return
	}

	// Parse filter parameters
	search := r.URL.Query().Get("search")
	isActiveStr := r.URL.Query().Get("is_active")
	roleIDStr := r.URL.Query().Get("role_id")

	var isActive *bool
	if isActiveStr != "" {
		if isActiveStr == "true" {
			val := true
			isActive = &val
		} else if isActiveStr == "false" {
			val := false
			isActive = &val
		}
	}

	var roleID *int
	if roleIDStr != "" {
		if id, err := strconv.Atoi(roleIDStr); err == nil {
			roleID = &id
		}
	}

	// Convert page-based pagination to offset-based
	offset := (page - 1) * limit

	// Get users with filters
	users, err := h.userService.GetAllUsersWithFilters(r.Context(), limit, offset, search, isActive, roleID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get users", err)
		return
	}

	// Get total count for pagination with filters
	total, err := h.userService.GetUsersCountWithFilters(r.Context(), search, isActive, roleID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get users count", err)
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	// Return response in the format expected by frontend
	response := map[string]interface{}{
		"users":       users,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	httpx.OK(w, "fetched users", response)
}

func (h *authHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract claims to check if user is admin
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	// Only admins can update other users
	if claims.Role != "admin" {
		httpx.Error(w, http.StatusForbidden, "only admins can update users", nil)
		return
	}

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

	var req dto.UpdateUserRequest
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
		MiddleName: domain.NewNullString(strings.TrimSpace(req.MiddleName)),
		LastName:   domain.NewNullString(strings.TrimSpace(req.LastName)),
		Avatar:     domain.NewNullString(req.Avatar),
		Gender:     domain.NewNullString(req.Gender),
	}

	// Handle DateOfBirth separately since it's a pointer in the request
	if req.DateOfBirth != nil {
		updateData.DateOfBirth = domain.NewNullTime(*req.DateOfBirth)
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
		"lastName":  updatedUser.LastName.String,
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

	var req dto.ChangePasswordRequest
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

	// Revoke http-only cookies
	helpers.ClearAuthCookies(w, h.environment)

	httpx.OK(w, "user deleted successfully", nil)
}

func (h *authHandler) CleanupExpiredTokens(w http.ResponseWriter, r *http.Request) {
	if err := h.authService.CleanupExpiredTokens(r.Context()); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to cleanup expired tokens", err)
		return
	}

	httpx.OK(w, "expired tokens cleaned up successfully", nil)
}

// GetProfile retrieves the current user's profile information
func (h *authHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT claims
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	// Get user profile using GetUserByID
	user, err := h.userService.GetUserByID(r.Context(), int(claims.UserID))
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get profile", err)
		return
	}

	// Convert to profile response format
	profile := dto.GetProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName.String,
		LastName:    user.LastName.String,
		PhoneNumber: user.PhoneNumber.String,
		DateOfBirth: user.DateOfBirth.Time,
		Avatar:      user.Avatar.String,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	httpx.OK(w, "profile retrieved successfully", profile)
}

// UpdateProfile updates the current user's profile information
func (h *authHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT claims
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate the request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	// Convert request to domain.User for the service
	updateData := &domain.User{
		FirstName:   strings.TrimSpace(req.FirstName),
		MiddleName:  domain.NewNullString(strings.TrimSpace(req.MiddleName)),
		LastName:    domain.NewNullString(strings.TrimSpace(req.LastName)),
		PhoneNumber: domain.NewNullString(strings.TrimSpace(req.PhoneNumber)),
		Gender:      domain.NewNullString(req.Gender),
		Avatar:      domain.NewNullString(req.Avatar),
	}

	// Handle DateOfBirth separately since it's a pointer in the request
	if req.DateOfBirth != nil {
		updateData.DateOfBirth = domain.NewNullTime(*req.DateOfBirth)
	}

	// Call UpdateUser service method
	updatedUser, err := h.userService.UpdateUser(r.Context(), int(claims.UserID), updateData)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to update profile", err)
		return
	}

	// Convert to profile response format
	profile := dto.GetProfileResponse{
		ID:          updatedUser.ID,
		Username:    updatedUser.Username,
		Email:       updatedUser.Email,
		FirstName:   updatedUser.FirstName,
		MiddleName:  updatedUser.MiddleName.String,
		LastName:    updatedUser.LastName.String,
		PhoneNumber: updatedUser.PhoneNumber.String,
		DateOfBirth: updatedUser.DateOfBirth.Time,
		Avatar:      updatedUser.Avatar.String,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
	}

	httpx.OK(w, "profile updated successfully", profile)
}

// GetUserAnalytics handles GET /api/v1/users/analytics
func (h *authHandler) GetUserAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.userService.GetUserAnalytics(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get user analytics", err)
		return
	}

	httpx.OK(w, "user analytics retrieved successfully", analytics)
}

// ForgotPassword handles POST /api/v1/auth/forgot-password
// This endpoint initiates the password reset process by generating a token and sending it via email
// Accepts either email or username to find the user
func (h *authHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate that at least one identifier is provided
	if req.Email == "" && req.Username == "" {
		httpx.Error(w, http.StatusBadRequest, "either email or username is required", nil)
		return
	}

	// Validate the request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	// Process forgot password request
	// Note: We always return success to prevent enumeration attacks
	if err := h.authService.ForgotPassword(r.Context(), req.Email, req.Username, h.emailService); err != nil {
		// Log the error but still return success to the client
		// This prevents attackers from determining which emails/usernames are registered
		httpx.OK(w, "if an account with that email or username exists, a password reset link has been sent", nil)
		return
	}

	httpx.OK(w, "if an account with that email or username exists, a password reset link has been sent", nil)
}

// ResetPassword handles POST /api/v1/auth/reset-password
// This endpoint resets the user's password using a valid reset token
func (h *authHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate the request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	// Reset password
	if err := h.authService.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		httpx.Error(w, http.StatusBadRequest, "failed to reset password", err)
		return
	}

	httpx.OK(w, "password reset successfully", nil)
}
