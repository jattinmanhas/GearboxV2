# 🔐 JWT Authentication System

This document describes the comprehensive JWT-based authentication system implemented in the Auth Service.

## 🚀 **Features**

- **JWT Access Tokens**: Short-lived (15 minutes) for API access
- **JWT Refresh Tokens**: Long-lived (7 days) for token renewal
- **HTTP-Only Cookies**: Secure token storage
- **Token Revocation**: Ability to revoke individual or all user tokens
- **Security Features**: bcrypt password hashing, secure token generation
- **Middleware**: Authentication and authorization middleware
- **CORS Support**: Cross-origin request handling

## 🏗️ **Architecture**

### **Layers**
1. **Handlers**: HTTP request/response handling
2. **Services**: Business logic and token management
3. **Repositories**: Data persistence
4. **Middleware**: Authentication and CORS handling

### **Token Flow**
```
User Login → Generate Access + Refresh Tokens → Store in HTTP-Only Cookies
     ↓
API Request → Validate Access Token → Allow/Deny Access
     ↓
Token Expired → Use Refresh Token → Generate New Access + Refresh Tokens
```

## 📋 **API Endpoints**

### **Public Routes (No Authentication Required)**
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh tokens
- `POST /api/v1/auth/logout` - Logout (revoke refresh token)

### **Protected Routes (Authentication Required)**
- `GET /api/v1/auth/user/{id}` - Get user by ID
- `PUT /api/v1/auth/user/{id}` - Update user
- `DELETE /api/v1/auth/user/{id}` - Delete user
- `POST /api/v1/auth/user/{id}/change-password` - Change password
- `GET /api/v1/auth/users` - Get all users (with pagination)
- `POST /api/v1/auth/logout-all` - Logout from all devices

## 🔧 **Configuration**

### **Environment Variables**
```bash
# Database
DATABASE_URL=postgres://username:password@localhost:5432/auth_service?sslmode=disable

# Server
PORT=8081
ENVIRONMENT=development

# JWT Secrets (REQUIRED - Set these in production!)
JWT_SECRET=your-super-secret-jwt-key-here-make-it-long-and-random
JWT_REFRESH_SECRET=your-super-secret-refresh-token-key-here-make-it-long-and-random
```

### **Security Requirements**
- **JWT_SECRET**: Minimum 32 characters, cryptographically secure
- **JWT_REFRESH_SECRET**: Different from JWT_SECRET, minimum 32 characters
- **Environment**: Set to "production" in production environments

## 🗄️ **Database Schema**

### **Users Table**
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    middle_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar TEXT,
    gender VARCHAR(20),
    date_of_birth TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### **Refresh Tokens Table**
```sql
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_revoked BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## 🔐 **Authentication Flow**

### **1. User Registration**
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "john_doe",
    "password": "SecurePass123",
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe"
}
```

**Response**: User created with hashed password

### **2. User Login**
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "john_doe",
    "password": "SecurePass123"
}
```

**Response**: 
- Access token stored in HTTP-only cookie (`access_token`)
- Refresh token stored in HTTP-only cookie (`refresh_token`)
- User information returned in response body

### **3. API Access**
```bash
GET /api/v1/auth/user/1
# Access token automatically sent via HTTP-only cookie
```

**Response**: User data if token is valid

### **4. Token Refresh**
```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response**: New access and refresh tokens stored in cookies

### **5. Logout**
```bash
POST /api/v1/auth/logout
# Refresh token automatically sent via HTTP-only cookie
```

**Response**: Tokens revoked and cookies cleared

## 🛡️ **Security Features**

### **Password Security**
- bcrypt hashing with default cost (10)
- Unique salt per password
- Secure password validation

### **Token Security**
- JWT tokens signed with HMAC-SHA256
- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- Tokens stored in HTTP-only, secure cookies
- SameSite=Strict cookie policy

### **Database Security**
- Refresh tokens stored with user agent and IP tracking
- Token revocation capability
- Automatic cleanup of expired tokens
- Cascade deletion of user tokens

## 🔧 **Middleware Usage**

### **Authentication Middleware**
```go
// Protect routes
router.Use(middleware.AuthMiddleware(authService))

// Optional authentication
router.Use(middleware.OptionalAuthMiddleware(authService))
```

### **CORS Middleware**
```go
// Enable CORS for all origins
router.Use(middleware.CORSMiddleware([]string{"*"}))

// Enable CORS for specific origins
router.Use(middleware.CORSMiddleware([]string{"https://example.com", "https://app.example.com"}))
```

## 🧪 **Testing**

### **Run All Tests**
```bash
go test ./... -v
```

### **Run Specific Test Suites**
```bash
# Test authentication services
go test ./internal/services -v -run TestUserService

# Test configuration
go test ./internal/config -v

# Test HTTP responses
go test ./internal/httpx -v
```

## 🚀 **Running the Service**

### **1. Set Environment Variables**
```bash
export DATABASE_URL="postgres://username:password@localhost:5432/auth_service?sslmode=disable"
export JWT_SECRET="your-super-secret-jwt-key-here-make-it-long-and-random"
export JWT_REFRESH_SECRET="your-super-secret-refresh-token-key-here-make-it-long-and-random"
export ENVIRONMENT="development"
export PORT="8081"
```

### **2. Run Database Migrations**
```bash
# Apply migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback if needed
migrate -path migrations -database "$DATABASE_URL" down
```

### **3. Start the Service**
```bash
go run cmd/api/main.go
```

## 🔍 **Troubleshooting**

### **Common Issues**

1. **"JWT_SECRET environment variable is not set"**
   - Set the JWT_SECRET environment variable
   - Ensure it's at least 32 characters long

2. **"invalid access token"**
   - Check if token is expired (15 minutes)
   - Verify token is being sent via cookie or Authorization header
   - Ensure JWT_SECRET matches the one used to generate the token

3. **"refresh token not found or expired"**
   - Refresh token may have expired (7 days)
   - Token may have been revoked
   - User may need to login again

4. **Database Connection Issues**
   - Verify DATABASE_URL is correct
   - Ensure PostgreSQL is running
   - Check database permissions

### **Debug Mode**
Set `ENVIRONMENT=development` to enable detailed error messages.

## 🔒 **Production Considerations**

### **Security Checklist**
- [ ] Use strong, randomly generated JWT secrets
- [ ] Set ENVIRONMENT=production
- [ ] Use HTTPS in production
- [ ] Regularly rotate JWT secrets
- [ ] Monitor token usage and revoke suspicious tokens
- [ ] Implement rate limiting
- [ ] Use secure database connections
- [ ] Regular security audits

### **Performance Optimization**
- [ ] Database connection pooling
- [ ] Redis for token caching (optional)
- [ ] CDN for static assets
- [ ] Load balancing for high availability

## 📚 **Additional Resources**

- [JWT.io](https://jwt.io/) - JWT token debugging and validation
- [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) - Password hashing algorithm
- [HTTP-Only Cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies#Secure_and_HttpOnly_cookies) - Cookie security
- [CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) - Cross-origin resource sharing

## 🤝 **Contributing**

When contributing to the authentication system:

1. **Security First**: All changes must maintain or improve security
2. **Test Coverage**: New features must include comprehensive tests
3. **Documentation**: Update this README for any API changes
4. **Code Review**: Security-related changes require thorough review

---

**⚠️ Security Note**: This authentication system is designed for production use but should be thoroughly tested in your specific environment before deployment.
