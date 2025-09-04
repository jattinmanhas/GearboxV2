package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error)
	GetCategory(ctx context.Context, id int64) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error)
	GetCategoryHierarchy(ctx context.Context) ([]*domain.CategoryHierarchy, error)
	GetCategoryChildren(ctx context.Context, parentID int64) ([]*domain.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, productRepo repository.ProductRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
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

func (s *categoryService) CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	// Validate parent category exists if provided
	if cat.ParentID != nil {
		exists, err := s.categoryRepo.Exists(ctx, *cat.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate parent category: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("parent category with ID %d does not exist", *cat.ParentID)
		}
	}

	// Check if slug already exists
	exists, err := s.categoryRepo.ExistsBySlug(ctx, cat.Slug, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("category with slug '%s' already exists", cat.Slug)
	}

	// Generate slug if not provided
	if cat.Slug == "" {
		cat.Slug = s.generateSlug(cat.Name)
	}

	if err := s.categoryRepo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return cat, nil
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

func (s *categoryService) UpdateCategory(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) (*domain.Category, error) {
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

		// check if the gived parent id exists...
		exists, err := s.categoryRepo.Exists(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate parent category: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("parent category with ID %d does not exist", *req.ParentID)
		}
	}

	// Validate slug uniqueness if provided
	if req.Slug != "" && req.Slug != existing.Slug {
		exists, err := s.categoryRepo.ExistsBySlug(ctx, req.Slug, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("category with slug '%s' already exists", req.Slug)
		}
	}

	updatedCategory := *existing

	// Apply updates only for provided fields
	if req.Name != "" {
		updatedCategory.Name = req.Name
	}
	if req.Description != "" {
		updatedCategory.Description = req.Description
	}
	if req.Slug != "" {
		updatedCategory.Slug = req.Slug
	}
	if req.ParentID != nil {
		updatedCategory.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		updatedCategory.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		updatedCategory.SortOrder = *req.SortOrder
	}
	if req.ImageURL != "" {
		updatedCategory.ImageURL = req.ImageURL
	}
	if req.MetaTitle != "" {
		updatedCategory.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		updatedCategory.MetaDesc = req.MetaDescription
	}

	updatedCategory.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, &updatedCategory); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &updatedCategory, nil
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

	// Check if category has products
	hasProducts, err := s.productRepo.CheckCategoryHasProducts(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check category has products: %w", err)
	}
	if hasProducts {
		return fmt.Errorf("cannot delete category with ID %d: it has products", id)
	}

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
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	filter := &domain.CategoryFilter{
		ParentID: req.ParentID,
		IsActive: req.IsActive,
		Search:   req.Search,
		Page:     req.Page,
		Limit:    req.Limit,
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