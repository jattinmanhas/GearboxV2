package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/middleware"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

type IRoleHandler interface {
	GetAllRoles(w http.ResponseWriter, r *http.Request)
	GetUserRole(w http.ResponseWriter, r *http.Request)
	AssignRoleToUser(w http.ResponseWriter, r *http.Request)
	RemoveUserRole(w http.ResponseWriter, r *http.Request)
	GetMyRole(w http.ResponseWriter, r *http.Request)
	CheckPermission(w http.ResponseWriter, r *http.Request)
}

type RoleHandler struct {
	roleService services.IRoleService
}

func NewRoleHandler(roleService services.IRoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// GetAllRoles returns all available roles
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleService.GetAllRoles(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to fetch roles", err)
		return
	}

	httpx.OK(w, "roles retrieved successfully", roles)
}

// GetUserRole returns the role of a specific user
func (h *RoleHandler) GetUserRole(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path or query params
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		httpx.Error(w, http.StatusBadRequest, "user_id is required", nil)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid user_id", err)
		return
	}

	role, err := h.roleService.GetUserRole(r.Context(), uint(userID))
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "user role not found", err)
		return
	}

	httpx.OK(w, "user role retrieved successfully", role)
}

// AssignRoleToUser assigns a role to a user
func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		UserID uint `json:"user_id" validate:"required"`
		RoleID uint `json:"role_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Get the current user from context (who is assigning the role)
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "authentication required", nil)
		return
	}

	// Type assert to get the actual claims
	c, ok := claims.(*services.Claims)
	if !ok {
		httpx.Error(w, http.StatusInternalServerError, "invalid claims format", nil)
		return
	}

	// Assign role to user
	err := h.roleService.AssignRoleToUser(r.Context(), req.UserID, req.RoleID, c.UserID)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "failed to assign role", err)
		return
	}

	httpx.OK(w, "role assigned successfully", nil)
}

// RemoveUserRole removes a user's role (sets to default user role)
func (h *RoleHandler) RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req struct {
		UserID uint `json:"user_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Get the current user from context (who is removing the role)
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "authentication required", nil)
		return
	}

	// Type assert to get the actual claims
	c, ok := claims.(*services.Claims)
	if !ok {
		httpx.Error(w, http.StatusInternalServerError, "invalid claims format", nil)
		return
	}

	// Remove user role (sets to default user role)
	err := h.roleService.RemoveUserRole(r.Context(), req.UserID, c.UserID)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "failed to remove user role", err)
		return
	}

	httpx.OK(w, "user role removed successfully", nil)
}

// GetMyRole returns the current user's role
func (h *RoleHandler) GetMyRole(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "authentication required", nil)
		return
	}

	// Type assert to get the actual claims
	c, ok := claims.(*services.Claims)
	if !ok {
		httpx.Error(w, http.StatusInternalServerError, "invalid claims format", nil)
		return
	}

	role, err := h.roleService.GetUserRole(r.Context(), c.UserID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "user role not found", err)
		return
	}

	httpx.OK(w, "role retrieved successfully", role)
}

// CheckPermission checks if the current user has a specific permission level
func (h *RoleHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "authentication required", nil)
		return
	}

	// Type assert to get the actual claims
	c, ok := claims.(*services.Claims)
	if !ok {
		httpx.Error(w, http.StatusInternalServerError, "invalid claims format", nil)
		return
	}

	// Get required role from query params
	requiredRole := r.URL.Query().Get("role")
	if requiredRole == "" {
		httpx.Error(w, http.StatusBadRequest, "role parameter is required", nil)
		return
	}

	// Check if user has the required permission
	err := h.roleService.CheckUserPermission(r.Context(), c.UserID, requiredRole)
	if err != nil {
		httpx.Error(w, http.StatusForbidden, "insufficient permissions", err)
		return
	}

	httpx.OK(w, "permission check passed", map[string]interface{}{
		"user_id":        c.UserID,
		"user_role":      c.Role,
		"required_role":  requiredRole,
		"has_permission": true,
	})
}
