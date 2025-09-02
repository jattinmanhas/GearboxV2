package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*domain.Category, error)
	GetCategory(ctx context.Context, id int64) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error)
	GetCategoryHierarchy(ctx context.Context) ([]*domain.CategoryHierarchy, error)
	GetCategoryChildren(ctx context.Context, parentID int64) ([]*domain.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

type CreateCategoryRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=255"`
	Description     string `json:"description"`
	Slug            string `json:"slug" validate:"required,min=1,max=255"`
	ParentID        *int64 `json:"parent_id"`
	IsActive        bool   `json:"is_active"`
	SortOrder       int    `json:"sort_order"`
	ImageURL        string `json:"image_url"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

type UpdateCategoryRequest struct {
	Name            *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description     *string `json:"description"`
	Slug            *string `json:"slug" validate:"omitempty,min=1,max=255"`
	ParentID        *int64  `json:"parent_id"`
	IsActive        *bool   `json:"is_active"`
	SortOrder       *int    `json:"sort_order"`
	ImageURL        *string `json:"image_url"`
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
}

type ListCategoriesRequest struct {
	ParentID *int64 `json:"parent_id"`
	IsActive *bool  `json:"is_active"`
	Search   string `json:"search"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type ListCategoriesResponse struct {
	Categories []*domain.Category `json:"categories"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

func (s *categoryService) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*domain.Category, error) {
	// Validate parent category exists if provided
	if req.ParentID != nil {
		exists, err := s.categoryRepo.Exists(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate parent category: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("parent category with ID %d does not exist", *req.ParentID)
		}
	}

	// Check if slug already exists
	exists, err := s.categoryRepo.ExistsBySlug(ctx, req.Slug, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("category with slug '%s' already exists", req.Slug)
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = s.generateSlug(req.Name)
	}

	now := time.Now()
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		Slug:        slug,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		ImageURL:    req.ImageURL,
		MetaTitle:   req.MetaTitle,
		MetaDesc:    req.MetaDescription,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id int64) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category with ID %d not found", id)
	}

	return category, nil
}

func (s *categoryService) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category with slug '%s' not found", slug)
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id int64, req *UpdateCategoryRequest) (*domain.Category, error) {
	// Get existing category
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing category: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("category with ID %d not found", id)
	}

	// Validate parent category exists if provided
	if req.ParentID != nil {
		// Prevent setting parent to self
		if *req.ParentID == id {
			return nil, fmt.Errorf("category cannot be its own parent")
		}

		exists, err := s.categoryRepo.Exists(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate parent category: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("parent category with ID %d does not exist", *req.ParentID)
		}
	}

	// Check slug uniqueness if slug is being updated
	if req.Slug != nil && *req.Slug != existing.Slug {
		exists, err := s.categoryRepo.ExistsBySlug(ctx, *req.Slug, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("category with slug '%s' already exists", *req.Slug)
		}
	}

	// Update fields
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Slug != nil {
		existing.Slug = *req.Slug
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		existing.SortOrder = *req.SortOrder
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}
	if req.MetaTitle != nil {
		existing.MetaTitle = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		existing.MetaDesc = *req.MetaDescription
	}

	existing.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return existing, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	// Check if category exists
	exists, err := s.categoryRepo.Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check category existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("category with ID %d not found", id)
	}

	// Check if category has children
	children, err := s.categoryRepo.GetChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check category children: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete category with ID %d: it has child categories", id)
	}

	// TODO: Check if category has products
	// This would require a product repository to check for products in this category

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (s *categoryService) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	filter := &domain.CategoryFilter{
		ParentID: req.ParentID,
		IsActive: req.IsActive,
		Search:   req.Search,
	}

	// Get total count
	total, err := s.categoryRepo.Count(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	// Get categories
	categories, err := s.categoryRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &ListCategoriesResponse{
		Categories: categories,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *categoryService) GetCategoryHierarchy(ctx context.Context) ([]*domain.CategoryHierarchy, error) {
	hierarchy, err := s.categoryRepo.GetHierarchy(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category hierarchy: %w", err)
	}

	return hierarchy, nil
}

func (s *categoryService) GetCategoryChildren(ctx context.Context, parentID int64) ([]*domain.Category, error) {
	children, err := s.categoryRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category children: %w", err)
	}

	return children, nil
}

// generateSlug creates a URL-friendly slug from a name
func (s *categoryService) generateSlug(name string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	// Remove multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
