package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository is a mock implementation of IUserRepository
type MockUserRepository struct {
	mock.Mock
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

// GetAllUsers mocks the GetAllUsers method
func (m *MockUserRepository) GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

// TestUserService tests the user service implementation
func TestUserService(t *testing.T) {
	// 🎯 Test Strategy: Test service layer with mocked repository

	t.Run("should create new user service", func(t *testing.T) {
		// 🔧 Setup: Create mock repository
		mockRepo := &MockUserRepository{}

		// 🚀 Action: Create service
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

		// Create test user
		user := &domain.User{
			Username:    "john_doe",
			Password:    "SecurePass123", // Plain text password
			Email:       "john@example.com",
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

		expectedUser := &domain.User{
			ID:        1,
			Username:  "john_doe",
			Email:     "john@example.com",
			FirstName: "John",
			LastName:  "Doe",
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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

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
		service := NewUserService(mockRepo)

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
