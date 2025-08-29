package domain

import "time"

// Role represents a user role in the system
type Role struct {
	ID          uint      `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
}

// UserWithRole represents a user with their role information
type UserWithRole struct {
	User *User `json:"user"`
	Role *Role `json:"role"`
}

// Role constants - these are the three main roles
const (
	RoleUser   = "user"   // Basic authenticated user
	RoleEditor = "editor" // Can create/edit content, moderate
	RoleAdmin  = "admin"  // Full system access
)

// Role IDs - these will be set when we create the initial roles
const (
	RoleIDUser   = 1
	RoleIDEditor = 2
	RoleIDAdmin  = 3
)

// GetDefaultRole returns the default role for new users
func GetDefaultRole() string {
	return RoleUser
}

// IsValidRole checks if a role name is valid
func IsValidRole(roleName string) bool {
	switch roleName {
	case RoleUser, RoleEditor, RoleAdmin:
		return true
	default:
		return false
	}
}

// GetRoleHierarchy returns the hierarchy level of a role (higher = more permissions)
func GetRoleHierarchy(roleName string) int {
	switch roleName {
	case RoleUser:
		return 1
	case RoleEditor:
		return 2
	case RoleAdmin:
		return 3
	default:
		return 0
	}
}

// CanAccess checks if a role can access another role's level
func CanAccess(userRole, requiredRole string) bool {
	return GetRoleHierarchy(userRole) >= GetRoleHierarchy(requiredRole)
}
