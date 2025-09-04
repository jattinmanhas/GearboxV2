package domain

import (
	"time"
)

// Product represents a product in the ecommerce system
type Product struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	ShortDesc   string `json:"short_description" db:"short_description"`
	SKU         string `json:"sku" db:"sku"`

	Price            float64   `json:"price" db:"price"`
	ComparePrice     float64   `json:"compare_price" db:"compare_price"`
	CostPrice        float64   `json:"cost_price" db:"cost_price"`
	Weight           float64   `json:"weight" db:"weight"`
	Dimensions       string    `json:"dimensions" db:"dimensions"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	IsDigital        bool      `json:"is_digital" db:"is_digital"`
	RequiresShipping bool      `json:"requires_shipping" db:"requires_shipping"`
	Taxable          bool      `json:"taxable" db:"taxable"`
	TrackQuantity    bool      `json:"track_quantity" db:"track_quantity"`
	Quantity         int       `json:"quantity" db:"quantity"`
	MinQuantity      int       `json:"min_quantity" db:"min_quantity"`
	MaxQuantity      int       `json:"max_quantity" db:"max_quantity"`
	MetaTitle        string    `json:"meta_title" db:"meta_title"`
	MetaDesc         string    `json:"meta_description" db:"meta_description"`
	Tags             string    `json:"tags" db:"tags"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ProductVariant represents different variations of a product (size, color, etc.)
type ProductVariant struct {
	ID           int64   `json:"id" db:"id"`
	ProductID    int64   `json:"product_id" db:"product_id"`
	Name         string  `json:"name" db:"name"`
	SKU          string  `json:"sku" db:"sku"`
	Price        float64 `json:"price" db:"price"`
	ComparePrice float64 `json:"compare_price" db:"compare_price"`
	CostPrice    float64 `json:"cost_price" db:"cost_price"`
	Weight       float64 `json:"weight" db:"weight"`
	Quantity     int     `json:"quantity" db:"quantity"`
	IsActive     bool    `json:"is_active" db:"is_active"`
	Position     int     `json:"position" db:"position"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
// ProductAttribute represents product attributes (color, size, material, etc.)
type ProductAttribute struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Value string `json:"value" db:"value"`
	Type  string `json:"type" db:"type"` // text, number, select, multiselect, etc.
}

// ProductAttributeValue represents the value of an attribute for a specific product
type ProductAttributeValue struct {
	ID          int64  `json:"id" db:"id"`
	ProductID   int64  `json:"product_id" db:"product_id"`
	AttributeID int64  `json:"attribute_id" db:"attribute_id"`
	Value       string `json:"value" db:"value"`
}

// ProductImage represents product images
type ProductImage struct {
	ID        int64  `json:"id" db:"id"`
	ProductID int64  `json:"product_id" db:"product_id"`
	URL       string `json:"url" db:"url"`
	Alt       string `json:"alt" db:"alt"`
	Position  int    `json:"position" db:"position"`
	IsPrimary bool   `json:"is_primary" db:"is_primary"`
}

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ID         int64 `json:"id" db:"id"`
	ProductID  int64 `json:"product_id" db:"product_id"`
	CategoryID int64 `json:"category_id" db:"category_id"`
	IsPrimary  bool  `json:"is_primary" db:"is_primary"`
}

// ProductFilter represents filters for product queries
type ProductFilter struct {
	CategoryID *int64   `json:"category_id"`
	IsActive   *bool    `json:"is_active"`
	IsDigital  *bool    `json:"is_digital"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	InStock    *bool    `json:"in_stock"`
	Search     string   `json:"search"`
	Tags       []string `json:"tags"`
	SortBy     string   `json:"sort_by"`    // name, price, created_at, etc.
	SortOrder  string   `json:"sort_order"` // asc, desc
}
