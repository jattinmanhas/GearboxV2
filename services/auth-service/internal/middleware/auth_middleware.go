package middleware

import (
	"context"
	"net/http"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/httpx"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
)

type contextKey string

const (
	UserContextKey   contextKey = "user"
	ClaimsContextKey contextKey = "claims"
)

// AuthMiddleware validates JWT tokens and extracts user information
// Implements smart token refresh: if access token is expired but refresh token is valid,
// it will automatically refresh the access token
func AuthMiddleware(authService services.IAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract access token from Authorization header or cookie
			accessToken := services.ExtractTokenFromHeader(r)
			if accessToken == "" {
				accessToken = services.ExtractTokenFromCookie(r, "access_token")
			}

			var claims *services.Claims
			var err error

			// Try to validate access token first
			if accessToken != "" {
				claims, err = authService.ValidateAccessToken(r.Context(), accessToken)
				if err == nil {
					// Access token is valid, proceed normally
					ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// Access token is invalid or expired, try refresh token
			refreshToken := services.ExtractTokenFromCookie(r, "refresh_token")
			if refreshToken == "" {
				httpx.Error(w, http.StatusUnauthorized, "access token required", nil)
				return
			}

			// Validate refresh token
			_, err = authService.ValidateRefreshToken(r.Context(), refreshToken)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid refresh token", err)
				return
			}

			// Refresh token is valid, generate new access token
			// Use refresh token claims directly (no DB query needed)
			refreshClaims, err := authService.ValidateRefreshToken(r.Context(), refreshToken)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid refresh token", err)
				return
			}

			// Generate new access token using JWT service directly
			// Create minimal user object from claims (no DB query)
			minimalUser := &domain.User{
				ID:       refreshClaims.UserID,
				Username: refreshClaims.Username,
				Email:    refreshClaims.Email,
			}

			newAccessToken, err := authService.GenerateAccessTokenFromUser(r.Context(), minimalUser)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "failed to generate new access token", err)
				return
			}

			// Set new access token in cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    newAccessToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   900, // 15 minutes
			})

			// Use the claims from the minimal user object
			claims = &services.Claims{
				UserID:   minimalUser.ID,
				Username: minimalUser.Username,
				Email:    minimalUser.Email,
				Role:     minimalUser.Role,
			}

			// Add claims to context (no user DB query needed)
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)

			// Call next handler with updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware validates JWT tokens if present, but doesn't require them
func OptionalAuthMiddleware(authService services.IAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header or cookie
			token := services.ExtractTokenFromHeader(r)
			if token == "" {
				token = services.ExtractTokenFromCookie(r, "access_token")
			}

			if token != "" {
				// Try to validate token
				if claims, err := authService.ValidateAccessToken(r.Context(), token); err == nil {
					// Only add claims to context (no user DB query needed)
					ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			// Continue without authentication
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user from request context
func GetUserFromContext(ctx context.Context) interface{} {
	return ctx.Value(UserContextKey)
}

// GetClaimsFromContext extracts claims from request context
func GetClaimsFromContext(ctx context.Context) interface{} {
	return ctx.Value(ClaimsContextKey)
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) uint {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return 0
	}

	// Type assert to get the actual claims
	if c, ok := claims.(*services.Claims); ok {
		return c.UserID
	}

	return 0
}

// RequireRole middleware checks if the authenticated user has a specific role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
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

			// Check if user has the required role
			if !domain.CanAccess(c.Role, requiredRole) {
				httpx.Error(w, http.StatusForbidden, "insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware checks if the authenticated user is an admin
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRole(domain.RoleAdmin)
}

// RequireEditor middleware checks if the authenticated user is an editor or admin
func RequireEditor() func(http.Handler) http.Handler {
	return RequireRole(domain.RoleEditor)
}

// CORS middleware for handling cross-origin requests
func CORSMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
