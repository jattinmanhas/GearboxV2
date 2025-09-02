package domain

import (
	"time"
)

// Category represents a product category in the ecommerce system
type Category struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Slug        string    `json:"slug" db:"slug"`
	ParentID    *int64    `json:"parent_id" db:"parent_id"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	MetaTitle   string    `json:"meta_title" db:"meta_title"`
	MetaDesc    string    `json:"meta_description" db:"meta_description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryHierarchy represents the hierarchical structure of categories
type CategoryHierarchy struct {
	Category
	Children []CategoryHierarchy `json:"children,omitempty"`
}

// CategoryFilter represents filters for category queries
type CategoryFilter struct {
	ParentID *int64 `json:"parent_id"`
	IsActive *bool  `json:"is_active"`
	Search   string `json:"search"`
}
