# 🚀 **Complete Code Explanation & Learning Guide**

## 🎯 **What This Document Teaches You**

This document explains **every line of your existing code** and teaches you:
- **Go patterns and best practices**
- **Clean Architecture principles**
- **HTTP handling and routing**
- **Database operations**
- **Interface design**
- **Error handling**
- **And much more!**

---

## 🏗️ **1. MAIN FUNCTION - The Entry Point**

### **File: `cmd/api/main.go`**

```go
func main() {
    cfg := config.LoadConfig()
    log.Println("Configuration loaded successfully.")

    // Connect to DB
    database, err := db.NewConnection(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("❌ Failed to connect to the database: %v", err)
    }
    defer database.Close()
    log.Println("✅ Database connection established successfully.")

    // Initialize repository
    userRepo := repository.NewUserRepository(database)

    // Initialize service
    userService := services.NewUserService(userRepo)

    // Initialize handler
    authHandler := handlers.NewAuthHandler(userService)

    // Initialize router
    r := router.NewRouter(authHandler)

    log.Println("🚀 Auth Service Running on port:", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
```

### **🔍 What's Happening Here?**

This is a **Dependency Injection** pattern! Let me break it down:

1. **Load Configuration** → Get database URL and port
2. **Connect to Database** → Establish PostgreSQL connection
3. **Create Repository** → Data access layer
4. **Create Service** → Business logic layer
5. **Create Handler** → HTTP request handling layer
6. **Create Router** → URL routing layer
7. **Start Server** → Listen for HTTP requests

### **🎓 Learning Points:**

- **`defer database.Close()`**: Ensures database connection is closed when program exits
- **`log.Fatal()`**: Logs error and exits program (for critical failures)
- **Dependency Chain**: Each layer depends on the layer below it

---

## ⚙️ **2. CONFIGURATION - Environment & Settings**

### **File: `internal/config/config.go`**

```go
type Config struct {
    Port        string
    DatabaseURL string
}

var (
    cfg  *Config
    once sync.Once // Ensures the config is loaded only once.
)

func LoadConfig() *Config {
    once.Do(func() {
        if err := godotenv.Load(); err != nil {
            log.Fatal("Error loading .env file")
        }

        port := os.Getenv("PORT")
        if port == "" {
            port = "8081"  // Default port
        }

        databaseURL := os.Getenv("DATABASE_URL")
        if databaseURL == "" {
            log.Fatal("Error: DATABASE_URL environment variable is not set.")
        }

        cfg = &Config{
            Port:        port,
            DatabaseURL: databaseURL,
        }
    })
    return cfg
}
```

### **🔍 What's Happening Here?**

This implements the **Singleton Pattern** with **Thread Safety**:

1. **`sync.Once`**: Ensures config is loaded only once, even if called from multiple goroutines
2. **Environment Variables**: Gets settings from `.env` file or system environment
3. **Default Values**: Provides sensible defaults (like port 8081)
4. **Error Handling**: Fails fast if critical config is missing

### **🎓 Learning Points:**

- **`sync.Once`**: Go's way of ensuring something happens only once
- **Environment Variables**: Secure way to store configuration (not in code)
- **Fail Fast**: If critical config is missing, stop immediately
- **Default Values**: Always provide sensible defaults

---

## 🗄️ **3. DATABASE CONNECTION - PostgreSQL Setup**

### **File: `internal/db/db.go`**

```go
import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func NewConnection(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

### **🔍 What's Happening Here?**

This uses **SQLx** (enhanced SQL library) and **PostgreSQL driver**:

1. **`_ "github.com/lib/pq"`**: The underscore means "import for side effects" (registers PostgreSQL driver)
2. **`sqlx.Connect`**: Creates connection and validates it works
3. **`dsn`**: "Data Source Name" - connection string like `postgres://user:pass@localhost/dbname`

### **🎓 Learning Points:**

- **Blank Import (`_`)**: Import package just to run its `init()` function
- **SQLx**: Enhanced version of Go's `database/sql` with better struct scanning
- **Connection String**: Standard format for database connections

---

## 🏛️ **4. DOMAIN MODELS - Your Data Structures**

### **File: `internal/domain/user.go`**

```go
type User struct {
    ID          uint      `json:"id" db:"id"`
    Username    string    `json:"username" db:"username"`
    Password    string    `json:"-" db:"password"`
    Email       string    `json:"email" db:"email"`
    FirstName   string    `json:"first_name" db:"first_name"`
    MiddleName  string    `json:"middle_name" db:"middle_name"`
    LastName    string    `json:"last_name" db:"last_name"`
    Avatar      string    `json:"avatar" db:"avatar"`
    Gender      string    `json:"gender" db:"gender"`
    DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
}
```

### **🔍 What's Happening Here?**

This defines your **data model** with **struct tags**:

1. **`json:"id"`**: When converting to JSON, use "id" as the key
2. **`json:"-"`**: Exclude this field from JSON (password is hidden!)
3. **`db:"username"`**: When scanning database results, map to this field
4. **`uint`**: Unsigned integer (can't be negative)
5. **`time.Time`**: Go's built-in time type

### **🎓 Learning Points:**

- **Struct Tags**: Metadata attached to struct fields
- **JSON Marshaling**: Automatic conversion between Go structs and JSON
- **Database Mapping**: Automatic mapping between database columns and struct fields
- **Security**: Password field is hidden from JSON responses

---

## 🗃️ **5. REPOSITORY LAYER - Data Access**

### **File: `internal/repository/auth_repository.go`**

```go
type IUserRepository interface {
    RegisterNewUser(ctx context.Context, u *domain.User) error
    GetUserByID(ctx context.Context, id int) (*domain.User, error)
    GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
}

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) RegisterNewUser(ctx context.Context, u *domain.User) error {
    query := `
        INSERT INTO users (
            username, password, email, first_name, middle_name, last_name, avatar, gender, date_of_birth
        ) VALUES (
            :username, :password, :email, :first_name, :middle_name, :last_name, :avatar, :gender, :date_of_birth
        ) RETURNING id;
    `

    // use NamedQueryRowx to bind struct fields by db tags
    rows, err := r.db.NamedQueryContext(ctx, query, u)
    if err != nil {
        return err
    }
    defer rows.Close()
    if rows.Next() {
        if err := rows.Scan(&u.ID); err != nil {
            return err
        }
    }
    return rows.Err()
}
```

### **🔍 What's Happening Here?**

This implements the **Repository Pattern** with **SQLx named queries**:

1. **Interface**: Defines what methods a repository must have
2. **Struct**: Implements the interface with actual database operations
3. **Named Queries**: `:username` gets replaced with `u.Username` value
4. **Context**: Allows cancellation and timeouts
5. **`RETURNING id`**: PostgreSQL returns the generated ID

### **🎓 Learning Points:**

- **Interface Design**: Define behavior, not implementation
- **Named Queries**: Safer than string concatenation (prevents SQL injection)
- **Context**: Go's way of handling cancellation, timeouts, and request-scoped values
- **Error Handling**: Always check and return errors

---

## 🧠 **6. SERVICE LAYER - Business Logic**

### **File: `internal/services/auth_service.go`**

```go
type IUserService interface {
    RegisterNewUser(ctx context.Context, u *domain.User) error
    GetUserByID(ctx context.Context, id int) (*domain.User, error)
    GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
}

type userService struct {
    userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
    return &userService{userRepo: userRepo}
}

func (s *userService) RegisterNewUser(ctx context.Context, u *domain.User) error {
    // Hash the password
    hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    u.Password = string(hash)

    return s.userRepo.RegisterNewUser(ctx, u)
}
```

### **🔍 What's Happening Here?**

This implements the **Service Layer Pattern** with **Password Hashing**:

1. **Interface**: Defines business logic contract
2. **Dependency Injection**: Service receives repository as parameter
3. **Password Hashing**: Uses bcrypt for secure password storage
4. **Business Rules**: Enforces password security before saving

### **🎓 Learning Points:**

- **Service Layer**: Where business logic lives (not in handlers or repositories)
- **Dependency Injection**: Pass dependencies instead of creating them
- **bcrypt**: Industry-standard password hashing (slow, secure)
- **Separation of Concerns**: Each layer has one responsibility

---

## 🌐 **7. HANDLER LAYER - HTTP Request Handling**

### **File: `internal/handlers/auth_handler.go`**

```go
type registerRequest struct {
    Username    string     `json:"username" validate:"required,username"`
    Password    string     `json:"password" validate:"required,password"`
    Email       string     `json:"email" validate:"required,email"`
    FirstName   string     `json:"first_name" validate:"required,min=1,max=50"`
    MiddleName  string     `json:"middle_name" validate:"omitempty,max=50"`
    LastName    string     `json:"last_name" validate:"required,min=1,max=50"`
    Avatar      string     `json:"avatar" validate:"omitempty,url"`
    Gender      string     `json:"gender" validate:"omitempty,oneof=male female other prefer_not_to_say"`
    DateOfBirth *time.Time `json:"date_of_birth" validate:"required,date_of_birth"`
}

func (h *authHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
        return
    }

    // Validate the request using our validation package
    if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
        httpx.Error(w, http.StatusBadRequest, "validation failed", validationErrors)
        return
    }

    user := &domain.User{
        Username:    strings.TrimSpace(req.Username),
        Password:    req.Password,
        Email:       strings.TrimSpace(req.Email),
        FirstName:   strings.TrimSpace(req.FirstName),
        MiddleName:  strings.TrimSpace(req.MiddleName),
        LastName:    strings.TrimSpace(req.LastName),
        Avatar:      req.Avatar,
        Gender:      req.Gender,
        DateOfBirth: *req.DateOfBirth,
    }

    if err := h.userService.RegisterNewUser(r.Context(), user); err != nil {
        httpx.Error(w, http.StatusInternalServerError, "failed to register user", err)
        return
    }

    httpx.Created(w, "user registered", map[string]any{"id": user.ID, "username": user.Username, "email": user.Email})
}
```

### **🔍 What's Happening Here?**

This handles **HTTP requests** with **validation** and **data transformation**:

1. **Request Struct**: Separate struct for incoming data (different from domain model)
2. **JSON Decoding**: Converts HTTP body to Go struct
3. **Validation**: Checks data before processing
4. **Data Cleaning**: `strings.TrimSpace()` removes whitespace
5. **Service Call**: Delegates business logic to service layer
6. **Response**: Returns success/error with appropriate HTTP status

### **🎓 Learning Points:**

- **Request/Response Models**: Separate from domain models
- **Input Sanitization**: Always clean user input
- **HTTP Status Codes**: Use appropriate status codes (200, 201, 400, 500)
- **Error Handling**: Return meaningful error messages
- **Context**: Pass request context through the call chain

---

## 🛣️ **8. ROUTER LAYER - URL Routing**

### **File: `internal/router/router.go`**

```go
func NewRouter(authHandler handlers.IAuthHandler) *chi.Mux {
    router := chi.NewRouter()

    // Auth routes
    router.Route("/api/v1/auth", func(r chi.Router) {
        r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("Auth Service is running"))
        })

        // User management routes
        r.Post("/register", authHandler.RegisterUser)
        r.Get("/user/{id}", authHandler.GetUserByID)
        r.Put("/user/{id}", authHandler.UpdateUser)
        r.Delete("/user/{id}", authHandler.DeleteUser)
        r.Post("/user/{id}/change-password", authHandler.ChangePassword)
        r.Get("/users", authHandler.GetAllUsers)
    })

    return router
}
```

### **🔍 What's Happening Here?**

This uses **Chi Router** for **RESTful API routing**:

1. **Chi Router**: Lightweight, fast HTTP router
2. **Route Groups**: `/api/v1/auth` groups related endpoints
3. **HTTP Methods**: GET, POST, PUT, DELETE for different operations
4. **URL Parameters**: `{id}` captures dynamic values
5. **Handler Mapping**: Each route maps to a handler function

### **🎓 Learning Points:**

- **RESTful Design**: Use HTTP methods appropriately
- **URL Structure**: `/api/v1/auth/users` follows conventions
- **Route Groups**: Organize related endpoints together
- **URL Parameters**: Dynamic values in URLs

---

## 📡 **9. HTTP RESPONSE HELPERS - Standardized Responses**

### **File: `internal/httpx/response.go`**

```go
type APIResponse struct {
    Timestamp time.Time   `json:"timestamp"`
    Status    int         `json:"status"`
    Success   bool        `json:"success"`
    Message   string      `json:"message,omitempty"`
    Data      interface{} `json:"data,omitempty"`
    Error     interface{} `json:"error,omitempty"`
}

func OK(w http.ResponseWriter, message string, data interface{}) {
    WriteJSON(w, http.StatusOK, true, message, data, nil)
}

func Error(w http.ResponseWriter, status int, message string, err error) {
    payload := map[string]string{"message": message}
    if err != nil {
        payload["detail"] = err.Error()
    }
    WriteJSON(w, status, false, message, nil, payload)
}
```

### **🔍 What's Happening Here?**

This provides **consistent API responses** across your service:

1. **Standardized Format**: All responses have the same structure
2. **Helper Functions**: Easy to use functions for common responses
3. **JSON Marshaling**: Automatic conversion to JSON
4. **Error Details**: Include error information when available

### **🎓 Learning Points:**

- **Consistent APIs**: Users know what to expect
- **Helper Functions**: Reduce code duplication
- **JSON Structure**: Well-defined response format
- **Error Handling**: Include useful error information

---

## 🔄 **10. COMPLETE REQUEST FLOW**

Let me show you how a request flows through your entire system:

```
HTTP Request: POST /api/v1/auth/register
     ↓
Router: Matches route to authHandler.RegisterUser
     ↓
Handler: Decodes JSON, validates data, calls service
     ↓
Service: Hashes password, calls repository
     ↓
Repository: Executes SQL INSERT, returns user ID
     ↓
Service: Returns success
     ↓
Handler: Returns HTTP 201 with user data
     ↓
Router: Sends response back to client
```

### **🎓 Learning Points:**

- **Request Flow**: Each layer has a specific responsibility
- **Data Transformation**: Data changes form as it flows through layers
- **Error Propagation**: Errors bubble up from lower layers
- **Separation of Concerns**: Each layer focuses on one thing

---

## 🎯 **11. KEY GO PATTERNS YOU'RE USING**

### **1. Interface Segregation**
```go
type IUserService interface {
    RegisterNewUser(ctx context.Context, u *domain.User) error
    GetUserByID(ctx context.Context, id int) (*domain.User, error)
    GetAllUsers(ctx context.Context, limit int, offset int) ([]domain.User, error)
}
```

### **2. Dependency Injection**
```go
func NewUserService(userRepo repository.IUserRepository) IUserService {
    return &userService{userRepo: userRepo}
}
```

### **3. Error Handling**
```go
if err != nil {
    return err  // Propagate errors up the call chain
}
```

### **4. Context Usage**
```go
func (s *userService) RegisterNewUser(ctx context.Context, u *domain.User) error
```

---

## 🚀 **12. WHAT MAKES YOUR CODE GREAT**

1. **Clean Architecture**: Clear separation between layers
2. **Interface Design**: Easy to test and mock
3. **Error Handling**: Proper error propagation
4. **Security**: Password hashing, input validation
5. **Standards**: RESTful API, proper HTTP status codes
6. **Maintainability**: Each piece has one responsibility

---

## 🎓 **13. NEXT STEPS TO LEARN**

1. **Testing**: Add unit tests for each layer
2. **Middleware**: Add logging, authentication, CORS
3. **Database Migrations**: Version control your database schema
4. **Configuration**: Add more config options (timeouts, connection pools)
5. **Monitoring**: Add metrics and health checks
6. **Documentation**: Generate API documentation

---

## 🎉 **CONCLUSION**

Your code demonstrates **excellent Go practices** and **clean architecture principles**! You've built a:

- ✅ **Well-structured** microservice
- ✅ **Secure** authentication system  
- ✅ **Maintainable** codebase
- ✅ **Scalable** architecture
- ✅ **Professional-grade** application

This is exactly how production Go services should be built. Keep up the great work! 🚀
