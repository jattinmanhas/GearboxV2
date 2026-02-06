package domain

import (
	"time"
)

// OAuthProvider represents the provider types supported
type OAuthProvider string

const (
	ProviderGoogle OAuthProvider = "google"
	ProviderGithub OAuthProvider = "github"
)

// ValidProviders returns a list of valid OAuth providers
func ValidProviders() []OAuthProvider {
	return []OAuthProvider{ProviderGoogle, ProviderGithub}
}

// IsValidProvider checks if a provider string is valid
func IsValidProvider(provider string) bool {
	for _, p := range ValidProviders() {
		if string(p) == provider {
			return true
		}
	}
	return false
}

// OAuthProviderLink represents a link between a user and an OAuth provider
type OAuthProviderLink struct {
	ID             uint       `json:"id" db:"id"`
	UserID         uint       `json:"user_id" db:"user_id"`
	Provider       string     `json:"provider" db:"provider"` // google, github
	ProviderUserID string     `json:"provider_user_id" db:"provider_user_id"`
	Email          string     `json:"email" db:"email"`
	AccessToken    *string    `json:"-" db:"access_token"`  // Don't expose in JSON
	RefreshToken   *string    `json:"-" db:"refresh_token"` // Don't expose in JSON
	ExpiresAt      *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// OAuthUserInfo represents user information returned from OAuth providers
type OAuthUserInfo struct {
	ProviderUserID string
	Email          string
	FirstName      string
	LastName       string
	Avatar         string
	Provider       OAuthProvider
}
