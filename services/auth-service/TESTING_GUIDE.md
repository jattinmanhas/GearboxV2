# 🧪 **Complete Testing Guide & Learning Resource**

## 🎯 **What This Guide Teaches You**

This document explains **every testing pattern and technique** used in your project and teaches you:
- **Go testing fundamentals**
- **Unit testing strategies**
- **Mocking and dependency injection**
- **HTTP testing techniques**
- **Database testing with mocks**
- **Best practices and patterns**
- **And much more!**

---

## 🏗️ **1. TESTING ARCHITECTURE OVERVIEW**

Your project follows a **layered testing approach** that mirrors your clean architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    TESTING LAYERS                          │
├─────────────────────────────────────────────────────────────┤
│  🧪 Validation Tests    │ Test input validation rules     │
│  🧪 Domain Tests        │ Test data models & structs      │
│  🧪 Repository Tests    │ Test data access (with mocks)   │
│  🧪 Service Tests       │ Test business logic (with mocks)│
│  🧪 Handler Tests       │ Test HTTP handling (with mocks) │
│  🧪 HTTP Response Tests │ Test response formatting        │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔧 **2. TESTING TOOLS & DEPENDENCIES**

### **Core Testing Libraries**
```go
import (
    "testing"           // Go's built-in testing package
    "github.com/stretchr/testify/assert"  // Cleaner assertions
    "github.com/stretchr/testify/require" // Fail-fast assertions
    "github.com/stretchr/testify/mock"    // Mocking framework
    "github.com/DATA-DOG/go-sqlmock"      // SQL database mocking
    "net/http/httptest"                   // HTTP testing utilities
)
```

### **Why These Tools?**
- **`testing`**: Go's standard testing package
- **`testify/assert`**: More readable assertions than `if err != nil`
- **`testify/require`**: Stops test immediately on failure
- **`testify/mock`**: Easy mocking of interfaces
- **`go-sqlmock`**: Test database operations without real database
- **`httptest`**: Test HTTP handlers without starting a server

---

## 🎭 **3. MOCKING PATTERNS**

### **What is Mocking?**
Mocking is creating **fake implementations** of dependencies to:
- **Control behavior** in tests
- **Avoid external dependencies** (databases, APIs)
- **Test error scenarios** easily
- **Speed up tests** (no real network calls)

### **Example: Mock Repository**
```go
// MockUserRepository is a mock implementation of IUserRepository
type MockUserRepository struct {
    mock.Mock  // Embeds testify/mock functionality
}

// RegisterNewUser mocks the RegisterNewUser method
func (m *MockUserRepository) RegisterNewUser(ctx context.Context, u *domain.User) error {
    args := m.Called(ctx, u)  // Record the call
    return args.Error(0)      // Return the mocked error
}

// In your test:
mockRepo := &MockUserRepository{}
mockRepo.On("RegisterNewUser", mock.Anything, mock.Anything).Return(nil)
```

### **Mock Expectations**
```go
// Set up what the mock should expect
mockRepo.On("GetUserByID", mock.Anything, 1).Return(expectedUser, nil)

// Verify expectations were met
mockRepo.AssertExpectations(t)
```

---

## 🧪 **4. TESTING PATTERNS BY LAYER**

### **A. Configuration Layer Testing**

**File: `internal/config/config_test.go`**

```go
func TestLoadConfig(t *testing.T) {
    t.Run("should load config successfully", func(t *testing.T) {
        // 🔧 Setup: Set environment variables
        os.Setenv("PORT", "8080")
        defer os.Unsetenv("PORT")  // Cleanup
        
        // 🚀 Action: Load configuration
        cfg := LoadConfig()
        
        // ✅ Assertions: Verify config values
        assert.Equal(t, "8080", cfg.Port)
    })
}
```

**🎓 Learning Points:**
- **Environment Variable Testing**: Test different config scenarios
- **Cleanup with defer**: Always restore environment after tests
- **Singleton Pattern Testing**: Verify same instance is returned
- **Thread Safety Testing**: Test concurrent access

---

### **B. Domain Model Testing**

**File: `internal/domain/user_test.go`**

```go
func TestUserJSONMarshaling(t *testing.T) {
    t.Run("should hide password in JSON", func(t *testing.T) {
        // 🔧 Setup: Create user with password
        user := &User{
            Username: "john_doe",
            Password: "secret123",  // This should be hidden
        }
        
        // 🚀 Action: Marshal to JSON
        jsonData, err := json.Marshal(user)
        require.NoError(t, err)
        
        // ✅ Assertions: Password should not appear
        jsonString := string(jsonData)
        assert.False(t, contains(jsonString, "secret123"))
    })
}
```

**🎓 Learning Points:**
- **Struct Tag Testing**: Verify `json:"-"` works correctly
- **JSON Marshaling**: Test serialization/deserialization
- **Security Testing**: Ensure sensitive data is hidden
- **Field Mapping**: Test JSON field name conventions

---

### **C. Repository Layer Testing**

**File: `internal/repository/auth_repository_test.go`**

```go
func TestUserRepository_RegisterNewUser(t *testing.T) {
    t.Run("should register new user successfully", func(t *testing.T) {
        // 🔧 Setup: Create mock database
        mockDB, mock, err := sqlmock.New()
        require.NoError(t, err)
        defer mockDB.Close()
        
        // 🎭 Mock Expectations: Set up SQL query expectations
        mock.ExpectQuery("INSERT INTO users").
            WithArgs("john_doe", "hashedpass", "john@example.com").
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
        
        // 🚀 Action: Register user
        repo := NewUserRepository(mockDB)
        user := &domain.User{Username: "john_doe", Email: "john@example.com"}
        err = repo.RegisterNewUser(context.Background(), user)
        
        // ✅ Assertions: Should succeed and set user ID
        assert.NoError(t, err)
        assert.Equal(t, uint(1), user.ID)
        
        // Verify all expectations were met
        assert.NoError(t, mock.ExpectationsWereMet())
    })
}
```

**🎓 Learning Points:**
- **SQL Mocking**: Test database operations without real DB
- **Query Expectations**: Verify SQL queries are correct
- **Result Mocking**: Mock database results and errors
- **Context Testing**: Test context propagation
- **Cleanup**: Always close mock connections

---

### **D. Service Layer Testing**

**File: `internal/services/auth_service_test.go`**

```go
func TestUserService_RegisterNewUser(t *testing.T) {
    t.Run("should hash password securely", func(t *testing.T) {
        // 🔧 Setup: Create mock repository and service
        mockRepo := &MockUserRepository{}
        service := NewUserService(mockRepo)
        
        user := &domain.User{
            Username: "testuser",
            Password: "MySecurePass123",  // Plain text
        }
        
        // 🎭 Mock Expectations: Repository should succeed
        mockRepo.On("RegisterNewUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
            // Verify password was hashed
            return u.Password != "MySecurePass123" && len(u.Password) > 0
        })).Return(nil)
        
        // 🚀 Action: Register user
        err := service.RegisterNewUser(context.Background(), user)
        
        // ✅ Assertions: Should succeed and hash password
        require.NoError(t, err)
        assert.NotEqual(t, "MySecurePass123", user.Password)
        assert.Contains(t, user.Password, "$2")  // bcrypt hash format
        
        // Verify repository was called
        mockRepo.AssertExpectations(t)
    })
}
```

**🎓 Learning Points:**
- **Business Logic Testing**: Test password hashing, validation
- **Mock Verification**: Ensure mocks are called correctly
- **Security Testing**: Verify passwords are properly hashed
- **Dependency Injection**: Test with mocked dependencies
- **Error Propagation**: Test error scenarios

---

### **E. HTTP Response Testing**

**File: `internal/httpx/response_test.go`**

```go
func TestWriteJSON(t *testing.T) {
    t.Run("should write success response correctly", func(t *testing.T) {
        // 🔧 Setup: Create response recorder
        rr := httptest.NewRecorder()
        
        // 🚀 Action: Write JSON response
        WriteJSON(rr, http.StatusOK, true, "Success", map[string]string{"key": "value"}, nil)
        
        // ✅ Assertions: Verify response
        assert.Equal(t, http.StatusOK, rr.Code)
        assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
        
        // Parse and verify response body
        var response APIResponse
        err := json.Unmarshal(rr.Body.Bytes(), &response)
        require.NoError(t, err)
        
        assert.True(t, response.Success)
        assert.Equal(t, "Success", response.Message)
    })
}
```

**🎓 Learning Points:**
- **HTTP Testing**: Use `httptest.NewRecorder()` for testing
- **Response Verification**: Check status codes, headers, body
- **JSON Parsing**: Test response structure and content
- **Header Testing**: Verify content-type and custom headers
- **Status Code Testing**: Test different HTTP status codes

---

## 🎯 **5. TESTING STRATEGIES & PATTERNS**

### **A. Table-Driven Tests**
```go
func TestValidationRules(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected bool
    }{
        {"valid username", "john_doe", true},
        {"too short", "jo", false},
        {"too long", "very_long_username_that_exceeds_thirty_characters", false},
        {"invalid chars", "user@name", false},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := validateUsername(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### **B. Test Organization with t.Run()**
```go
func TestUserService(t *testing.T) {
    t.Run("user registration", func(t *testing.T) {
        t.Run("should succeed with valid data", func(t *testing.T) {
            // Test implementation
        })
        
        t.Run("should fail with invalid data", func(t *testing.T) {
            // Test implementation
        })
    })
    
    t.Run("user retrieval", func(t *testing.T) {
        // More tests...
    })
}
```

### **C. Setup and Teardown**
```go
func TestWithDatabase(t *testing.T) {
    // 🔧 Setup: Create test database
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)  // 🧹 Cleanup
    
    // Your tests here...
}

func setupTestDatabase(t *testing.T) *sql.DB {
    // Create test database
    return db
}

func cleanupTestDatabase(t *testing.T, db *sql.DB) {
    // Clean up test data
    db.Close()
}
```

---

## 🚀 **6. ADVANCED TESTING TECHNIQUES**

### **A. Context Testing**
```go
func TestContextHandling(t *testing.T) {
    t.Run("should respect context cancellation", func(t *testing.T) {
        // Create a context that's already cancelled
        ctx, cancel := context.WithCancel(context.Background())
        cancel() // Cancel immediately
        
        // Test that your function handles cancellation
        result, err := someFunction(ctx)
        
        // Verify appropriate behavior
        if err != nil {
            assert.Contains(t, err.Error(), "context")
        }
    })
}
```

### **B. Concurrency Testing**
```go
func TestConcurrency(t *testing.T) {
    // Test that your function is thread-safe
    results := make(chan int, 10)
    
    for i := 0; i < 10; i++ {
        go func() {
            result := someFunction()
            results <- result
        }()
    }
    
    // Collect results and verify
    for i := 0; i < 10; i++ {
        result := <-results
        assert.Equal(t, expectedValue, result)
    }
}
```

### **C. Benchmark Testing**
```go
func BenchmarkPasswordHashing(b *testing.B) {
    password := "MySecurePassword123"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## 💡 **7. TESTING BEST PRACTICES**

### **1. Test Naming**
```go
// ❌ Bad
func Test1(t *testing.T)

// ✅ Good
func TestUserService_RegisterNewUser_ShouldHashPassword(t *testing.T)
```

### **2. Test Structure (AAA Pattern)**
```go
func TestSomething(t *testing.T) {
    // 🔧 Arrange/Setup
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)
    
    // 🚀 Act/Action
    result, err := service.DoSomething()
    
    // ✅ Assert/Verify
    assert.NoError(t, err)
    assert.Equal(t, expectedValue, result)
}
```

### **3. Test Independence**
```go
// Each test should be independent
func TestA(t *testing.T) {
    // Don't rely on TestB running first
    // Set up your own test data
}

func TestB(t *testing.T) {
    // Don't rely on TestA running first
    // Set up your own test data
}
```

### **4. Error Testing**
```go
func TestErrorScenarios(t *testing.T) {
    t.Run("should handle database connection error", func(t *testing.T) {
        // Test error scenarios, not just success
        mockRepo.On("GetUser").Return(nil, errors.New("connection failed"))
        
        _, err := service.GetUser()
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "connection failed")
    })
}
```

---

## 🔍 **8. COMMON TESTING PATTERNS**

### **A. Interface Compliance Testing**
```go
func TestInterfaceCompliance(t *testing.T) {
    // Verify struct implements interface
    var _ IUserService = &userService{}
    var _ IUserRepository = &userRepository{}
}
```

### **B. Mock Verification**
```go
func TestMockVerification(t *testing.T) {
    mockRepo := &MockUserRepository{}
    
    // Set expectations
    mockRepo.On("GetUser", 1).Return(&User{ID: 1}, nil)
    
    // Call function
    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)
    
    // Verify result
    assert.NoError(t, err)
    assert.Equal(t, uint(1), user.ID)
    
    // Verify mock was called correctly
    mockRepo.AssertExpectations(t)
}
```

### **C. Edge Case Testing**
```go
func TestEdgeCases(t *testing.T) {
    t.Run("empty string", func(t *testing.T) {
        result := validateUsername("")
        assert.False(t, result)
    })
    
    t.Run("nil pointer", func(t *testing.T) {
        var user *User
        err := validateUser(user)
        assert.Error(t, err)
    })
    
    t.Run("zero values", func(t *testing.T) {
        user := &User{} // All fields are zero values
        err := validateUser(user)
        assert.Error(t, err)
    })
}
```

---

## 🧹 **9. TEST CLEANUP & MAINTENANCE**

### **A. Resource Cleanup**
```go
func TestWithResources(t *testing.T) {
    // Create resources
    db := createTestDB()
    defer db.Close()  // Always cleanup
    
    // Create temporary files
    tmpFile := createTempFile()
    defer os.Remove(tmpFile.Name())  // Clean up files
    
    // Your tests...
}
```

### **B. Test Data Management**
```go
func TestWithTestData(t *testing.T) {
    // Use unique test data to avoid conflicts
    username := fmt.Sprintf("test_user_%d", time.Now().Unix())
    email := fmt.Sprintf("test_%d@example.com", time.Now().Unix())
    
    // Your tests...
}
```

### **C. Test Isolation**
```go
func TestIsolation(t *testing.T) {
    // Each test should be independent
    t.Run("test 1", func(t *testing.T) {
        // Set up test data
        // Run test
        // Clean up
    })
    
    t.Run("test 2", func(t *testing.T) {
        // Set up test data (independent of test 1)
        // Run test
        // Clean up
    })
}
```

---

## 🎯 **10. RUNNING YOUR TESTS**

### **A. Run All Tests**
```bash
go test ./...
```

### **B. Run Tests in Specific Package**
```bash
go test ./internal/services
go test ./internal/validation
```

### **C. Run Tests with Verbose Output**
```bash
go test -v ./internal/services
```

### **D. Run Tests with Coverage**
```bash
go test -cover ./internal/services
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### **E. Run Specific Test**
```bash
go test -run TestUserService_RegisterNewUser ./internal/services
```

---

## 🎓 **11. LEARNING ROADMAP**

### **Beginner Level**
1. ✅ **Basic Testing**: Write simple tests with `testing` package
2. ✅ **Assertions**: Use `testify/assert` for cleaner tests
3. ✅ **Test Organization**: Use `t.Run()` for structured tests

### **Intermediate Level**
1. 🔄 **Mocking**: Learn `testify/mock` for dependency mocking
2. 🔄 **HTTP Testing**: Use `httptest` for handler testing
3. 🔄 **Database Mocking**: Use `go-sqlmock` for repository testing

### **Advanced Level**
1. 🎯 **Integration Testing**: Test multiple components together
2. 🎯 **Performance Testing**: Write benchmarks
3. 🎯 **Test Coverage**: Aim for high test coverage
4. 🎯 **Test-Driven Development**: Write tests before code

---

## 🎉 **CONCLUSION**

You now have a **comprehensive testing foundation** that covers:

- ✅ **Unit Testing**: Test individual functions and methods
- ✅ **Mocking**: Test with fake dependencies
- ✅ **HTTP Testing**: Test your API endpoints
- ✅ **Database Testing**: Test data access without real DB
- ✅ **Best Practices**: Follow Go testing conventions
- ✅ **Patterns**: Use proven testing patterns

### **Next Steps:**
1. **Run all tests** to see them in action
2. **Add more tests** for edge cases and error scenarios
3. **Practice mocking** with different dependencies
4. **Explore integration testing** for end-to-end scenarios
5. **Add benchmarks** for performance-critical code

### **Remember:**
- **Good tests** make your code more reliable
- **Test coverage** helps catch bugs early
- **Testing skills** are essential for professional development
- **Practice** makes perfect - write tests for everything!

Your testing setup is now **production-ready** and follows **industry best practices**. Keep writing tests and your code quality will improve dramatically! 🚀
