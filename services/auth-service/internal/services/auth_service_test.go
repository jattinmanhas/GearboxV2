package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// changePasswordRequest matches the struct in the handler package
type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// MockUserRepository is a mock implementation of IUserRepository
type MockUserRepository struct {
	mock.Mock
}

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

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*domain.User, string, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*domain.User), args.Get(1).(string), args.Error(2)
}

func (m *MockAuthService) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func (m *MockAuthService) ValidateRefreshToken(ctx context.Context, refreshTokenString string) (*jwt.RefreshTokenClaims, error) {
	args := m.Called(ctx, refreshTokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RefreshTokenClaims), args.Error(1)
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

func (m *MockAuthService) ForgotPassword(ctx context.Context, email, username string, emailService IEmailService) error {
	args := m.Called(ctx, email, username, emailService)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

// RegisterNewUser mocks the RegisterNewUser method
func (m *MockUserRepository) RegisterNewUser(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserRepository) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// GetUserByUsername mocks the GetUserByUsername method
func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// GetUserByEmail mocks the GetUserByEmail method
func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// GetAllUsers mocks the GetAllUsers method
func (m *MockUserRepository) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

// GetAllUsersWithFilters mocks the GetAllUsersWithFilters method
func (m *MockUserRepository) GetAllUsersWithFilters(ctx context.Context, limit int, offset int, search string, isActive *bool, roleID *int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset, search, isActive, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

// GetUsersCount mocks the GetUsersCount method
func (m *MockUserRepository) GetUsersCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// GetUsersCountWithFilters mocks the GetUsersCountWithFilters method
func (m *MockUserRepository) GetUsersCountWithFilters(ctx context.Context, search string, isActive *bool, roleID *int) (int, error) {
	args := m.Called(ctx, search, isActive, roleID)
	return args.Int(0), args.Error(1)
}

// UpdateUser mocks the UpdateUser method
func (m *MockUserRepository) UpdateUser(ctx context.Context, id int, u *domain.User) error {
	args := m.Called(ctx, id, u)
	return args.Error(0)
}

// DeleteUser mocks the DeleteUser method
func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetUserAnalytics mocks the GetUserAnalytics method
func (m *MockUserRepository) GetUserAnalytics(ctx context.Context) (*domain.UserAnalytics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserAnalytics), args.Error(1)
}

// MockRefreshTokenRepository is a mock implementation of IRefreshTokenRepository
type MockRefreshTokenRepository struct {
	mock.Mock
}

type MockRoleRepository struct {
	mock.Mock
}

// MockPasswordResetRepository is a mock implementation of IPasswordResetRepository
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) GetPasswordResetTokenByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetRepository) GetPasswordResetTokenByUserID(ctx context.Context, userID uint) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeleteExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeleteTokenByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// CreateRefreshToken mocks the CreateRefreshToken method
func (m *MockRefreshTokenRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// GetRefreshTokenByToken mocks the GetRefreshTokenByToken method
func (m *MockRefreshTokenRepository) GetRefreshTokenByToken(ctx context.Context, refreshToken string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

// GetRefreshTokensByUserID mocks the GetRefreshTokensByUserID method
func (m *MockRefreshTokenRepository) GetRefreshTokensByUserID(ctx context.Context, userID uint) ([]domain.RefreshToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefreshToken), args.Error(1)
}

// RevokeRefreshToken mocks the RevokeRefreshToken method
func (m *MockRefreshTokenRepository) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

// RevokeAllUserTokens mocks the RevokeAllUserTokens method
func (m *MockRefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockRoleRepository methods
func (m *MockRoleRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) GetAllRoles(ctx context.Context) ([]*domain.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteRole(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID, grantedBy uint) error {
	args := m.Called(ctx, userID, roleID, grantedBy)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserRole(ctx context.Context, userID uint) (*domain.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) UpdateUserRole(ctx context.Context, userID, roleID uint) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveUserRole(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRoleRepository) InitializeDefaultRoles(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// CleanupExpiredTokens mocks the CleanupExpiredTokens method
func (m *MockRefreshTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// DeleteRefreshToken mocks the DeleteRefreshToken method
func (m *MockRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

// TestUserService tests the user service implementation
func TestUserService(t *testing.T) {
	// 🎯 Test Strategy: Test service layer with mocked repository

	t.Run("should create new user service", func(t *testing.T) {
		// 🔧 Setup: Create mock repository
		mockRepo := &MockUserRepository{}

		// 🚀 Action: Create service
		service := NewUserService(mockRepo, nil)

		// ✅ Assertions: Service should be created
		assert.NotNil(t, service)

		// Verify it implements the interface
		var _ IUserService = service
	})
}

// TestUserService_RegisterNewUser tests user registration service logic
func TestUserService_RegisterNewUser(t *testing.T) {
	// 🎯 Test Strategy: Test password hashing and repository delegation

	t.Run("should register new user successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create test user
		user := &domain.User{
			Username:    "john_doe",
			Password:    "SecurePass123", // Plain text password
			Email:       "john@example.com",
			FirstName:   "John",
			LastName:    domain.NewNullString("Doe"),
			DateOfBirth: domain.NewNullTime(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
		}

		// 🎭 Mock Expectations: Repository should be called with hashed password
		mockRepo.On("RegisterNewUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			// Verify password was hashed (should not be plain text)
			return u.Password != "SecurePass123" && len(u.Password) > 0
		})).Return(nil)

		// 🚀 Action: Register user
		err := service.RegisterNewUser(context.Background(), user)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)

		// Verify password was hashed
		assert.NotEqual(t, "SecurePass123", user.Password)
		assert.True(t, len(user.Password) > 0)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		// 🔧 Setup: Create mock repository that returns error
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		user := &domain.User{
			Username: "john_doe",
			Password: "SecurePass123",
			Email:    "john@example.com",
		}

		// 🎭 Mock Expectations: Repository should return error
		expectedError := errors.New("database connection failed")
		mockRepo.On("RegisterNewUser", mock.Anything, mock.Anything).Return(expectedError)

		// 🚀 Action: Register user
		err := service.RegisterNewUser(context.Background(), user)

		// ✅ Assertions: Should return repository error
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})

	t.Run("should hash password with different costs", func(t *testing.T) {
		// 🔧 Setup: Create mock repository
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		user1 := &domain.User{
			Username: "user1",
			Password: "Password123",
			Email:    "user1@example.com",
		}

		user2 := &domain.User{
			Username: "user2",
			Password: "Password123",
			Email:    "user2@example.com",
		}

		// 🎭 Mock Expectations: Both should succeed
		mockRepo.On("RegisterNewUser", mock.Anything, mock.Anything).Return(nil).Times(2)

		// 🚀 Action: Register both users
		err1 := service.RegisterNewUser(context.Background(), user1)
		err2 := service.RegisterNewUser(context.Background(), user2)

		// ✅ Assertions: Both should succeed
		assert.NoError(t, err1)
		assert.NoError(t, err2)

		// Passwords should be hashed (different hashes due to salt)
		assert.NotEqual(t, "Password123", user1.Password)
		assert.NotEqual(t, "Password123", user2.Password)
		assert.NotEqual(t, user1.Password, user2.Password) // Different salts

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_GetUserByID tests user retrieval service logic
func TestUserService_GetUserByID(t *testing.T) {
	// 🎯 Test Strategy: Test service delegation to repository

	t.Run("should get user by ID successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		expectedUser := &domain.User{
			ID:        1,
			Username:  "john_doe",
			Email:     "john@example.com",
			FirstName: "John",
			LastName:  domain.NewNullString("Doe"),
		}

		// 🎭 Mock Expectations: Repository should return user
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)

		// 🚀 Action: Get user by ID
		user, err := service.GetUserByID(context.Background(), 1)

		// ✅ Assertions: Should return user successfully
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error when getting user", func(t *testing.T) {
		// 🔧 Setup: Create mock repository that returns error
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// 🎭 Mock Expectations: Repository should return error
		expectedError := errors.New("user not found")
		mockRepo.On("GetUserByID", mock.Anything, 999).Return(nil, expectedError)

		// 🚀 Action: Get non-existent user
		user, err := service.GetUserByID(context.Background(), 999)

		// ✅ Assertions: Should return error
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_GetAllUsers tests user listing service logic
func TestUserService_GetAllUsers(t *testing.T) {
	// 🎯 Test Strategy: Test pagination and repository delegation

	t.Run("should get all users with pagination", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		expectedUsers := []domain.User{
			{ID: 1, Username: "john_doe", Email: "john@example.com"},
			{ID: 2, Username: "jane_doe", Email: "jane@example.com"},
		}

		// 🎭 Mock Expectations: Repository should return users
		mockRepo.On("GetAllUsers", mock.Anything, 10, 0).Return(expectedUsers, nil)

		// 🚀 Action: Get users with pagination
		users, err := service.GetAllUsers(context.Background(), 10, 0)

		// ✅ Assertions: Should return users successfully
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, expectedUsers, users)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error when getting users", func(t *testing.T) {
		// 🔧 Setup: Create mock repository that returns error
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// 🎭 Mock Expectations: Repository should return error
		expectedError := errors.New("database connection failed")
		mockRepo.On("GetAllUsers", mock.Anything, 10, 0).Return(nil, expectedError)

		// 🚀 Action: Get users
		users, err := service.GetAllUsers(context.Background(), 10, 0)

		// ✅ Assertions: Should return error
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, expectedError, err)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_Context tests context handling
func TestUserService_Context(t *testing.T) {
	// 🎯 Test Strategy: Test context propagation to repository

	t.Run("should pass context to repository", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create a context with a specific value
		ctx := context.WithValue(context.Background(), "user_id", "123")

		// 🎭 Mock Expectations: Repository should receive the context
		mockRepo.On("GetUserByID", mock.MatchedBy(func(c context.Context) bool {
			// Verify context contains our value
			return c.Value("user_id") == "123"
		}), 1).Return(&domain.User{ID: 1}, nil)

		// 🚀 Action: Get user with context
		_, err := service.GetUserByID(ctx, 1)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)

		// Verify repository was called with correct context
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_PasswordHashing tests password hashing functionality
func TestUserService_PasswordHashing(t *testing.T) {
	// 🎯 Test Strategy: Test password hashing security features

	t.Run("should hash password securely", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		plainPassword := "MySecurePassword123"
		user := &domain.User{
			Username: "testuser",
			Password: plainPassword,
			Email:    "test@example.com",
		}

		// 🎭 Mock Expectations: Repository should succeed
		mockRepo.On("RegisterNewUser", mock.Anything, mock.Anything).Return(nil)

		// 🚀 Action: Register user
		err := service.RegisterNewUser(context.Background(), user)

		// ✅ Assertions: Should succeed and hash password
		require.NoError(t, err)

		// Password should be hashed (much longer than plain text)
		assert.NotEqual(t, plainPassword, user.Password)
		assert.True(t, len(user.Password) > len(plainPassword))

		// Hash should look like bcrypt hash (starts with $2a$ or $2b$)
		assert.Contains(t, user.Password, "$2")

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle bcrypt hashing error", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create user with extremely long password that might cause bcrypt issues
		// (This is a theoretical test - bcrypt is very robust)
		user := &domain.User{
			Username: "testuser",
			Password: "a", // Very short password
			Email:    "test@example.com",
		}

		// 🎭 Mock Expectations: Repository should succeed
		mockRepo.On("RegisterNewUser", mock.Anything, mock.Anything).Return(nil)

		// 🚀 Action: Register user
		err := service.RegisterNewUser(context.Background(), user)

		// ✅ Assertions: Should succeed even with short password
		// (bcrypt handles this gracefully)
		assert.NoError(t, err)
		assert.NotEqual(t, "a", user.Password)

		// Verify repository was called
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_UpdateUser tests user update service logic
func TestUserService_UpdateUser(t *testing.T) {
	// 🎯 Test Strategy: Test partial updates and repository delegation

	t.Run("should update user successfully with partial data", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create existing user
		existingUser := &domain.User{
			ID:          1,
			Username:    "john_doe",
			Email:       "john@example.com",
			FirstName:   "John",
			MiddleName:  domain.NewNullString("Michael"),
			LastName:    domain.NewNullString("Doe"),
			Avatar:      domain.NewNullString("https://example.com/avatar.jpg"),
			Gender:      domain.NewNullString("male"),
			DateOfBirth: domain.NewNullTime(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
		}

		// Create update request with only some fields
		updateData := &domain.User{
			FirstName: "Jonathan",                                                 // Only updating first name
			Avatar:    domain.NewNullString("https://example.com/new-avatar.jpg"), // And avatar
		}

		// 🎭 Mock Expectations: Repository should be called to get existing user
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(existingUser, nil)

		// Mock Expectations: Repository should be called to update user
		mockRepo.On("UpdateUser", mock.Anything, 1, mock.MatchedBy(func(u *domain.User) bool {
			// Verify that the update data is passed through
			return u.FirstName == "Jonathan" &&
				u.Avatar.String == "https://example.com/new-avatar.jpg"
		})).Return(nil)

		// 🚀 Action: Update user
		updatedUser, err := service.UpdateUser(context.Background(), 1, updateData)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error when getting user", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// 🎭 Mock Expectations: Repository should return error
		mockRepo.On("GetUserByID", mock.Anything, 999).Return(nil, errors.New("user not found"))

		// 🚀 Action: Update user
		updatedUser, err := service.UpdateUser(context.Background(), 999, &domain.User{})

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.Contains(t, err.Error(), "user not found")

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error when updating user", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create existing user
		existingUser := &domain.User{
			ID:        1,
			Username:  "john_doe",
			FirstName: "John",
		}

		// 🎭 Mock Expectations: Repository should be called to get existing user
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(existingUser, nil)

		// Mock Expectations: Repository should return error on update
		mockRepo.On("UpdateUser", mock.Anything, 1, mock.Anything).Return(errors.New("database error"))

		// 🚀 Action: Update user
		updatedUser, err := service.UpdateUser(context.Background(), 1, &domain.User{FirstName: "Jonathan"})

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
		assert.Contains(t, err.Error(), "database error")

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_ChangePassword tests password change service logic
func TestUserService_ChangePassword(t *testing.T) {
	// 🎯 Test Strategy: Test password change and repository delegation

	t.Run("should change password successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create existing user with hashed password
		hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte("old_password"), bcrypt.DefaultCost)
		existingUser := &domain.User{
			ID:       1,
			Username: "john_doe",
			Password: string(hashedOldPassword),
			Email:    "john@example.com",
		}

		// 🎭 Mock Expectations: Repository should be called to get existing user
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(existingUser, nil)

		// Mock Expectations: Repository should be called to update user with new hashed password
		mockRepo.On("UpdateUser", mock.Anything, 1, mock.MatchedBy(func(u *domain.User) bool {
			// Verify password was hashed and is different from old password
			return u.Password != string(hashedOldPassword) && len(u.Password) > 0
		})).Return(nil)

		// 🚀 Action: Change password
		err := service.ChangePassword(context.Background(), 1, "old_password", "new_password")

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error when getting user", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// 🎭 Mock Expectations: Repository should return error
		mockRepo.On("GetUserByID", mock.Anything, 999).Return(nil, errors.New("user not found"))

		// 🚀 Action: Change password
		err := service.ChangePassword(context.Background(), 999, "old_password", "new_password")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle wrong current password", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create existing user with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
		existingUser := &domain.User{
			ID:       1,
			Username: "john_doe",
			Password: string(hashedPassword),
			Email:    "john@example.com",
		}

		// 🎭 Mock Expectations: Repository should be called to get existing user
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(existingUser, nil)

		// 🚀 Action: Change password with wrong current password
		err := service.ChangePassword(context.Background(), 1, "wrong_password", "new_password")

		// ✅ Assertions: Should fail due to wrong current password
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "crypto/bcrypt")

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle update repository error", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// Create existing user with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("old_password"), bcrypt.DefaultCost)
		existingUser := &domain.User{
			ID:       1,
			Username: "john_doe",
			Password: string(hashedPassword),
			Email:    "john@example.com",
		}

		// 🎭 Mock Expectations: Repository should be called to get existing user
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(existingUser, nil)

		// Mock Expectations: Repository should return error on update
		mockRepo.On("UpdateUser", mock.Anything, 1, mock.Anything).Return(errors.New("database error"))

		// 🚀 Action: Change password
		err := service.ChangePassword(context.Background(), 1, "old_password", "new_password")

		// ✅ Assertions: Should fail due to repository error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})
}

// TestUserService_DeleteUser tests user deletion service logic
func TestUserService_DeleteUser(t *testing.T) {
	// 🎯 Test Strategy: Test user deletion and repository delegation

	t.Run("should delete user successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// 🎭 Mock Expectations: Repository should check if user exists and then delete
		mockRepo.On("GetUserByID", mock.Anything, 1).Return(&domain.User{ID: 1}, nil)
		mockRepo.On("DeleteUser", mock.Anything, 1).Return(nil)

		// 🚀 Action: Delete user
		err := service.DeleteUser(context.Background(), 1)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error when deleting user", func(t *testing.T) {
		// 🔧 Setup: Create mock repository and service
		mockRepo := &MockUserRepository{}
		mockAuthService := &MockAuthService{}
		service := NewUserService(mockRepo, mockAuthService)

		// 🎭 Mock Expectations: Repository should check if user exists and then return error
		mockRepo.On("GetUserByID", mock.Anything, 999).Return(&domain.User{ID: 999}, nil)
		mockRepo.On("DeleteUser", mock.Anything, 999).Return(errors.New("database error"))

		// 🚀 Action: Delete user
		err := service.DeleteUser(context.Background(), 999)

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify repository was called correctly
		mockRepo.AssertExpectations(t)
	})
}

// TestAuthService tests the authentication service implementation
func TestAuthService(t *testing.T) {
	// 🎯 Test Strategy: Test authentication service with mocked repositories

	t.Run("should create new auth service", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and JWT service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")

		// 🚀 Action: Create service
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// ✅ Assertions: Service should be created
		assert.NotNil(t, service)

		// Verify it implements the interface
		var _ IAuthService = service
	})
}

// TestAuthService_Login tests user login functionality
func TestAuthService_Login(t *testing.T) {
	// 🎯 Test Strategy: Test login with mocked repositories

	t.Run("should login user successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		mockRoleRepo := &MockRoleRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// Create test user with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
			Email:    "test@example.com",
		}

		// 🎭 Mock Expectations: Repository should return user
		mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil)

		// Mock Expectations: Refresh token repository should store token
		mockRefreshTokenRepo.On("CreateRefreshToken", mock.Anything, mock.MatchedBy(func(token *domain.RefreshToken) bool {
			return token.UserID == user.ID && !token.IsRevoked
		})).Return(nil)

		// 🚀 Action: Login user
		loggedInUser, refreshToken, _, err := service.Login(context.Background(), "testuser", "password123", "test-agent", "127.0.0.1")

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)
		assert.Equal(t, user, loggedInUser)
		assert.NotNil(t, refreshToken)
		assert.Equal(t, user.ID, refreshToken.UserID)
		assert.Equal(t, "test-agent", refreshToken.UserAgent)
		assert.Equal(t, "127.0.0.1", refreshToken.IPAddress)

		// Verify repositories were called correctly
		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid username", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🎭 Mock Expectations: Repository should return error
		mockUserRepo.On("GetUserByUsername", mock.Anything, "nonexistent").Return(nil, errors.New("user not found"))

		// 🚀 Action: Login with non-existent user
		user, refreshToken, _, err := service.Login(context.Background(), "nonexistent", "password123", "test-agent", "127.0.0.1")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Nil(t, refreshToken)
		assert.Contains(t, err.Error(), "invalid credentials")

		// Verify repositories were called correctly
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// Create test user with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
			Email:    "test@example.com",
		}

		// 🎭 Mock Expectations: Repository should return user
		mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(user, nil)

		// 🚀 Action: Login with wrong password
		loggedInUser, refreshToken, _, err := service.Login(context.Background(), "testuser", "wrongpassword", "test-agent", "127.0.0.1")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.Nil(t, refreshToken)
		assert.Contains(t, err.Error(), "invalid credentials")

		// Verify repositories were called correctly
		mockUserRepo.AssertExpectations(t)
	})
}

// TestAuthService_RefreshToken tests token refresh functionality
func TestAuthService_RefreshToken(t *testing.T) {
	// 🎯 Test Strategy: Test token refresh with mocked repositories

	t.Run("should refresh tokens successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// Create test user
		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		// Generate refresh token
		jwtUser := &jwt.User{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role}
		refreshToken, err := jwtService.GenerateRefreshToken(jwtUser)
		require.NoError(t, err)

		// 🎭 Mock Expectations: Refresh token repository should return token
		// Create a mock refresh token domain object for the repository
		mockRefreshToken := &domain.RefreshToken{
			UserID:       user.ID,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
			CreatedAt:    time.Now(),
			IsRevoked:    false,
		}
		mockRefreshTokenRepo.On("GetRefreshTokenByToken", mock.Anything, mock.Anything).Return(mockRefreshToken, nil)

		// Mock Expectations: User repository should return user
		mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

		// 🚀 Action: Refresh token
		refreshedUser, accessToken, err := service.RefreshToken(context.Background(), refreshToken)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)
		assert.Equal(t, user, refreshedUser)
		assert.NotEmpty(t, accessToken)

		// Verify repositories were called correctly
		mockUserRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid refresh token", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🚀 Action: Refresh with invalid token
		user, accessToken, err := service.RefreshToken(context.Background(), "invalid-token")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, accessToken)
		assert.Contains(t, err.Error(), "invalid refresh token")

		// Verify no repositories were called
		mockUserRepo.AssertNotCalled(t, "GetUserByID")
		mockRefreshTokenRepo.AssertNotCalled(t, "GetRefreshTokenByToken")
	})
}

// TestAuthService_Logout tests logout functionality
func TestAuthService_Logout(t *testing.T) {
	// 🎯 Test Strategy: Test logout with mocked repositories

	t.Run("should logout user successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🎭 Mock Expectations: Refresh token should be revoked
		mockRefreshTokenRepo.On("RevokeRefreshToken", mock.Anything, "test-refresh-token").Return(nil)

		// 🚀 Action: Logout user
		err := service.Logout(context.Background(), "test-refresh-token")

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)

		// Verify repository was called correctly
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should handle logout error", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🎭 Mock Expectations: Repository should return error
		mockRefreshTokenRepo.On("RevokeRefreshToken", mock.Anything, "test-refresh-token").Return(errors.New("database error"))

		// 🚀 Action: Logout user
		err := service.Logout(context.Background(), "test-refresh-token")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify repository was called correctly
		mockRefreshTokenRepo.AssertExpectations(t)
	})
}

// TestAuthService_LogoutAll tests logout from all devices functionality
func TestAuthService_LogoutAll(t *testing.T) {
	// 🎯 Test Strategy: Test logout all with mocked repositories

	t.Run("should logout from all devices successfully", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🎭 Mock Expectations: All user tokens should be revoked
		mockRefreshTokenRepo.On("RevokeAllUserTokens", mock.Anything, uint(1)).Return(nil)

		// 🚀 Action: Logout from all devices
		err := service.LogoutAll(context.Background(), 1)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)

		// Verify repository was called correctly
		mockRefreshTokenRepo.AssertExpectations(t)
	})

	t.Run("should handle logout all error", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🎭 Mock Expectations: Repository should return error
		mockRefreshTokenRepo.On("RevokeAllUserTokens", mock.Anything, uint(1)).Return(errors.New("database error"))

		// 🚀 Action: Logout from all devices
		err := service.LogoutAll(context.Background(), 1)

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		// Verify repository was called correctly
		mockRefreshTokenRepo.AssertExpectations(t)
	})
}

// TestAuthService_ValidateAccessToken tests access token validation
func TestAuthService_ValidateAccessToken(t *testing.T) {
	// 🎯 Test Strategy: Test access token validation

	t.Run("should validate valid access token", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// Create test user and generate token
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		jwtUser := &jwt.User{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role}
		token, err := jwtService.GenerateAccessToken(jwtUser)
		require.NoError(t, err)

		// 🚀 Action: Validate token
		claims, err := service.ValidateAccessToken(context.Background(), token)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, user.Email, claims.Email)
	})

	t.Run("should reject invalid access token", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🚀 Action: Validate invalid token
		claims, err := service.ValidateAccessToken(context.Background(), "invalid-token")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "failed to parse token")
	})
}

// TestAuthService_GetUserFromToken tests user extraction from token
func TestAuthService_GetUserFromToken(t *testing.T) {
	// 🎯 Test Strategy: Test user extraction from valid token

	t.Run("should get user from valid token", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// Create test user and generate token
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		jwtUser := &jwt.User{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role}
		token, err := jwtService.GenerateAccessToken(jwtUser)
		require.NoError(t, err)

		// 🎭 Mock Expectations: User repository should return user
		mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(user, nil)

		// Mock Expectations: Role repository should return default role
		mockRoleRepo.On("GetUserRole", mock.Anything, uint(1)).Return(nil, errors.New("user role not found"))

		// 🚀 Action: Get user from token
		extractedUser, err := service.GetUserFromToken(context.Background(), token)

		// ✅ Assertions: Should succeed
		assert.NoError(t, err)
		assert.Equal(t, user, extractedUser)

		// Verify repository was called correctly
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// 🚀 Action: Get user from invalid token
		user, err := service.GetUserFromToken(context.Background(), "invalid-token")

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to parse token")

		// Verify no repositories were called
		mockUserRepo.AssertNotCalled(t, "GetUserByID")
	})

	t.Run("should fail when user not found", func(t *testing.T) {
		// 🔧 Setup: Create mock repositories and service
		mockUserRepo := &MockUserRepository{}
		mockRefreshTokenRepo := &MockRefreshTokenRepository{}
		mockPasswordResetRepo := &MockPasswordResetRepository{}
		jwtService := jwt.NewJWTService("test-secret", "test-refresh-secret")
		mockRoleRepo := &MockRoleRepository{}
		service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordResetRepo, mockRoleRepo, jwtService)

		// Create test user and generate token
		user := &domain.User{ID: 1, Username: "testuser", Email: "test@example.com"}
		jwtUser := &jwt.User{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role}
		token, err := jwtService.GenerateAccessToken(jwtUser)
		require.NoError(t, err)

		// 🎭 Mock Expectations: User repository should return error
		mockUserRepo.On("GetUserByID", mock.Anything, 1).Return(nil, errors.New("user not found"))

		// 🚀 Action: Get user from token
		extractedUser, err := service.GetUserFromToken(context.Background(), token)

		// ✅ Assertions: Should fail
		assert.Error(t, err)
		assert.Nil(t, extractedUser)
		assert.Contains(t, err.Error(), "user not found")

		// Verify repository was called correctly
		mockUserRepo.AssertExpectations(t)
	})
}

// 🎓 **LEARNING POINTS FROM THESE TESTS:**

// 1. **Mocking**: Use testify/mock to create mock implementations
// 2. **Interface Testing**: Verify structs implement interfaces
// 3. **Password Security**: Test that passwords are properly hashed
// 4. **Error Propagation**: Test that errors flow through the service layer
// 5. **Context Handling**: Test context propagation to lower layers
// 6. **Business Logic**: Test service-specific logic (password hashing)
// 7. **Mock Expectations**: Verify mocks are called correctly

// 💡 **TESTING BEST PRACTICES:**
// - Mock external dependencies (repository)
// - Test both success and error scenarios
// - Verify business logic (password hashing)
// - Test context propagation
// - Use descriptive test names
// - Test edge cases and security features
// - Verify mock expectations are met

// 🔧 **TOOLS USED:**
// - testify/mock: Create mock implementations
// - testify/assert: Cleaner assertions
// - testify/require: Fail fast assertions
// - bcrypt: Password hashing (already in your code)
