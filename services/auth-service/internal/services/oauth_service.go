package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/helpers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type IOAuthService interface {
	InitiateOAuth(ctx context.Context, provider string) (authURL string, state string, err error)
	HandleOAuthCallback(ctx context.Context, provider string, code string, state string, ipAddress string, userAgent string) (*domain.User, *domain.RefreshToken, string, error)
	LinkOAuthProvider(ctx context.Context, userID uint, provider string, code string) error
	UnlinkOAuthProvider(ctx context.Context, userID uint, provider string) error
	GetLinkedProviders(ctx context.Context, userID uint) ([]*domain.OAuthProviderLink, error)
}

type oauthService struct {
	oauthRepo        repository.IOAuthRepository
	userRepo         repository.IUserRepository
	refreshTokenRepo repository.IRefreshTokenRepository
	jwtService       *jwt.JWTService
	googleConfig     *oauth2.Config
	githubConfig     *oauth2.Config
}

func NewOAuthService(
	oauthRepo repository.IOAuthRepository,
	userRepo repository.IUserRepository,
	refreshTokenRepo repository.IRefreshTokenRepository,
	jwtService *jwt.JWTService,
	googleClientID, googleClientSecret, googleRedirectURL string,
	githubClientID, githubClientSecret, githubRedirectURL string,
) IOAuthService {
	return &oauthService{
		oauthRepo:        oauthRepo,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		googleConfig: &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		githubConfig: &oauth2.Config{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  githubRedirectURL,
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
	}
}

// generateStateToken generates a random state token for OAuth flow
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// InitiateOAuth generates the OAuth authorization URL
func (s *oauthService) InitiateOAuth(ctx context.Context, provider string) (string, string, error) {
	if !domain.IsValidProvider(provider) {
		return "", "", fmt.Errorf("invalid provider: %s", provider)
	}

	state, err := generateStateToken()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate state token: %w", err)
	}

	var authURL string
	switch provider {
	case string(domain.ProviderGoogle):
		authURL = s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	case string(domain.ProviderGithub):
		authURL = s.githubConfig.AuthCodeURL(state)
	default:
		return "", "", fmt.Errorf("unsupported provider: %s", provider)
	}

	return authURL, state, nil
}

// HandleOAuthCallback handles the OAuth callback and creates or logs in a user
func (s *oauthService) HandleOAuthCallback(ctx context.Context, provider string, code string, state string, ipAddress string, userAgent string) (*domain.User, *domain.RefreshToken, string, error) {
	// Exchange code for token
	token, err := s.exchangeCodeForToken(ctx, provider, code)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Fetch user info from provider
	userInfo, err := s.fetchUserInfo(ctx, provider, token.AccessToken)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Check if OAuth provider already exists
	existingOAuthLink, err := s.oauthRepo.GetOAuthProviderByProviderAndUserID(ctx, provider, userInfo.ProviderUserID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to check existing oauth provider: %w", err)
	}

	var user *domain.User

	if existingOAuthLink != nil {
		// User already exists with this OAuth provider
		user, err = s.userRepo.GetUserByID(ctx, int(existingOAuthLink.UserID))
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to get user: %w", err)
		}

		// Update OAuth tokens
		// Update OAuth tokens
		existingOAuthLink.AccessToken = domain.StringPtr(token.AccessToken)
		if token.RefreshToken != "" {
			existingOAuthLink.RefreshToken = domain.StringPtr(token.RefreshToken)
		}
		if !token.Expiry.IsZero() {
			existingOAuthLink.ExpiresAt = domain.TimePtr(token.Expiry)
		}
		if err := s.oauthRepo.UpdateOAuthProvider(ctx, existingOAuthLink); err != nil {
			return nil, nil, "", fmt.Errorf("failed to update oauth provider: %w", err)
		}
	} else {
		// Check if user exists by email
		user, err = s.userRepo.GetUserByEmail(ctx, userInfo.Email)
		if err != nil && err != sql.ErrNoRows {
			return nil, nil, "", fmt.Errorf("failed to check user by email: %w", err)
		}

		if user == nil {
			// Create new user (OAuth user without password)
			username := s.generateUniqueUsername(userInfo.Email, userInfo.FirstName)
			user = &domain.User{
				Username:  username,
				Password:  "", // OAuth users don't have password initially
				Email:     userInfo.Email,
				FirstName: userInfo.FirstName,
				LastName:  domain.StringPtr(userInfo.LastName),
				Avatar:    domain.StringPtr(userInfo.Avatar),
				IsActive:  true,
			}

			if err := s.userRepo.RegisterNewUser(ctx, user); err != nil {
				return nil, nil, "", fmt.Errorf("failed to create user: %w", err)
			}
		}

		// Link OAuth provider to user
		oauthLink := &domain.OAuthProviderLink{
			UserID:         user.ID,
			Provider:       provider,
			ProviderUserID: userInfo.ProviderUserID,
			Email:          userInfo.Email,
			AccessToken:    domain.StringPtr(token.AccessToken),
		}
		if token.RefreshToken != "" {
			oauthLink.RefreshToken = domain.StringPtr(token.RefreshToken)
		}
		if !token.Expiry.IsZero() {
			oauthLink.ExpiresAt = domain.TimePtr(token.Expiry)
		}

		if err := s.oauthRepo.CreateOAuthProvider(ctx, oauthLink); err != nil {
			return nil, nil, "", fmt.Errorf("failed to create oauth provider link: %w", err)
		}
	}

	// Get fresh user data with role
	user, err = s.userRepo.GetUserByID(ctx, int(user.ID))
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Generate JWT tokens
	jwtUser := &jwt.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	accessToken, err := s.jwtService.GenerateAccessToken(jwtUser)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshTokenJWT, err := s.jwtService.GenerateRefreshToken(jwtUser)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash the refresh token before storing in database
	hashedToken := helpers.HashToken(refreshTokenJWT)

	// Create refresh token in database
	refreshToken := &domain.RefreshToken{
		UserID:       user.ID,
		RefreshToken: hashedToken, // Use hashed token
		ExpiresAt:    time.Now().Add(s.jwtService.GetRefreshTokenExpiry()),
		CreatedAt:    time.Now(),
		IsRevoked:    false,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
	}

	if err := s.refreshTokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Return the plain JWT to be sent to client (not the hash)
	refreshToken.RefreshToken = refreshTokenJWT

	return user, refreshToken, accessToken, nil
}

// LinkOAuthProvider links an OAuth provider to an existing authenticated user
func (s *oauthService) LinkOAuthProvider(ctx context.Context, userID uint, provider string, code string) error {
	if !domain.IsValidProvider(provider) {
		return fmt.Errorf("invalid provider: %s", provider)
	}

	// Check if provider is already linked
	existing, err := s.oauthRepo.GetOAuthProviderByUserAndProvider(ctx, userID, provider)
	if err != nil {
		return fmt.Errorf("failed to check existing provider link: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("provider already linked to this account")
	}

	// Exchange code for token
	token, err := s.exchangeCodeForToken(ctx, provider, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Fetch user info from provider
	userInfo, err := s.fetchUserInfo(ctx, provider, token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Check if this OAuth account is already linked to another user
	existingOAuthLink, err := s.oauthRepo.GetOAuthProviderByProviderAndUserID(ctx, provider, userInfo.ProviderUserID)
	if err != nil {
		return fmt.Errorf("failed to check oauth provider: %w", err)
	}
	if existingOAuthLink != nil && existingOAuthLink.UserID != userID {
		return fmt.Errorf("this %s account is already linked to another user", provider)
	}

	// Create OAuth provider link
	oauthLink := &domain.OAuthProviderLink{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: userInfo.ProviderUserID,
		Email:          userInfo.Email,
		AccessToken:    domain.StringPtr(token.AccessToken),
	}
	if token.RefreshToken != "" {
		oauthLink.RefreshToken = domain.StringPtr(token.RefreshToken)
	}
	if !token.Expiry.IsZero() {
		oauthLink.ExpiresAt = domain.TimePtr(token.Expiry)
	}

	if err := s.oauthRepo.CreateOAuthProvider(ctx, oauthLink); err != nil {
		return fmt.Errorf("failed to create oauth provider link: %w", err)
	}

	return nil
}

// UnlinkOAuthProvider unlinks an OAuth provider from a user
func (s *oauthService) UnlinkOAuthProvider(ctx context.Context, userID uint, provider string) error {
	if !domain.IsValidProvider(provider) {
		return fmt.Errorf("invalid provider: %s", provider)
	}

	// Get user to check if they have a password
	user, err := s.userRepo.GetUserByID(ctx, int(userID))
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Get all linked providers
	linkedProviders, err := s.oauthRepo.GetOAuthProvidersByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get linked providers: %w", err)
	}

	// Prevent unlinking if it's the only authentication method
	if !user.HasPassword() && len(linkedProviders) <= 1 {
		return fmt.Errorf("cannot unlink the only authentication method. Please set a password first")
	}

	if err := s.oauthRepo.UnlinkProvider(ctx, userID, provider); err != nil {
		return fmt.Errorf("failed to unlink provider: %w", err)
	}

	return nil
}

// GetLinkedProviders retrieves all linked OAuth providers for a user
func (s *oauthService) GetLinkedProviders(ctx context.Context, userID uint) ([]*domain.OAuthProviderLink, error) {
	return s.oauthRepo.GetOAuthProvidersByUserID(ctx, userID)
}

// exchangeCodeForToken exchanges an authorization code for an access token
func (s *oauthService) exchangeCodeForToken(ctx context.Context, provider string, code string) (*oauth2.Token, error) {
	var config *oauth2.Config
	switch provider {
	case string(domain.ProviderGoogle):
		config = s.googleConfig
	case string(domain.ProviderGithub):
		config = s.githubConfig
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// fetchUserInfo fetches user information from the OAuth provider
func (s *oauthService) fetchUserInfo(ctx context.Context, provider string, accessToken string) (*domain.OAuthUserInfo, error) {
	switch provider {
	case string(domain.ProviderGoogle):
		return s.fetchGoogleUserInfo(ctx, accessToken)
	case string(domain.ProviderGithub):
		return s.fetchGithubUserInfo(ctx, accessToken)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// fetchGoogleUserInfo fetches user info from Google
func (s *oauthService) fetchGoogleUserInfo(ctx context.Context, accessToken string) (*domain.OAuthUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google API error: %s", string(body))
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &domain.OAuthUserInfo{
		ProviderUserID: googleUser.ID,
		Email:          googleUser.Email,
		FirstName:      googleUser.GivenName,
		LastName:       googleUser.FamilyName,
		Avatar:         googleUser.Picture,
		Provider:       domain.ProviderGoogle,
	}, nil
}

// fetchGithubUserInfo fetches user info from GitHub
func (s *oauthService) fetchGithubUserInfo(ctx context.Context, accessToken string) (*domain.OAuthUserInfo, error) {
	// Get user profile
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error: %s", string(body))
	}

	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// If email is not public, fetch it from the emails endpoint
	email := githubUser.Email
	if email == "" {
		email, err = s.fetchGithubEmail(accessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch github email: %w", err)
		}
	}

	// Split name into first and last name
	firstName := githubUser.Name
	lastName := ""
	if githubUser.Name != "" {
		parts := strings.SplitN(githubUser.Name, " ", 2)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = parts[1]
		}
	} else {
		firstName = githubUser.Login
	}

	return &domain.OAuthUserInfo{
		ProviderUserID: fmt.Sprintf("%d", githubUser.ID),
		Email:          email,
		FirstName:      firstName,
		LastName:       lastName,
		Avatar:         githubUser.AvatarURL,
		Provider:       domain.ProviderGithub,
	}, nil
}

// fetchGithubEmail fetches the primary email from GitHub
func (s *oauthService) fetchGithubEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github API error: %s", string(body))
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode emails: %w", err)
	}

	// Find primary verified email
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// Fallback to first verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}

// generateUniqueUsername generates a unique username from email or name
func (s *oauthService) generateUniqueUsername(email string, firstName string) string {
	// Try email prefix first
	username := strings.Split(email, "@")[0]
	username = strings.ReplaceAll(username, ".", "_")
	username = strings.ReplaceAll(username, "+", "_")

	// If username is too short, use firstName
	if len(username) < 3 && firstName != "" {
		username = strings.ToLower(firstName)
	}

	// Add random suffix to ensure uniqueness (basic approach)
	// In production, you might want to check database for uniqueness
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d", username, timestamp%10000)
}
