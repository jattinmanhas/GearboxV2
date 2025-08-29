package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type IRoleRepository interface {
	// Role management - only read operations since roles are static
	GetRoleByID(ctx context.Context, id uint) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]*domain.Role, error)

	// User-role management - simplified since role_id is in users table
	UpdateUserRole(ctx context.Context, userID, roleID uint) error
	GetUserRole(ctx context.Context, userID uint) (*domain.Role, error)
	RemoveUserRole(ctx context.Context, userID uint) error

	// Initialize default roles
	InitializeDefaultRoles(ctx context.Context) error
}

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) IRoleRepository {
	return &roleRepository{db: db}
}

// CreateRole creates a new role
func (r *roleRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (name, description, is_active, created_at, updated_at, is_deleted)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		role.Name,
		role.Description,
		role.IsActive,
		role.CreatedAt,
		role.UpdatedAt,
		role.IsDeleted,
	).Scan(&role.ID)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetRoleByID retrieves a role by ID
func (r *roleRepository) GetRoleByID(ctx context.Context, id uint) (*domain.Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at, is_deleted
		FROM roles
		WHERE id = $1 AND is_deleted = false`

	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// GetRoleByName retrieves a role by name
func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at, is_deleted
		FROM roles
		WHERE name = $1 AND is_deleted = false`

	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// GetAllRoles retrieves all active roles
func (r *roleRepository) GetAllRoles(ctx context.Context) ([]*domain.Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at, is_deleted
		FROM roles
		WHERE is_deleted = false
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.IsActive,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateRole updates an existing role
func (r *roleRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, is_active = $3, updated_at = $4
		WHERE id = $5 AND is_deleted = false`

	role.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query,
		role.Name,
		role.Description,
		role.IsActive,
		role.UpdatedAt,
		role.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found or already deleted")
	}

	return nil
}

// DeleteRole soft deletes a role
func (r *roleRepository) DeleteRole(ctx context.Context, id uint) error {
	query := `
		UPDATE roles
		SET is_deleted = true, updated_at = $1
		WHERE id = $2 AND is_deleted = false`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found or already deleted")
	}

	return nil
}

// UpdateUserRole updates the role assigned to a user (directly in users table)
func (r *roleRepository) UpdateUserRole(ctx context.Context, userID, roleID uint) error {
	query := `
		UPDATE users
		SET role_id = $1, updated_at = $2
		WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, roleID, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

// GetUserRole retrieves the role assigned to a user
func (r *roleRepository) GetUserRole(ctx context.Context, userID uint) (*domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_active, r.created_at, r.updated_at, r.is_deleted
		FROM roles r
		JOIN users u ON r.id = u.role_id
		WHERE u.id = $1 AND r.is_deleted = false`

	role := &domain.Role{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user role not found")
		}
		return nil, fmt.Errorf("failed to get user role: %w", err)
	}

	return role, nil
}

// RemoveUserRole removes the role from a user (sets to default user role)
func (r *roleRepository) RemoveUserRole(ctx context.Context, userID uint) error {
	return r.UpdateUserRole(ctx, userID, domain.RoleIDUser)
}

// InitializeDefaultRoles creates the three default roles if they don't exist
func (r *roleRepository) InitializeDefaultRoles(ctx context.Context) error {
	defaultRoles := []*domain.Role{
		{
			ID:          domain.RoleIDUser,
			Name:        domain.RoleUser,
			Description: "Basic authenticated user with limited access",
			IsActive:    true,
		},
		{
			ID:          domain.RoleIDEditor,
			Name:        domain.RoleEditor,
			Description: "Content editor with create/edit/moderate permissions",
			IsActive:    true,
		},
		{
			ID:          domain.RoleIDAdmin,
			Name:        domain.RoleAdmin,
			Description: "Full system administrator with complete access",
			IsActive:    true,
		},
	}

	for _, role := range defaultRoles {
		// Try to get existing role
		existingRole, err := r.GetRoleByID(ctx, role.ID)
		if err != nil {
			// Role doesn't exist, create it
			if err := r.CreateRole(ctx, role); err != nil {
				return fmt.Errorf("failed to create default role %s: %w", role.Name, err)
			}
		} else {
			// Role exists, update it if needed
			if existingRole.Name != role.Name || existingRole.Description != role.Description {
				role.ID = existingRole.ID
				role.CreatedAt = existingRole.CreatedAt
				if err := r.UpdateRole(ctx, role); err != nil {
					return fmt.Errorf("failed to update default role %s: %w", role.Name, err)
				}
			}
		}
	}

	return nil
}
