package helpers

import (
	"net"
	"net/http"
	"strings"
	"time"
)

// SetRefreshTokenCookie sets the refresh token in an HTTP-only cookie
func SetRefreshTokenCookie(w http.ResponseWriter, token string, expiry time.Duration, environment string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   environment == "production", // Only secure in production
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(expiry.Seconds()),
	})
}

// ClearAuthCookies clears the refresh token cookie
// Note: Access tokens are no longer stored in cookies, so we only clear refresh_token
func ClearAuthCookies(w http.ResponseWriter, environment string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   environment == "production",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// ExtractClientIP extracts the client IP address from the request
// It checks for forwarded headers first (for proxy/load balancer scenarios)
// and falls back to RemoteAddr if no forwarded headers are present
func ExtractClientIP(r *http.Request) string {
	// Check for forwarded headers first (for proxy/load balancer scenarios)
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if commaIdx := strings.Index(forwardedFor, ","); commaIdx != -1 {
			return strings.TrimSpace(forwardedFor[:commaIdx])
		}
		return strings.TrimSpace(forwardedFor)
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	ipAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback, but may include port
	}
	return ipAddress
}
