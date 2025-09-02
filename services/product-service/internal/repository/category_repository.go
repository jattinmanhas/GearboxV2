package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filter *domain.CategoryFilter) ([]*domain.Category, error)
	GetHierarchy(ctx context.Context) ([]*domain.CategoryHierarchy, error)
	GetChildren(ctx context.Context, parentID int64) ([]*domain.Category, error)
	Exists(ctx context.Context, id int64) (bool, error)
	ExistsBySlug(ctx context.Context, slug string, excludeID *int64) (bool, error)
	Count(ctx context.Context, filter *domain.CategoryFilter) (int64, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	query := `
		INSERT INTO categories (
			name, description, slug, parent_id, is_active, sort_order,
			image_url, meta_title, meta_description, created_at, updated_at
		) VALUES (
			:name, :description, :slug, :parent_id, :is_active, :sort_order,
			:image_url, :meta_title, :meta_description, :created_at, :updated_at
		) RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, category)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&category.ID); err != nil {
			return fmt.Errorf("failed to get created category ID: %w", err)
		}
	}

	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	query := `
		SELECT id, name, description, slug, parent_id, is_active, sort_order,
			   image_url, meta_title, meta_description, created_at, updated_at
		FROM categories
		WHERE id = $1`

	var category domain.Category
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by ID: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := `
		SELECT id, name, description, slug, parent_id, is_active, sort_order,
			   image_url, meta_title, meta_description, created_at, updated_at
		FROM categories
		WHERE slug = $1`

	var category domain.Category
	err := r.db.GetContext(ctx, &category, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	query := `
		UPDATE categories SET
			name = :name,
			description = :description,
			slug = :slug,
			parent_id = :parent_id,
			is_active = :is_active,
			sort_order = :sort_order,
			image_url = :image_url,
			meta_title = :meta_title,
			meta_description = :meta_description,
			updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, category)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", category.ID)
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category with ID %d not found", id)
	}

	return nil
}

func (r *categoryRepository) List(ctx context.Context, filter *domain.CategoryFilter) ([]*domain.Category, error) {
	query := `
		SELECT id, name, description, slug, parent_id, is_active, sort_order,
			   image_url, meta_title, meta_description, created_at, updated_at
		FROM categories
		WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.ParentID != nil {
			query += fmt.Sprintf(" AND parent_id = $%d", argIndex)
			args = append(args, *filter.ParentID)
			argIndex++
		}

		if filter.IsActive != nil {
			query += fmt.Sprintf(" AND is_active = $%d", argIndex)
			args = append(args, *filter.IsActive)
			argIndex++
		}

		if filter.Search != "" {
			query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm)
			argIndex++
		}
	}

	query += " ORDER BY sort_order ASC, name ASC"

	var categories []*domain.Category
	err := r.db.SelectContext(ctx, &categories, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) GetHierarchy(ctx context.Context) ([]*domain.CategoryHierarchy, error) {
	// First, get all categories
	categories, err := r.List(ctx, &domain.CategoryFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get categories for hierarchy: %w", err)
	}

	// Build a map for quick lookup
	categoryMap := make(map[int64]*domain.CategoryHierarchy)
	var rootCategories []*domain.CategoryHierarchy

	// Convert to hierarchy structs
	for _, cat := range categories {
		hierarchy := &domain.CategoryHierarchy{
			Category: *cat,
			Children: []domain.CategoryHierarchy{},
		}
		categoryMap[cat.ID] = hierarchy

		if cat.ParentID == nil {
			rootCategories = append(rootCategories, hierarchy)
		}
	}

	// Build the hierarchy
	for _, cat := range categories {
		if cat.ParentID != nil {
			if parent, exists := categoryMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, *categoryMap[cat.ID])
			}
		}
	}

	return rootCategories, nil
}

func (r *categoryRepository) GetChildren(ctx context.Context, parentID int64) ([]*domain.Category, error) {
	filter := &domain.CategoryFilter{
		ParentID: &parentID,
	}
	return r.List(ctx, filter)
}

func (r *categoryRepository) Exists(ctx context.Context, id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check category existence: %w", err)
	}

	return exists, nil
}

func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1`
	args := []interface{}{slug}

	if excludeID != nil {
		query += ` AND id != $2`
		args = append(args, *excludeID)
	}
	query += `)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check category slug existence: %w", err)
	}

	return exists, nil
}

func (r *categoryRepository) Count(ctx context.Context, filter *domain.CategoryFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM categories WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filter != nil {
		if filter.ParentID != nil {
			query += fmt.Sprintf(" AND parent_id = $%d", argIndex)
			args = append(args, *filter.ParentID)
			argIndex++
		}

		if filter.IsActive != nil {
			query += fmt.Sprintf(" AND is_active = $%d", argIndex)
			args = append(args, *filter.IsActive)
			argIndex++
		}

		if filter.Search != "" {
			query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
			searchTerm := "%" + filter.Search + "%"
			args = append(args, searchTerm)
			argIndex++
		}
	}

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return count, nil
}
