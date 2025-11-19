package dto

type CreateCategoryRequest struct {
	Name            string `json:"name" validate:"name"`
	Description     string `json:"description" validate:"description"`
	Slug            string `json:"slug" validate:"slug"`
	ParentID        *int64 `json:"parent_id" validate:"omitempty"`
	IsActive        bool   `json:"is_active"`
	SortOrder       int    `json:"sort_order" validate:"min=0"`
	ImageURL        string `json:"image_url" validate:"omitempty,url"`
	ImagePublicID   string `json:"image_public_id"`
	MetaTitle       string `json:"meta_title" validate:"max=255"`
	MetaDescription string `json:"meta_description"`
}

type UpdateCategoryRequest struct {
	Name            string `json:"name" validate:"omitempty,min=2,max=255"`
	Description     string `json:"description"`
	Slug            string `json:"slug" validate:"omitempty,slug"`
	ParentID        *int64 `json:"parent_id"`
	IsActive        *bool  `json:"is_active"`
	SortOrder       *int   `json:"sort_order" validate:"omitempty,min=0"`
	ImageURL        string `json:"image_url" validate:"omitempty,url"`
	ImagePublicID   string `json:"image_public_id"`
	MetaTitle       string `json:"meta_title" validate:"omitempty,max=255"`
	MetaDescription string `json:"meta_description"`
}
