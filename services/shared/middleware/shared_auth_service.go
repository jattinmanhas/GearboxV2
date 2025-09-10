package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
)

// SharedAuthService provides an implementation of IAuthService using the shared JWT service
type SharedAuthService struct {
	jwtService *jwt.JWTService
}

// NewSharedAuthService creates a new SharedAuthService instance
func NewSharedAuthService(jwtService *jwt.JWTService) *SharedAuthService {
	return &SharedAuthService{
		jwtService: jwtService,
	}
}

// ValidateAccessToken validates an access token and returns claims
func (s *SharedAuthService) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	// Convert shared JWT claims to middleware claims
	return &Claims{
		UserID:           claims.UserID,
		Username:         claims.Username,
		Email:            claims.Email,
		Role:             claims.Role,
		RegisteredClaims: claims.RegisteredClaims,
	}, nil
}

// ValidateRefreshToken validates a refresh token and returns claims
func (s *SharedAuthService) ValidateRefreshToken(ctx context.Context, token string) (*Claims, error) {
	claims, err := s.jwtService.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	// Convert shared JWT claims to middleware claims
	return &Claims{
		UserID:           claims.UserID,
		Username:         claims.Username,
		Email:            claims.Email,
		Role:             claims.Role,
		RegisteredClaims: claims.RegisteredClaims,
	}, nil
}

// GenerateAccessTokenFromUser generates a new access token from user information
func (s *SharedAuthService) GenerateAccessTokenFromUser(ctx context.Context, user interface{}) (string, error) {
	// Type assert user to our User struct
	u, ok := user.(*User)
	if !ok {
		return "", fmt.Errorf("invalid user type")
	}

	// Convert middleware User to shared JWT User
	jwtUser := &jwt.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}

	return s.jwtService.GenerateAccessToken(jwtUser)
}

// ExtractTokenFromHeader extracts token from Authorization header
func (s *SharedAuthService) ExtractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if it's a Bearer token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return ""
	}

	return tokenString
}

// ExtractTokenFromCookie extracts token from cookie
func (s *SharedAuthService) ExtractTokenFromCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}
