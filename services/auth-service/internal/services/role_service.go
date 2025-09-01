package services

import (
	"context"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
)

type IRoleService interface {
	// Role management - only read operations since roles are static
	GetRoleByID(ctx context.Context, id uint) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)

	// User-role management
	AssignRoleToUser(ctx context.Context, userID, roleID, grantedBy uint) error
	GetUserRole(ctx context.Context, userID uint) (*domain.Role, error)
	UpdateUserRole(ctx context.Context, userID, roleID, grantedBy uint) error
	RemoveUserRole(ctx context.Context, userID uint, removedBy uint) error

	// Authorization helpers
	CheckUserPermission(ctx context.Context, userID uint, requiredRole string) error
	GetUserRoleName(ctx context.Context, userID uint) (string, error)
	IsUserAdmin(ctx context.Context, userID uint) (bool, error)
	IsUserEditor(ctx context.Context, userID uint) (bool, error)

	// Initialize default roles
	InitializeDefaultRoles(ctx context.Context) error
}

type roleService struct {
	roleRepo repository.IRoleRepository
	userRepo repository.IUserRepository
}

func NewRoleService(roleRepo repository.IRoleRepository, userRepo repository.IUserRepository) IRoleService {
	return &roleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// GetRoleByID retrieves a role by ID
func (r *roleService) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, id)
}

// GetRoleByName retrieves a role by name
func (r *roleService) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	return r.roleRepo.GetRoleByName(ctx, name)
}

// GetAllRoles retrieves all active roles
func (r *roleService) GetAllRoles(ctx context.Context) ([]*domain.Role, error) {
	return r.roleRepo.GetAllRoles(ctx)
}

// AssignRoleToUser assigns a role to a user
func (r *roleService) AssignRoleToUser(ctx context.Context, userID, roleID, grantedBy uint) error {
	// Validate user exists
	_, err := r.userRepo.GetUserByID(ctx, int(userID))
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Validate role exists
	role, err := r.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Check if granter has permission to assign this role
	granterRole, err := r.roleRepo.GetUserRole(ctx, grantedBy)
	if err != nil {
		return fmt.Errorf("granter role not found: %w", err)
	}

	// Only admins can assign admin roles, editors can assign user roles
	if role.Name == domain.RoleAdmin && granterRole.Name != domain.RoleAdmin {
		return fmt.Errorf("only administrators can assign admin roles")
	}

	if role.Name == domain.RoleEditor && !domain.CanAccess(granterRole.Name, domain.RoleEditor) {
		return fmt.Errorf("insufficient permissions to assign editor role")
	}

	return r.roleRepo.UpdateUserRole(ctx, userID, roleID)
}

// GetUserRole retrieves the role assigned to a user
func (r *roleService) GetUserRole(ctx context.Context, userID uint) (*domain.Role, error) {
	return r.roleRepo.GetUserRole(ctx, userID)
}

// UpdateUserRole updates the role assigned to a user
func (r *roleService) UpdateUserRole(ctx context.Context, userID, roleID, grantedBy uint) error {
	return r.AssignRoleToUser(ctx, userID, roleID, grantedBy)
}

// RemoveUserRole removes the role from a user
func (r *roleService) RemoveUserRole(ctx context.Context, userID uint, removedBy uint) error {
	// Check if remover has permission
	removerRole, err := r.roleRepo.GetUserRole(ctx, removedBy)
	if err != nil {
		return fmt.Errorf("remover role not found: %w", err)
	}

	// Get user's current role
	userRole, err := r.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("user role not found: %w", err)
	}

	// Only admins can remove admin roles, editors can remove user roles
	if userRole.Name == domain.RoleAdmin && removerRole.Name != domain.RoleAdmin {
		return fmt.Errorf("only administrators can remove admin roles")
	}

	if userRole.Name == domain.RoleEditor && !domain.CanAccess(removerRole.Name, domain.RoleEditor) {
		return fmt.Errorf("insufficient permissions to remove editor role")
	}

	return r.roleRepo.RemoveUserRole(ctx, userID)
}

// CheckUserPermission checks if a user has the required role level
func (r *roleService) CheckUserPermission(ctx context.Context, userID uint, requiredRole string) error {
	userRole, err := r.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		return fmt.Errorf("user role not found: %w", err)
	}

	if !domain.CanAccess(userRole.Name, requiredRole) {
		return fmt.Errorf("insufficient permissions: requires %s role, user has %s", requiredRole, userRole.Name)
	}

	return nil
}

// GetUserRoleName gets the role name for a user
func (r *roleService) GetUserRoleName(ctx context.Context, userID uint) (string, error) {
	role, err := r.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		return "", err
	}
	return role.Name, nil
}

// IsUserAdmin checks if a user is an admin
func (r *roleService) IsUserAdmin(ctx context.Context, userID uint) (bool, error) {
	role, err := r.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		return false, err
	}
	return role.Name == domain.RoleAdmin, nil
}

// IsUserEditor checks if a user is an editor or admin
func (r *roleService) IsUserEditor(ctx context.Context, userID uint) (bool, error) {
	role, err := r.roleRepo.GetUserRole(ctx, userID)
	if err != nil {
		return false, err
	}
	return domain.CanAccess(role.Name, domain.RoleEditor), nil
}

// InitializeDefaultRoles creates the three default roles
func (r *roleService) InitializeDefaultRoles(ctx context.Context) error {
	return r.roleRepo.InitializeDefaultRoles(ctx)
}
