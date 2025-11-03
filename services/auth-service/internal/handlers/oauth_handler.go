package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
	"github.com/jattinmanhas/GearboxV2/services/shared/middleware"
)

type IOAuthHandler interface {
	InitiateOAuth(w http.ResponseWriter, r *http.Request)
	HandleOAuthCallback(w http.ResponseWriter, r *http.Request)
	LinkOAuthProvider(w http.ResponseWriter, r *http.Request)
	UnlinkOAuthProvider(w http.ResponseWriter, r *http.Request)
	GetLinkedProviders(w http.ResponseWriter, r *http.Request)
}

type oauthHandler struct {
	oauthService services.IOAuthService
	authHandler  IAuthHandler
	environment  string
	frontendURL  string
}

func NewOAuthHandler(oauthService services.IOAuthService, authHandler IAuthHandler, environment string, frontendURL string) IOAuthHandler {
	return &oauthHandler{
		oauthService: oauthService,
		authHandler:  authHandler,
		environment:  environment,
		frontendURL:  frontendURL,
	}
}

// InitiateOAuth initiates the OAuth flow by redirecting to the provider
func (h *oauthHandler) InitiateOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		httpx.Error(w, http.StatusBadRequest, "provider is required", nil)
		return
	}

	authURL, state, err := h.oauthService.InitiateOAuth(r.Context(), provider)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "failed to initiate oauth", err)
		return
	}

	// Store state in cookie for validation
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.environment == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600, // 10 minutes
	})

	// Return the auth URL or redirect directly
	// For API, return JSON. For web flow, you might redirect
	httpx.OK(w, "oauth initiated", map[string]any{
		"auth_url": authURL,
		"state":    state,
	})
}

// HandleOAuthCallback handles the OAuth callback from the provider
func (h *oauthHandler) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		http.Redirect(w, r, h.frontendURL+"/auth/oauth/error?message=invalid_provider", http.StatusTemporaryRedirect)
		return
	}

	// Get code and state from query params
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		http.Redirect(w, r, h.frontendURL+"/auth/oauth/error?message=missing_parameters", http.StatusTemporaryRedirect)
		return
	}

	// Verify state (optional but recommended for security)
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		http.Redirect(w, r, h.frontendURL+"/auth/oauth/error?message=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.environment == "production",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// Handle OAuth callback
	_, refreshToken, accessToken, err := h.oauthService.HandleOAuthCallback(r.Context(), provider, code, state)
	if err != nil {
		http.Redirect(w, r, h.frontendURL+"/auth/oauth/error?message="+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	// Set auth cookies (similar to login handler)
	h.setAccessTokenCookie(w, accessToken)
	h.setRefreshTokenCookie(w, refreshToken.RefreshToken)

	// Redirect to frontend success page
	http.Redirect(w, r, h.frontendURL+"/auth/oauth/success", http.StatusTemporaryRedirect)
}

// LinkOAuthProvider links an OAuth provider to the authenticated user's account
func (h *oauthHandler) LinkOAuthProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		httpx.Error(w, http.StatusBadRequest, "provider is required", nil)
		return
	}

	// Get user ID from JWT claims
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	var req dto.LinkOAuthProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate the request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
		return
	}

	// Link OAuth provider
	if err := h.oauthService.LinkOAuthProvider(r.Context(), claims.UserID, provider, req.Code); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to link oauth provider", err)
		return
	}

	httpx.OK(w, "oauth provider linked successfully", map[string]any{
		"provider": provider,
		"message":  "OAuth provider linked successfully to your account",
	})
}

// UnlinkOAuthProvider unlinks an OAuth provider from the authenticated user's account
func (h *oauthHandler) UnlinkOAuthProvider(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		httpx.Error(w, http.StatusBadRequest, "provider is required", nil)
		return
	}

	// Get user ID from JWT claims
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	// Unlink OAuth provider
	if err := h.oauthService.UnlinkOAuthProvider(r.Context(), claims.UserID, provider); err != nil {
		httpx.Error(w, http.StatusBadRequest, "failed to unlink oauth provider", err)
		return
	}

	httpx.OK(w, "oauth provider unlinked successfully", map[string]any{
		"provider": provider,
		"message":  "OAuth provider unlinked successfully from your account",
	})
}

// GetLinkedProviders retrieves all OAuth providers linked to the authenticated user's account
func (h *oauthHandler) GetLinkedProviders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT claims
	claims, ok := r.Context().Value(middleware.ClaimsContextKey).(*middleware.Claims)
	if !ok || claims == nil {
		httpx.Error(w, http.StatusUnauthorized, "invalid claims", nil)
		return
	}

	// Get linked providers
	providers, err := h.oauthService.GetLinkedProviders(r.Context(), claims.UserID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get linked providers", err)
		return
	}

	// Convert to response format
	var response []dto.OAuthProviderResponse
	for _, p := range providers {
		response = append(response, dto.OAuthProviderResponse{
			ID:             p.ID,
			Provider:       p.Provider,
			ProviderUserID: p.ProviderUserID,
			Email:          p.Email,
			LinkedAt:       p.CreatedAt,
		})
	}

	httpx.OK(w, "linked providers retrieved successfully", dto.LinkedProvidersResponse{
		Providers: response,
		Count:     len(response),
	})
}

// Helper methods to set cookies (borrowed from auth handler pattern)
func (h *oauthHandler) setAccessTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.environment == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes
	})
}

func (h *oauthHandler) setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.environment == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})
}
