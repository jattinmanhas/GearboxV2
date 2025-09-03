package dto

type CreateCategoryRequest struct {
	Name            string `json:"name" validate:"name"`
	Description     string `json:"description" validate:"description"`
	Slug            string `json:"slug" validate:"slug"`
	ParentID        *int64 `json:"parent_id" validate:"omitempty"`
	IsActive        bool   `json:"is_active"`
	SortOrder       int    `json:"sort_order" validate:"sort_order"`
	ImageURL        string `json:"image_url" validate:"image_url"`
	MetaTitle       string `json:"meta_title" validate:"meta_title"`
	MetaDescription string `json:"meta_description" validate:"meta_description"`
}

type UpdateCategoryRequest struct {
    Name            string `json:"name" validate:"omitempty,name"`
    Description     string `json:"description" validate:"omitempty,description"`
    Slug            string `json:"slug" validate:"omitempty,slug"`
    ParentID        *int64 `json:"parent_id" validate:"omitempty"`
    IsActive        *bool  `json:"is_active"`
    SortOrder       *int   `json:"sort_order" validate:"omitempty,sort_order"`
    ImageURL        string `json:"image_url" validate:"omitempty,image_url"`
    MetaTitle       string `json:"meta_title" validate:"omitempty,meta_title"`
    MetaDescription string `json:"meta_description" validate:"omitempty,meta_description"`
}