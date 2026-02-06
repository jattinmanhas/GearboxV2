package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/helpers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, username, password, userAgent, ipAddress string) (*domain.User, *domain.RefreshToken, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.User, string, error) // Updated: no longer returns new refresh token
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uint) error
	ValidateAccessToken(ctx context.Context, tokenString string) (*jwt.Claims, error)
	ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*jwt.RefreshTokenClaims, error)
	GetUserFromToken(ctx context.Context, tokenString string) (*domain.User, error)
	GenerateAccessTokenFromUser(ctx context.Context, user *domain.User) (string, error)
	CleanupExpiredTokens(ctx context.Context) error
	// Password reset methods
	ForgotPassword(ctx context.Context, email, username string, emailService IEmailService) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	userRepo              repository.IUserRepository
	refreshTokenRepo      repository.IRefreshTokenRepository
	passwordResetRepo     repository.IPasswordResetRepository
	roleRepo              repository.IRoleRepository
	jwtService            *jwt.JWTService
	passwordResetTokenTTL time.Duration
}

func NewAuthService(userRepo repository.IUserRepository, refreshTokenRepo repository.IRefreshTokenRepository, passwordResetRepo repository.IPasswordResetRepository, roleRepo repository.IRoleRepository, jwtService *jwt.JWTService) IAuthService {
	return &authService{
		userRepo:              userRepo,
		refreshTokenRepo:      refreshTokenRepo,
		passwordResetRepo:     passwordResetRepo,
		roleRepo:              roleRepo,
		jwtService:            jwtService,
		passwordResetTokenTTL: 1 * time.Hour, // Password reset tokens expire in 1 hour
	}
}

// Login authenticates a user and generates access and refresh tokens
func (a *authService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*domain.User, *domain.RefreshToken, string, error) {
	// Get user by username
	user, err := a.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Println("Invalid credentials", err)
		return nil, nil, "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("Invalid credentials", err)
		return nil, nil, "", fmt.Errorf("invalid credentials")
	}

	// Convert domain user to shared JWT user
	jwtUser := &jwt.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	// Generate access token (returned in response body; stored in memory on client)
	accessToken, err := a.jwtService.GenerateAccessToken(jwtUser)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshTokenJWT, err := a.jwtService.GenerateRefreshToken(jwtUser)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash the refresh token before storing in database
	hashedToken := helpers.HashToken(refreshTokenJWT)

	// Create refresh token domain object
	refreshToken := &domain.RefreshToken{
		UserID:       user.ID,
		RefreshToken: hashedToken, // Store hashed token
		ExpiresAt:    time.Now().Add(a.jwtService.GetRefreshTokenExpiry()),
		CreatedAt:    time.Now(),
		IsRevoked:    false,
	}

	// Set additional fields for refresh token
	refreshToken.UserAgent = userAgent
	refreshToken.IPAddress = ipAddress

	// Store refresh token in database
	if err := a.refreshTokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Return the plain JWT to be sent to client (not the hash)
	refreshToken.RefreshToken = refreshTokenJWT

	return user, refreshToken, accessToken, nil
}

// RefreshToken validates a refresh token and generates a new access token
// Note: Refresh token is NOT rotated for performance reasons (no DB writes)
func (a *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*domain.User, string, error) {
	// Validate refresh token JWT
	claims, err := a.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Hash the token before looking it up in database
	hashedToken := helpers.HashToken(refreshTokenString)

	// Get refresh token from database using hashed value
	dbToken, err := a.refreshTokenRepo.GetRefreshTokenByToken(ctx, hashedToken)
	if err != nil {
		return nil, "", fmt.Errorf("refresh token not found or expired: %w", err)
	}

	// Verify token belongs to the same user
	if dbToken.UserID != claims.UserID {
		return nil, "", fmt.Errorf("token mismatch")
	}

	// Get user details
	user, err := a.userRepo.GetUserByID(ctx, int(claims.UserID))
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	roleName, ok := domain.RoleNames[int(user.RoleID)]
	if !ok {
		user.RoleID, user.Role = domain.GetDefaultRole()
	} else {
		user.Role = roleName
	}

	// Convert domain user to shared JWT user
	jwtUser := &jwt.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	// Generate new access token only (no refresh token rotation)
	accessToken, err := a.jwtService.GenerateAccessToken(jwtUser)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Return user and new access token (refresh token remains unchanged)
	return user, accessToken, nil
}

// Logout revokes a specific refresh token
func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	// Hash the token before revoking
	hashedToken := helpers.HashToken(refreshToken)
	return a.refreshTokenRepo.RevokeRefreshToken(ctx, hashedToken)
}

// LogoutAll revokes all refresh tokens for a user
func (a *authService) LogoutAll(ctx context.Context, userID uint) error {
	return a.refreshTokenRepo.RevokeAllUserTokens(ctx, userID)
}

// ValidateAccessToken validates an access token and returns claims
func (a *authService) ValidateAccessToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	return a.jwtService.ValidateAccessToken(tokenString)
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (a *authService) ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*jwt.RefreshTokenClaims, error) {
	// Validate refresh token JWT
	claims, err := a.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Hash the token before looking it up in database
	hashedToken := helpers.HashToken(refreshTokenString)

	// Get refresh token from database to ensure it's not revoked
	dbToken, err := a.refreshTokenRepo.GetRefreshTokenByToken(ctx, hashedToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired: %w", err)
	}

	// Verify token belongs to the same user
	if dbToken.UserID != claims.UserID {
		return nil, fmt.Errorf("token mismatch")
	}

	return claims, nil
}

// GenerateAccessTokenFromUser generates a new access token for a user
func (a *authService) GenerateAccessTokenFromUser(ctx context.Context, user *domain.User) (string, error) {
	// Convert domain user to shared JWT user
	jwtUser := &jwt.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
	return a.jwtService.GenerateAccessToken(jwtUser)
}

// GetUserFromToken extracts user information from a valid access token
func (a *authService) GetUserFromToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := a.jwtService.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	user, err := a.userRepo.GetUserByID(ctx, int(claims.UserID))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	roleName, ok := domain.RoleNames[int(user.RoleID)]
	if !ok {
		user.RoleID, user.Role = domain.GetDefaultRole()
	} else {
		user.Role = roleName
	}

	return user, nil
}

// CleanupExpiredTokens revokes all expired refresh tokens
func (a *authService) CleanupExpiredTokens(ctx context.Context) error {
	return a.refreshTokenRepo.CleanupExpiredTokens(ctx)
}

// ExtractTokenFromHeader extracts the token from Authorization header
func ExtractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenFromCookie extracts the token from cookies
func ExtractTokenFromCookie(r *http.Request, cookieName string) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ForgotPassword generates a password reset token and sends it via email
// Accepts either email or username to find the user
func (a *authService) ForgotPassword(ctx context.Context, email, username string, emailService IEmailService) error {
	var user *domain.User
	var err error
	identifier := ""

	// Try to find user by email first (if provided), then by username
	if email != "" {
		identifier = email
		user, err = a.userRepo.GetUserByEmail(ctx, email)
	} else if username != "" {
		identifier = username
		user, err = a.userRepo.GetUserByUsername(ctx, username)
	} else {
		// Neither email nor username provided
		log.Printf("Neither email nor username provided for password reset")
		return nil
	}

	if err != nil {
		// Don't reveal if user exists or not (security best practice)
		// Return success even if user doesn't exist to prevent enumeration
		log.Printf("User not found for identifier: %s", identifier)
		return nil
	}

	// Check if user has a password (OAuth users might not have passwords)
	if !user.HasPassword() {
		log.Printf("User %s does not have a password set (OAuth user)", identifier)
		return nil
	}

	// Ensure user has an email (required for sending reset link)
	if user.Email == "" {
		log.Printf("User %s does not have an email address set", identifier)
		return nil
	}

	// Delete any existing password reset tokens for this user
	if err := a.passwordResetRepo.DeleteTokenByUserID(ctx, user.ID); err != nil {
		log.Printf("Failed to delete existing tokens: %v", err)
		// Continue anyway - we'll create a new token
	}

	// Generate a secure random token
	token, err := generateSecureToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset token
	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(a.passwordResetTokenTTL),
		Used:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store token in database
	if err := a.passwordResetRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	// Send password reset email
	if err := emailService.SendPasswordResetEmail(user.Email, token, ""); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		// Don't return error - token is already created, user can request again if needed
	}

	return nil
}

// ResetPassword resets a user's password using a valid reset token
func (a *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get the password reset token
	resetToken, err := a.passwordResetRepo.GetPasswordResetTokenByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Get the user
	user, err := a.userRepo.GetUserByID(ctx, int(resetToken.UserID))
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user's password
	user.Password = string(hashedPassword)
	if err := a.userRepo.UpdateUser(ctx, int(user.ID), user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := a.passwordResetRepo.MarkTokenAsUsed(ctx, token); err != nil {
		log.Printf("Failed to mark token as used: %v", err)
		// Don't return error - password is already reset
	}

	// Revoke all refresh tokens for this user (security best practice)
	if err := a.refreshTokenRepo.RevokeAllUserTokens(ctx, user.ID); err != nil {
		log.Printf("Failed to revoke refresh tokens: %v", err)
		// Don't return error - password is already reset
	}

	return nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
