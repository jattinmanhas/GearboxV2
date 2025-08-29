package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAuthService is a mock implementation of IAuthService for testing
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

func (m *MockAuthService) ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*services.RefreshTokenClaims, error) {
	args := m.Called(ctx, refreshTokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.RefreshTokenClaims), args.Error(1)
}

func (m *MockAuthService) GenerateAccessTokenFromUser(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(string), args.Error(1)
}

// TestAuthMiddleware tests the authentication middleware
func TestAuthMiddleware(t *testing.T) {
	// 🎯 Test Strategy: Test authentication middleware with mocked auth service

	t.Run("should allow access with valid access token in cookie", func(t *testing.T) {
		// 🔧 Setup: Create mock auth service and middleware
		mockAuthService := &MockAuthService{}
		middleware := AuthMiddleware(mockAuthService)

		// Create test claims (no user object needed)
		claims := &services.Claims{UserID: 1, Username: "testuser", Email: "test@example.com"}

		// 🎭 Mock Expectations: Auth service should validate token (no user DB query needed)
		mockAuthService.On("ValidateAccessToken", mock.Anything, "valid-token").Return(claims, nil)

		// Create request with valid token in cookie
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid-token",
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler that checks context
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			// Check if claims are in context (no user object needed)
			ctxClaims := GetClaimsFromContext(r.Context())

			assert.Equal(t, claims, ctxClaims)

			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})

	t.Run("should allow access with valid access token in Authorization header", func(t *testing.T) {
		// 🔧 Setup: Create mock auth service and middleware
		mockAuthService := &MockAuthService{}
		middleware := AuthMiddleware(mockAuthService)

		// Create test claims (no user object needed)
		claims := &services.Claims{UserID: 1, Username: "testuser", Email: "test@example.com"}

		// 🎭 Mock Expectations: Auth service should validate token (no user DB query needed)
		mockAuthService.On("ValidateAccessToken", mock.Anything, "valid-token").Return(claims, nil)

		// Create request with valid token in Authorization header
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid-token")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler that checks context
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			// Check if claims are in context (no user object needed)
			ctxClaims := GetClaimsFromContext(r.Context())

			assert.Equal(t, claims, ctxClaims)

			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})

	t.Run("should reject access without token", func(t *testing.T) {
		// 🔧 Setup: Create mock auth service and middleware
		mockAuthService := &MockAuthService{}
		middleware := AuthMiddleware(mockAuthService)

		// Create request without token
		req := httptest.NewRequest("GET", "/protected", nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, handlerCalled)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})

	t.Run("should reject access with invalid token", func(t *testing.T) {
		// 🔧 Setup: Create mock auth service and middleware
		mockAuthService := &MockAuthService{}
		middleware := AuthMiddleware(mockAuthService)

		// 🎭 Mock Expectations: Auth service should reject invalid token
		mockAuthService.On("ValidateAccessToken", mock.Anything, "invalid-token").Return(nil, assert.AnError)

		// Create request with invalid token
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "invalid-token",
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, handlerCalled)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})
}

// TestOptionalAuthMiddleware tests the optional authentication middleware
func TestOptionalAuthMiddleware(t *testing.T) {
	// 🎯 Test Strategy: Test optional authentication middleware with mocked auth service

	t.Run("should allow access with valid token", func(t *testing.T) {
		// 🔧 Setup: Create mock auth service and middleware
		mockAuthService := &MockAuthService{}
		middleware := OptionalAuthMiddleware(mockAuthService)

		// Create test claims (no user object needed)
		claims := &services.Claims{UserID: 1, Username: "testuser", Email: "test@example.com"}

		// 🎭 Mock Expectations: Auth service should validate token (no user DB query needed)
		mockAuthService.On("ValidateAccessToken", mock.Anything, "valid-token").Return(claims, nil)

		// Create request with valid token
		req := httptest.NewRequest("GET", "/optional", nil)
		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "valid-token",
		})

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler that checks context
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			// Check if claims are in context (no user object needed)
			ctxClaims := GetClaimsFromContext(r.Context())

			assert.Equal(t, claims, ctxClaims)

			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})

	t.Run("should allow access without token", func(t *testing.T) {
		// 🔧 Setup: Create mock auth service and middleware
		mockAuthService := &MockAuthService{}
		middleware := OptionalAuthMiddleware(mockAuthService)

		// Create request without token
		req := httptest.NewRequest("GET", "/optional", nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)

		// Verify service was called correctly
		mockAuthService.AssertExpectations(t)
	})
}

// TestCORSMiddleware tests the CORS middleware
func TestCORSMiddleware(t *testing.T) {
	// 🎯 Test Strategy: Test CORS middleware functionality

	t.Run("should handle preflight request", func(t *testing.T) {
		// 🔧 Setup: Create CORS middleware
		middleware := CORSMiddleware([]string{"*"})

		// Create preflight request
		req := httptest.NewRequest("OPTIONS", "/api", nil)
		req.Header.Set("Origin", "https://example.com")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should handle preflight
		assert.Equal(t, http.StatusOK, w.Code)
		assert.False(t, handlerCalled) // Handler should not be called for OPTIONS

		// Check CORS headers
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization, X-Requested-With", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("should allow specific origins", func(t *testing.T) {
		// 🔧 Setup: Create CORS middleware with specific origins
		middleware := CORSMiddleware([]string{"https://example.com", "https://app.example.com"})

		// Create request with allowed origin
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Set("Origin", "https://example.com")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)

		// Check CORS headers
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("should reject disallowed origins", func(t *testing.T) {
		// 🔧 Setup: Create CORS middleware with specific origins
		middleware := CORSMiddleware([]string{"https://example.com"})

		// Create request with disallowed origin
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Set("Origin", "https://malicious.com")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed but without CORS headers
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)

		// Check that no CORS origin header is set
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// TestRequireRole tests the role requirement middleware
func TestRequireRole(t *testing.T) {
	// 🎯 Test Strategy: Test role requirement middleware

	t.Run("should allow authenticated user", func(t *testing.T) {
		// 🔧 Setup: Create role middleware
		middleware := RequireRole("admin")

		// Create request with claims in context
		req := httptest.NewRequest("GET", "/admin", nil)
		ctx := context.WithValue(req.Context(), ClaimsContextKey, &services.Claims{UserID: 1, Username: "admin", Role: "admin"})
		req = req.WithContext(ctx)

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should succeed
		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, handlerCalled)
	})

	t.Run("should reject unauthenticated user", func(t *testing.T) {
		// 🔧 Setup: Create role middleware
		middleware := RequireRole("admin")

		// Create request without user in context
		req := httptest.NewRequest("GET", "/admin", nil)

		// Create response recorder
		w := httptest.NewRecorder()

		// Create handler
		handlerCalled := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
			w.WriteHeader(http.StatusOK)
		})

		// 🚀 Action: Call middleware
		middleware(handler).ServeHTTP(w, req)

		// ✅ Assertions: Should fail
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.False(t, handlerCalled)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "authentication required", response["message"])
	})
}

// TestContextHelpers tests the context helper functions
func TestContextHelpers(t *testing.T) {
	// 🎯 Test Strategy: Test context helper functions

	t.Run("should extract user from context", func(t *testing.T) {
		// 🔧 Setup: Create context with user
		user := &domain.User{ID: 1, Username: "testuser"}
		ctx := context.WithValue(context.Background(), UserContextKey, user)

		// 🚀 Action: Extract user
		extractedUser := GetUserFromContext(ctx)

		// ✅ Assertions: Should return user
		assert.Equal(t, user, extractedUser)
	})

	t.Run("should extract claims from context", func(t *testing.T) {
		// 🔧 Setup: Create context with claims
		claims := &services.Claims{UserID: 1, Username: "testuser"}
		ctx := context.WithValue(context.Background(), ClaimsContextKey, claims)

		// 🚀 Action: Extract claims
		extractedClaims := GetClaimsFromContext(ctx)

		// ✅ Assertions: Should return claims
		assert.Equal(t, claims, extractedClaims)
	})

	t.Run("should return nil for missing values", func(t *testing.T) {
		// 🔧 Setup: Create empty context
		ctx := context.Background()

		// 🚀 Action: Extract values
		user := GetUserFromContext(ctx)
		claims := GetClaimsFromContext(ctx)

		// ✅ Assertions: Should return nil
		assert.Nil(t, user)
		assert.Nil(t, claims)
	})
}
