package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// IAuthService defines the interface for authentication services
type IAuthService interface {
	ValidateAccessToken(ctx context.Context, token string) (*Claims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*Claims, error)
	GenerateAccessTokenFromUser(ctx context.Context, user interface{}) (string, error)
	ExtractTokenFromHeader(r *http.Request) string
	ExtractTokenFromCookie(r *http.Request, name string) string
}

// User represents a minimal user structure for token generation
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type contextKey string

const (
	UserContextKey      contextKey = "user"
	ClaimsContextKey    contextKey = "claims"
	AuthTokenContextKey contextKey = "auth_token"
)

// AuthMiddleware validates JWT tokens and extracts user information
// Simplified version: Only validates access tokens, no automatic refresh
// Frontend is responsible for calling /refresh endpoint when access token expires
func AuthMiddleware(authService IAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Debug logging
			authHeader := r.Header.Get("Authorization")
			fmt.Printf("[AuthMiddleware] Request: %s %s | Auth Header Present: %v\n", r.Method, r.URL.Path, authHeader != "")

			// Extract access token from Authorization header (preferred)
			accessToken := authService.ExtractTokenFromHeader(r)

			// Fallback to cookie for backward compatibility during transition
			if accessToken == "" {
				accessToken = authService.ExtractTokenFromCookie(r, "access_token")
			}

			// If no token found, return 401
			if accessToken == "" {
				httpx.Error(w, http.StatusUnauthorized, "access token required", nil)
				return
			}

			// Validate access token
			claims, err := authService.ValidateAccessToken(r.Context(), accessToken)
			if err != nil {
				// Token is invalid or expired - return 401
				// Frontend will detect this and call /refresh endpoint
				httpx.Error(w, http.StatusUnauthorized, "invalid or expired access token", err)
				return
			}

			// Token is valid, add claims to context and continue
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			ctx = context.WithValue(ctx, AuthTokenContextKey, accessToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware validates JWT tokens if present, but doesn't require them
func OptionalAuthMiddleware(authService IAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header or cookie
			token := authService.ExtractTokenFromHeader(r)
			if token == "" {
				token = authService.ExtractTokenFromCookie(r, "access_token")
			}

			if token != "" {
				// Try to validate token
				if claims, err := authService.ValidateAccessToken(r.Context(), token); err == nil {
					// Add claims and token to context (no user DB query needed)
					ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
					ctx = context.WithValue(ctx, AuthTokenContextKey, token)
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
	if c, ok := claims.(*Claims); ok {
		return c.UserID
	}

	return 0
}

// ExtractAuthTokenFromContext extracts auth token from context
func ExtractAuthTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AuthTokenContextKey).(string)
	return token, ok
}

// CanAccess checks if a user role can access a required role
// This is a simplified role hierarchy check
func CanAccess(userRole, requiredRole string) bool {
	// Define role hierarchy (higher number = more permissions)
	roleLevels := map[string]int{
		"user":   1,
		"editor": 2,
		"admin":  3,
	}

	userLevel, userExists := roleLevels[userRole]
	requiredLevel, requiredExists := roleLevels[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
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
			c, ok := claims.(*Claims)
			if !ok {
				httpx.Error(w, http.StatusInternalServerError, "invalid claims format", nil)
				return
			}

			// Check if user has the required role
			if !CanAccess(c.Role, requiredRole) {
				httpx.Error(w, http.StatusForbidden, "insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware checks if the authenticated user is an admin
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRole("admin")
}

// RequireEditor middleware checks if the authenticated user is an editor or admin
func RequireEditor() func(http.Handler) http.Handler {
	return RequireRole("editor")
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

// GetAuthTokenFromContext extracts the auth token from the request context
func GetAuthTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(AuthTokenContextKey).(string); ok {
		return token
	}
	return ""
}
