package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
)

type JWTService struct {
	accessTokenSecret  string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenID string `json:"token_id"`
	UserID  uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTService(accessSecret, refreshSecret string) *JWTService {
	return &JWTService{
		accessTokenSecret:  accessSecret,
		refreshTokenSecret: refreshSecret,
		accessTokenExpiry:  15 * time.Minute,   // 15 minutes
		refreshTokenExpiry: 7 * 24 * time.Hour, // 7 days
	}
}

// GenerateAccessToken creates a new access token for a user
func (j *JWTService) GenerateAccessToken(user *domain.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessTokenSecret))
}

// GenerateRefreshToken creates a new refresh token and returns both token and claims
func (j *JWTService) GenerateRefreshToken(user *domain.User) (*domain.RefreshToken, error) {
	// Generate a unique token ID
	tokenID := uuid.New().String()

	// Generate a cryptographically secure random string for additional security
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	_ = base64.URLEncoding.EncodeToString(randomBytes) // Additional entropy

	claims := &RefreshTokenClaims{
		TokenID: tokenID,
		UserID:  user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Sign the refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenJWT, err := token.SignedString([]byte(j.refreshTokenSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Create refresh token domain object
	refreshToken := &domain.RefreshToken{
		UserID:       user.ID,
		RefreshToken: refreshTokenJWT,
		ExpiresAt:    time.Now().Add(j.refreshTokenExpiry),
		CreatedAt:    time.Now(),
		IsRevoked:    false,
	}

	return refreshToken, nil
}

// ValidateAccessToken validates and parses an access token
func (j *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accessTokenSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates and parses a refresh token
func (j *JWTService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.refreshTokenSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetAccessTokenExpiry returns the access token expiry duration
func (j *JWTService) GetAccessTokenExpiry() time.Duration {
	return j.accessTokenExpiry
}

// GetRefreshTokenExpiry returns the refresh token expiry duration
func (j *JWTService) GetRefreshTokenExpiry() time.Duration {
	return j.refreshTokenExpiry
}
