package middleware

import (
	"context"
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
	UserContextKey   contextKey = "user"
	ClaimsContextKey contextKey = "claims"
)

// AuthMiddleware validates JWT tokens and extracts user information
// Implements smart token refresh: if access token is expired but refresh token is valid,
// it will automatically refresh the access token
func AuthMiddleware(authService IAuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract access token from Authorization header or cookie
			accessToken := authService.ExtractTokenFromHeader(r)
			if accessToken == "" {
				accessToken = authService.ExtractTokenFromCookie(r, "access_token")
			}

			var claims *Claims
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
			refreshToken := authService.ExtractTokenFromCookie(r, "refresh_token")
			if refreshToken == "" {
				httpx.Error(w, http.StatusUnauthorized, "access token required", nil)
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
			minimalUser := &User{
				ID:       refreshClaims.UserID,
				Username: refreshClaims.Username,
				Email:    refreshClaims.Email,
				Role:     refreshClaims.Role,
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
			claims = &Claims{
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
	if c, ok := claims.(*Claims); ok {
		return c.UserID
	}

	return 0
}

// ExtractAuthTokenFromContext extracts auth token from context
func ExtractAuthTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value("auth_token").(string)
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
	if token, ok := ctx.Value("auth_token").(string); ok {
		return token
	}
	return ""
}
