package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserService is a mock implementation of IUserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterNewUser(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id int, u *domain.User) (*domain.User, error) {
	args := m.Called(ctx, id, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, id int, currentPassword, newPassword string) error {
	args := m.Called(ctx, id, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAuthService is a mock implementation of IAuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*domain.User, *domain.RefreshToken, string, error) {
	args := m.Called(ctx, username, password, userAgent, ipAddress)
	if args.Get(0) == nil {
		return nil, nil, "", args.Error(3)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*domain.User), nil, "", args.Error(3)
	}
	return args.Get(0).(*domain.User), args.Get(1).(*domain.RefreshToken), args.Get(2).(string), args.Error(3)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.User, *domain.RefreshToken, string, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, nil, "", args.Error(3)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*domain.User), nil, "", args.Error(3)
	}
	return args.Get(0).(*domain.User), args.Get(1).(*domain.RefreshToken), args.Get(2).(string), args.Error(3)
}

func (m *MockAuthService) ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*services.RefreshTokenClaims, error) {
	args := m.Called(ctx, refreshTokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.RefreshTokenClaims), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, tokenString string) (*services.Claims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Claims), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(ctx context.Context, tokenString string) (*domain.User, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) GenerateAccessTokenFromUser(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockAuthService) CleanupExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestAuthHandler_Login tests the login handler
func TestAuthHandler_Login(t *testing.T) {
	// 🎯 Test Strategy: Test login handler with mocked services

	t.Run("should login user successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// Create test user and refresh token
		user := &domain.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}
		refreshToken := &domain.RefreshToken{
			ID:           1,
			UserID:       1,
			RefreshToken: "test-refresh-token",
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
			IsRevoked:    false,
		}

		// 🎭 Mock Expectations: Auth service should handle login
		mockAuthService.On("Login", mock.Anything, "testuser", "password123", "test-agent", "127.0.0.1").
			Return(user, refreshToken, "test-access-token", nil)

		// Create request
		loginReq := loginRequest{
			Username: "testuser",
			Password: "password123",
		}
		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "test-agent")
		req.RemoteAddr = "127.0.0.1:12345"

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call login handler
		handler.Login(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "login successful", response["message"])

		// Check cookies
		cookies := w.Result().Cookies()
		accessTokenCookie := findCookie(cookies, "access_token")
		refreshTokenCookie := findCookie(cookies, "refresh_token")

		assert.NotNil(t, accessTokenCookie)
		assert.NotNil(t, refreshTokenCookie)
		assert.True(t, accessTokenCookie.HttpOnly)
		assert.True(t, refreshTokenCookie.HttpOnly)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})

	t.Run("should fail with invalid request body", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// Create invalid request
		req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call login handler
		handler.Login(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid request body", response["message"])

		// Verify no services were called
		mockAuthService.AssertNotCalled(t, "Login")
	})

	t.Run("should fail with validation errors", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// Create request with missing fields
		loginReq := loginRequest{
			Username: "", // Missing username
			Password: "password123",
		}
		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call login handler
		handler.Login(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "validation failed", response["message"])

		// Verify no services were called
		mockAuthService.AssertNotCalled(t, "Login")
	})

	t.Run("should fail with authentication error", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// 🎭 Mock Expectations: Auth service should return error
		mockAuthService.On("Login", mock.Anything, "testuser", "wrongpassword", "test-agent", "127.0.0.1").
			Return(nil, nil, "", assert.AnError)

		// Create request
		loginReq := loginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "test-agent")
		req.RemoteAddr = "127.0.0.1:12345"

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call login handler
		handler.Login(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "invalid credentials", response["message"])

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})
}

// TestAuthHandler_RefreshToken tests the refresh token handler
func TestAuthHandler_RefreshToken(t *testing.T) {
	// 🎯 Test Strategy: Test refresh token handler with mocked services

	t.Run("should refresh tokens successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// Create test user and refresh token
		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}
		newRefreshToken := &domain.RefreshToken{
			ID:           2,
			UserID:       1,
			RefreshToken: "new-refresh-token",
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
			IsRevoked:    false,
		}

		// 🎭 Mock Expectations: Auth service should handle refresh
		mockAuthService.On("RefreshToken", mock.Anything, "old-refresh-token").
			Return(user, newRefreshToken, "new-access-token", nil)

		// Create request with refresh token cookie
		req := httptest.NewRequest("POST", "/refresh", nil)
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "old-refresh-token",
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call refresh handler
		handler.RefreshToken(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "tokens refreshed successfully", response["message"])

		// Check cookies
		cookies := w.Result().Cookies()
		accessTokenCookie := findCookie(cookies, "access_token")
		refreshTokenCookie := findCookie(cookies, "refresh_token")

		assert.NotNil(t, accessTokenCookie)
		assert.NotNil(t, refreshTokenCookie)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})
}

// TestAuthHandler_Logout tests the logout handler
func TestAuthHandler_Logout(t *testing.T) {
	// 🎯 Test Strategy: Test logout handler with mocked services

	t.Run("should logout user successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// 🎭 Mock Expectations: Auth service should handle logout
		mockAuthService.On("Logout", mock.Anything, "test-refresh-token").Return(nil)

		// Create request with refresh token cookie
		req := httptest.NewRequest("POST", "/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "test-refresh-token",
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call logout handler
		handler.Logout(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "logout successful", response["message"])

		// Check that cookies are cleared
		cookies := w.Result().Cookies()
		accessTokenCookie := findCookie(cookies, "access_token")
		refreshTokenCookie := findCookie(cookies, "refresh_token")

		assert.NotNil(t, accessTokenCookie)
		assert.NotNil(t, refreshTokenCookie)
		assert.Equal(t, "", accessTokenCookie.Value)
		assert.Equal(t, "", refreshTokenCookie.Value)
		assert.Equal(t, -1, accessTokenCookie.MaxAge)
		assert.Equal(t, -1, refreshTokenCookie.MaxAge)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})

	t.Run("should fail without refresh token", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// Create request without refresh token cookie
		req := httptest.NewRequest("POST", "/logout", nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call logout handler
		handler.Logout(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "refresh token not found", response["message"])

		// Verify no services were called
		mockAuthService.AssertNotCalled(t, "Logout")
	})
}

// TestAuthHandler_LogoutAll tests the logout all handler
func TestAuthHandler_LogoutAll(t *testing.T) {
	// 🎯 Test Strategy: Test logout all handler with mocked services

	t.Run("should logout from all devices successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock services and handler
		mockUserService := &MockUserService{}
		mockAuthService := &MockAuthService{}
		jwtService := services.NewJWTService("test-secret", "test-refresh-secret")
		handler := NewAuthHandler(mockUserService, mockAuthService, jwtService)

		// Create test user and access token
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		accessToken, err := jwtService.GenerateAccessToken(user)
		require.NoError(t, err)

		// 🎭 Mock Expectations: Auth service should handle logout
		claims := &services.Claims{UserID: 1, Username: "testuser", Email: "test@example.com"}
		mockAuthService.On("LogoutAll", mock.Anything, uint(1)).Return(nil)

		// Create request with access token cookie
		req := httptest.NewRequest("POST", "/logout-all", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: accessToken,
		})

		// Set claims in context (as middleware would do)
		ctx := context.WithValue(req.Context(), "claims", claims)
		req = req.WithContext(ctx)

		// Create response recorder
		w := httptest.NewRecorder()

		// 🚀 Action: Call logout all handler
		handler.LogoutAll(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)

		// Check response body
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "logout from all devices successful", response["message"])

		// Check that cookies are cleared
		cookies := w.Result().Cookies()
		accessTokenCookie := findCookie(cookies, "access_token")
		refreshTokenCookie := findCookie(cookies, "refresh_token")

		assert.NotNil(t, accessTokenCookie)
		assert.NotNil(t, refreshTokenCookie)
		assert.Equal(t, "", accessTokenCookie.Value)
		assert.Equal(t, "", refreshTokenCookie.Value)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})
}

// Helper function to find a cookie by name
func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
