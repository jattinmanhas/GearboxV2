package dto

import "time"

// OAuthInitiateResponse is returned when initiating OAuth flow
type OAuthInitiateResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// OAuthCallbackRequest represents the OAuth callback parameters
type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

// OAuthUserInfoResponse represents user information from OAuth provider
type OAuthUserInfoResponse struct {
	ProviderUserID string `json:"provider_user_id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Avatar         string `json:"avatar"`
	Provider       string `json:"provider"`
}

// LinkOAuthProviderRequest is used to link an OAuth provider to an existing account
type LinkOAuthProviderRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google facebook github"`
	Code     string `json:"code" validate:"required"`
}

// UnlinkOAuthProviderRequest is used to unlink an OAuth provider from an account
type UnlinkOAuthProviderRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google facebook github"`
}

// OAuthProviderResponse represents a linked OAuth provider
type OAuthProviderResponse struct {
	ID             uint      `json:"id"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	Email          string    `json:"email"`
	LinkedAt       time.Time `json:"linked_at"`
}

// LinkedProvidersResponse represents all linked OAuth providers for a user
type LinkedProvidersResponse struct {
	Providers []OAuthProviderResponse `json:"providers"`
	Count     int                     `json:"count"`
}
