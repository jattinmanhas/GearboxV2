package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, username, password, userAgent, ipAddress string) (*domain.User, *domain.RefreshToken, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.User, *domain.RefreshToken, string, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uint) error
	ValidateAccessToken(ctx context.Context, tokenString string) (*Claims, error)
	ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*RefreshTokenClaims, error)
	GetUserFromToken(ctx context.Context, tokenString string) (*domain.User, error)
	GenerateAccessTokenFromUser(ctx context.Context, user *domain.User) (string, error)
}

type authService struct {
	userRepo         repository.IUserRepository
	refreshTokenRepo repository.IRefreshTokenRepository
	roleRepo         repository.IRoleRepository
	jwtService       *JWTService
}

func NewAuthService(userRepo repository.IUserRepository, refreshTokenRepo repository.IRefreshTokenRepository, roleRepo repository.IRoleRepository, jwtService *JWTService) IAuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		roleRepo:         roleRepo,
		jwtService:       jwtService,
	}
}

// Login authenticates a user and generates access and refresh tokens
func (a *authService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*domain.User, *domain.RefreshToken, string, error) {
	// Get user by username
	user, err := a.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate access token (stored in cookie by handler)
	accessToken, err := a.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := a.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Set additional fields for refresh token
	refreshToken.UserAgent = userAgent
	refreshToken.IPAddress = ipAddress

	// Store refresh token in database
	if err := a.refreshTokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return user, refreshToken, accessToken, nil
}

// RefreshToken validates a refresh token and generates new access and refresh tokens
func (a *authService) RefreshToken(ctx context.Context, refreshTokenString string) (*domain.User, *domain.RefreshToken, string, error) {
	// Validate refresh token JWT
	claims, err := a.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, nil, "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get refresh token from database
	dbToken, err := a.refreshTokenRepo.GetRefreshTokenByToken(ctx, refreshTokenString)
	if err != nil {
		return nil, nil, "", fmt.Errorf("refresh token not found or expired: %w", err)
	}

	// Verify token belongs to the same user
	if dbToken.UserID != claims.UserID {
		return nil, nil, "", fmt.Errorf("token mismatch")
	}

	// Get user details
	user, err := a.userRepo.GetUserByID(ctx, int(claims.UserID))
	if err != nil {
		return nil, nil, "", fmt.Errorf("user not found: %w", err)
	}

	// Get user's role information
	role, err := a.roleRepo.GetUserRole(ctx, user.ID)
	if err != nil {
		// If no role found, assign default user role
		user.Role = domain.GetDefaultRole()
		user.RoleID = domain.RoleIDUser
	} else {
		user.Role = role.Name
		user.RoleID = role.ID
	}

	// Revoke old refresh token
	if err := a.refreshTokenRepo.RevokeRefreshToken(ctx, refreshTokenString); err != nil {
		return nil, nil, "", fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new access token (stored in cookie by handler)
	accessToken, err := a.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := a.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Copy user agent and IP from old token
	newRefreshToken.UserAgent = dbToken.UserAgent
	newRefreshToken.IPAddress = dbToken.IPAddress

	// Store new refresh token in database
	if err := a.refreshTokenRepo.CreateRefreshToken(ctx, newRefreshToken); err != nil {
		return nil, nil, "", fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return user, newRefreshToken, accessToken, nil
}

// Logout revokes a specific refresh token
func (a *authService) Logout(ctx context.Context, refreshToken string) error {
	return a.refreshTokenRepo.RevokeRefreshToken(ctx, refreshToken)
}

// LogoutAll revokes all refresh tokens for a user
func (a *authService) LogoutAll(ctx context.Context, userID uint) error {
	return a.refreshTokenRepo.RevokeAllUserTokens(ctx, userID)
}

// ValidateAccessToken validates an access token and returns claims
func (a *authService) ValidateAccessToken(ctx context.Context, tokenString string) (*Claims, error) {
	return a.jwtService.ValidateAccessToken(tokenString)
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (a *authService) ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*RefreshTokenClaims, error) {
	// Validate refresh token JWT
	claims, err := a.jwtService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get refresh token from database to ensure it's not revoked
	dbToken, err := a.refreshTokenRepo.GetRefreshTokenByToken(ctx, refreshTokenString)
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
	return a.jwtService.GenerateAccessToken(user)
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

	// Get user's role
	role, err := a.roleRepo.GetUserRole(ctx, user.ID)
	if err != nil {
		// If no role found, assign default user role
		user.Role = domain.GetDefaultRole()
		user.RoleID = domain.RoleIDUser
	} else {
		user.Role = role.Name
		user.RoleID = role.ID
	}

	return user, nil
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
