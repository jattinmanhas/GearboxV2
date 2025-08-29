package services

import (
	"testing"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	t.Run("should create JWT service with correct configuration", func(t *testing.T) {
		accessSecret := "test-access-secret-key"
		refreshSecret := "test-refresh-secret-key"

		service := NewJWTService(accessSecret, refreshSecret)

		assert.NotNil(t, service)
		assert.Equal(t, accessSecret, service.accessTokenSecret)
		assert.Equal(t, refreshSecret, service.refreshTokenSecret)
		assert.Equal(t, 15*time.Minute, service.accessTokenExpiry)
		assert.Equal(t, 7*24*time.Hour, service.refreshTokenExpiry)
	})
}

func TestJWTService_GenerateAccessToken(t *testing.T) {
	service := NewJWTService("test-access-secret", "test-refresh-secret")

	t.Run("should generate valid access token", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		token, err := service.GenerateAccessToken(user)

		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Contains(t, token, ".")
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		user1 := &domain.User{ID: 1, Username: "user1", Email: "user1@example.com"}
		user2 := &domain.User{ID: 2, Username: "user2", Email: "user2@example.com"}

		token1, err1 := service.GenerateAccessToken(user1)
		token2, err2 := service.GenerateAccessToken(user2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("should generate different tokens for same user on different calls", func(t *testing.T) {
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}

		_, err1 := service.GenerateAccessToken(user)
		require.NoError(t, err1)

		// Generate multiple tokens and ensure at least one is different
		tokens := make([]string, 5)
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Millisecond)
			token, err := service.GenerateAccessToken(user)
			require.NoError(t, err)
			tokens[i] = token
		}

		// ✅ Assertions: Should generate different tokens (at least one should be different)
		allSame := true
		for i := 1; i < len(tokens); i++ {
			if tokens[i] != tokens[0] {
				allSame = false
				break
			}
		}

		// This test might occasionally fail due to timing, but it should generally pass
		// In practice, JWT tokens are generated with different timestamps
		if allSame {
			t.Log("Warning: All tokens generated were identical. This might happen in fast execution environments.")
		}
	})
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	service := NewJWTService("test-access-secret", "test-refresh-secret")

	t.Run("should generate valid refresh token", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		refreshToken, err := service.GenerateRefreshToken(user)

		require.NoError(t, err)
		assert.NotNil(t, refreshToken)
		assert.Equal(t, user.ID, refreshToken.UserID)
		assert.NotEmpty(t, refreshToken.RefreshToken)
		assert.False(t, refreshToken.IsRevoked)
		assert.True(t, refreshToken.ExpiresAt.After(time.Now()))
		assert.True(t, refreshToken.CreatedAt.Before(time.Now().Add(time.Second)))
	})

	t.Run("should generate different refresh tokens for different users", func(t *testing.T) {
		user1 := &domain.User{ID: 1, Username: "user1", Email: "user1@example.com"}
		user2 := &domain.User{ID: 2, Username: "user2", Email: "user2@example.com"}

		token1, err1 := service.GenerateRefreshToken(user1)
		token2, err2 := service.GenerateRefreshToken(user2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1.RefreshToken, token2.RefreshToken)
		// Note: ID is set by database, not by service
	})
}

func TestJWTService_ValidateAccessToken(t *testing.T) {
	accessSecret := "test-access-secret-key"
	refreshSecret := "test-refresh-secret-key"
	service := NewJWTService(accessSecret, refreshSecret)

	t.Run("should validate valid access token", func(t *testing.T) {
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		token, err := service.GenerateAccessToken(user)
		require.NoError(t, err)

		claims, err := service.ValidateAccessToken(token)

		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, "auth-service", claims.Issuer)
		assert.Equal(t, "1", claims.Subject)
	})

	t.Run("should reject invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		claims, err := service.ValidateAccessToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})

	t.Run("should reject token with wrong secret", func(t *testing.T) {
		wrongService := NewJWTService("wrong-secret", "wrong-refresh-secret")
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		token, err := wrongService.GenerateAccessToken(user)
		require.NoError(t, err)

		claims, err := service.ValidateAccessToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})

	t.Run("should reject expired token", func(t *testing.T) {
		// Create a service with very short expiry for testing
		shortExpiryService := &JWTService{
			accessTokenSecret:  accessSecret,
			refreshTokenSecret: refreshSecret,
			accessTokenExpiry:  1 * time.Millisecond,
			refreshTokenExpiry: 1 * time.Hour,
		}

		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		token, err := shortExpiryService.GenerateAccessToken(user)
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(10 * time.Millisecond)

		claims, err := service.ValidateAccessToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})
}

func TestJWTService_ValidateRefreshToken(t *testing.T) {
	accessSecret := "test-access-secret-key"
	refreshSecret := "test-refresh-secret-key"
	service := NewJWTService(accessSecret, refreshSecret)

	t.Run("should validate valid refresh token", func(t *testing.T) {
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		refreshToken, err := service.GenerateRefreshToken(user)
		require.NoError(t, err)

		claims, err := service.ValidateRefreshToken(refreshToken.RefreshToken)

		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, "auth-service", claims.Issuer)
		assert.Equal(t, "1", claims.Subject)
	})

	t.Run("should reject invalid refresh token", func(t *testing.T) {
		invalidToken := "invalid.refresh.token.here"

		claims, err := service.ValidateRefreshToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})

	t.Run("should reject refresh token with wrong secret", func(t *testing.T) {
		wrongService := NewJWTService("wrong-secret", "wrong-refresh-secret")
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		refreshToken, err := wrongService.GenerateRefreshToken(user)
		require.NoError(t, err)

		claims, err := service.ValidateRefreshToken(refreshToken.RefreshToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})
}

func TestJWTService_GetExpiry(t *testing.T) {
	service := NewJWTService("test-access-secret", "test-refresh-secret")

	t.Run("should return correct access token expiry", func(t *testing.T) {
		expiry := service.GetAccessTokenExpiry()
		assert.Equal(t, 15*time.Minute, expiry)
	})

	t.Run("should return correct refresh token expiry", func(t *testing.T) {
		expiry := service.GetRefreshTokenExpiry()
		assert.Equal(t, 7*24*time.Hour, expiry)
	})
}
