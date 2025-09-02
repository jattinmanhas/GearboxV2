package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
)

var (
	// CategoryValidator is the global instance for category validation
	CategoryValidator *categoryValidator
)

// CategoryValidator handles validation for category-related requests
type categoryValidator struct{}

// init initializes the global validator instance
func init() {
	CategoryValidator = &categoryValidator{}
}

// ValidateCreateCategoryRequest validates a create category request
func (v *categoryValidator) ValidateCreateCategoryRequest(req *services.CreateCategoryRequest) error {
	if err := v.validateName(req.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	if err := v.validateSlug(req.Slug); err != nil {
		return fmt.Errorf("slug: %w", err)
	}

	if err := v.validateDescription(req.Description); err != nil {
		return fmt.Errorf("description: %w", err)
	}

	if err := v.validateImageURL(req.ImageURL); err != nil {
		return fmt.Errorf("image_url: %w", err)
	}

	if err := v.validateMetaTitle(req.MetaTitle); err != nil {
		return fmt.Errorf("meta_title: %w", err)
	}

	if err := v.validateMetaDescription(req.MetaDescription); err != nil {
		return fmt.Errorf("meta_description: %w", err)
	}

	if err := v.validateSortOrder(req.SortOrder); err != nil {
		return fmt.Errorf("sort_order: %w", err)
	}

	return nil
}

// ValidateUpdateCategoryRequest validates an update category request
func (v *categoryValidator) ValidateUpdateCategoryRequest(req *services.UpdateCategoryRequest) error {
	if req.Name != nil {
		if err := v.validateName(*req.Name); err != nil {
			return fmt.Errorf("name: %w", err)
		}
	}

	if req.Slug != nil {
		if err := v.validateSlug(*req.Slug); err != nil {
			return fmt.Errorf("slug: %w", err)
		}
	}

	if req.Description != nil {
		if err := v.validateDescription(*req.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
	}

	if req.ImageURL != nil {
		if err := v.validateImageURL(*req.ImageURL); err != nil {
			return fmt.Errorf("image_url: %w", err)
		}
	}

	if req.MetaTitle != nil {
		if err := v.validateMetaTitle(*req.MetaTitle); err != nil {
			return fmt.Errorf("meta_title: %w", err)
		}
	}

	if req.MetaDescription != nil {
		if err := v.validateMetaDescription(*req.MetaDescription); err != nil {
			return fmt.Errorf("meta_description: %w", err)
		}
	}

	if req.SortOrder != nil {
		if err := v.validateSortOrder(*req.SortOrder); err != nil {
			return fmt.Errorf("sort_order: %w", err)
		}
	}

	return nil
}

// validateName validates category name
func (v *categoryValidator) validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("is required")
	}

	if len(name) < 1 {
		return fmt.Errorf("must be at least 1 character")
	}

	if len(name) > 255 {
		return fmt.Errorf("must be less than 255 characters")
	}

	// Check for invalid characters
	if strings.ContainsAny(name, "<>\"'&") {
		return fmt.Errorf("contains invalid characters")
	}

	return nil
}

// validateSlug validates category slug
func (v *categoryValidator) validateSlug(slug string) error {
	if strings.TrimSpace(slug) == "" {
		return fmt.Errorf("is required")
	}

	if len(slug) < 1 {
		return fmt.Errorf("must be at least 1 character")
	}

	if len(slug) > 255 {
		return fmt.Errorf("must be less than 255 characters")
	}

	// Slug should only contain lowercase letters, numbers, and hyphens
	slugRegex := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !slugRegex.MatchString(slug) {
		return fmt.Errorf("must contain only lowercase letters, numbers, and hyphens")
	}

	// Slug should not start or end with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return fmt.Errorf("cannot start or end with hyphen")
	}

	// Slug should not contain consecutive hyphens
	if strings.Contains(slug, "--") {
		return fmt.Errorf("cannot contain consecutive hyphens")
	}

	return nil
}

// validateDescription validates category description
func (v *categoryValidator) validateDescription(description string) error {
	if len(description) > 2000 {
		return fmt.Errorf("must be less than 2000 characters")
	}

	return nil
}

// validateImageURL validates image URL
func (v *categoryValidator) validateImageURL(imageURL string) error {
	if imageURL == "" {
		return nil // Optional field
	}

	if len(imageURL) > 500 {
		return fmt.Errorf("must be less than 500 characters")
	}

	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(imageURL) {
		return fmt.Errorf("must be a valid URL")
	}

	return nil
}

// validateMetaTitle validates meta title
func (v *categoryValidator) validateMetaTitle(metaTitle string) error {
	if metaTitle == "" {
		return nil // Optional field
	}

	if len(metaTitle) > 255 {
		return fmt.Errorf("must be less than 255 characters")
	}

	return nil
}

// validateMetaDescription validates meta description
func (v *categoryValidator) validateMetaDescription(metaDescription string) error {
	if metaDescription == "" {
		return nil // Optional field
	}

	if len(metaDescription) > 500 {
		return fmt.Errorf("must be less than 500 characters")
	}

	return nil
}

// validateSortOrder validates sort order
func (v *categoryValidator) validateSortOrder(sortOrder int) error {
	if sortOrder < 0 {
		return fmt.Errorf("must be a non-negative number")
	}

	if sortOrder > 999999 {
		return fmt.Errorf("must be less than 1,000,000")
	}

	return nil
}

// ValidateListCategoriesRequest validates a list categories request
func (v *categoryValidator) ValidateListCategoriesRequest(req *services.ListCategoriesRequest) error {
	if req.Page < 0 {
		return fmt.Errorf("page must be a non-negative number")
	}

	if req.Limit < 0 {
		return fmt.Errorf("limit must be a non-negative number")
	}

	if req.Limit > 100 {
		return fmt.Errorf("limit cannot exceed 100")
	}

	if len(req.Search) > 255 {
		return fmt.Errorf("search term must be less than 255 characters")
	}

	return nil
}
