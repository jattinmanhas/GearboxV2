package middleware

import (
	"context"
	"net/http"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/httpx"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
)

type contextKey string

const (
	UserContextKey   contextKey = "user"
	ClaimsContextKey contextKey = "claims"
)

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware(authService services.IAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header or cookie
			token := services.ExtractTokenFromHeader(r)
			if token == "" {
				token = services.ExtractTokenFromCookie(r, "access_token")
			}

			if token == "" {
				httpx.Error(w, http.StatusUnauthorized, "access token required", nil)
				return
			}

			// Validate access token
			claims, err := authService.ValidateAccessToken(r.Context(), token)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid access token", err)
				return
			}

			// Get user from token
			user, err := authService.GetUserFromToken(r.Context(), token)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "user not found", err)
				return
			}

			// Add user and claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, ClaimsContextKey, claims)

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
					if user, err := authService.GetUserFromToken(r.Context(), token); err == nil {
						// Add user and claims to context
						ctx := context.WithValue(r.Context(), UserContextKey, user)
						ctx = context.WithValue(ctx, ClaimsContextKey, claims)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
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

// RequireRole middleware checks if the authenticated user has a specific role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r.Context())
			if user == nil {
				httpx.Error(w, http.StatusUnauthorized, "authentication required", nil)
				return
			}

			// TODO: Implement role checking when you add roles to your user model
			// For now, just allow authenticated users
			next.ServeHTTP(w, r)
		})
	}
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
