package dto

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	Name             string  `json:"name" validate:"required,min=1,max=255"`
	Description      string  `json:"description" validate:"required,min=1,max=5000"`
	ShortDesc        string  `json:"short_description" validate:"omitempty,max=500"`
	SKU              string  `json:"sku" validate:"required,sku"`
	Price            float64 `json:"price" validate:"required,price"`
	ComparePrice     float64 `json:"compare_price" validate:"omitempty,min=0"`
	CostPrice        float64 `json:"cost_price" validate:"omitempty,min=0"`
	Weight           float64 `json:"weight" validate:"omitempty,weight"`
	Dimensions       string  `json:"dimensions" validate:"omitempty,dimensions"`
	IsActive         bool    `json:"is_active"`
	IsDigital        bool    `json:"is_digital"`
	RequiresShipping bool    `json:"requires_shipping"`
	Taxable          bool    `json:"taxable"`
	TrackQuantity    bool    `json:"track_quantity"`
	Quantity         int     `json:"quantity" validate:"omitempty,min=0"`
	MinQuantity      int     `json:"min_quantity" validate:"omitempty,min=0"`
	MaxQuantity      int     `json:"max_quantity" validate:"omitempty,min=0"`
	MetaTitle        string  `json:"meta_title" validate:"omitempty,meta_title"`
	MetaDescription  string  `json:"meta_description" validate:"omitempty,meta_description"`
	Tags             string  `json:"tags" validate:"omitempty,tags"`
	CategoryIDs      []int64 `json:"category_ids" validate:"omitempty"`
}

// UpdateProductRequest represents the request to update an existing product
type UpdateProductRequest struct {
	Name             *string  `json:"name" validate:"omitempty,min=1,max=255"`
	Description      *string  `json:"description" validate:"omitempty,min=1,max=5000"`
	ShortDesc        *string  `json:"short_description" validate:"omitempty,max=500"`
	SKU              *string  `json:"sku" validate:"omitempty,sku"`
	Price            *float64 `json:"price" validate:"omitempty,price"`
	ComparePrice     *float64 `json:"compare_price" validate:"omitempty,min=0"`
	CostPrice        *float64 `json:"cost_price" validate:"omitempty,min=0"`
	Weight           *float64 `json:"weight" validate:"omitempty,weight"`
	Dimensions       *string  `json:"dimensions" validate:"omitempty,dimensions"`
	IsActive         *bool    `json:"is_active"`
	IsDigital        *bool    `json:"is_digital"`
	RequiresShipping *bool    `json:"requires_shipping"`
	Taxable          *bool    `json:"taxable"`
	TrackQuantity    *bool    `json:"track_quantity"`
	Quantity         *int     `json:"quantity" validate:"omitempty,min=0"`
	MinQuantity      *int     `json:"min_quantity" validate:"omitempty,min=0"`
	MaxQuantity      *int     `json:"max_quantity" validate:"omitempty,min=0"`
	MetaTitle        *string  `json:"meta_title" validate:"omitempty,meta_title"`
	MetaDescription  *string  `json:"meta_description" validate:"omitempty,meta_description"`
	Tags             *string  `json:"tags" validate:"omitempty,tags"`
	CategoryIDs      []int64  `json:"category_ids" validate:"omitempty"`
}

// ProductResponse represents the response for product data
type ProductResponse struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDesc        string  `json:"short_description"`
	SKU              string  `json:"sku"`
	Price            float64 `json:"price"`
	ComparePrice     float64 `json:"compare_price"`
	CostPrice        float64 `json:"cost_price"`
	Weight           float64 `json:"weight"`
	Dimensions       string  `json:"dimensions"`
	IsActive         bool    `json:"is_active"`
	IsDigital        bool    `json:"is_digital"`
	RequiresShipping bool    `json:"requires_shipping"`
	Taxable          bool    `json:"taxable"`
	TrackQuantity    bool    `json:"track_quantity"`
	Quantity         int     `json:"quantity"`
	MinQuantity      int     `json:"min_quantity"`
	MaxQuantity      int     `json:"max_quantity"`
	MetaTitle        string  `json:"meta_title"`
	MetaDescription  string  `json:"meta_description"`
	Tags             string  `json:"tags"`
	CategoryIDs      []int64 `json:"category_ids"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// ListProductsRequest represents the request to list products with filters
type ListProductsRequest struct {
	CategoryID *int64   `json:"category_id"`
	IsActive   *bool    `json:"is_active"`
	IsDigital  *bool    `json:"is_digital"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	InStock    *bool    `json:"in_stock"`
	Search     string   `json:"search"`
	Tags       []string `json:"tags"`
	SortBy     string   `json:"sort_by"`
	SortOrder  string   `json:"sort_order"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
}

// ListProductsResponse represents the response for listing products
type ListProductsResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// CreateProductVariantRequest represents the request to create a product variant
type CreateProductVariantRequest struct {
	ProductID    int64   `json:"product_id" validate:"required"`
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	SKU          string  `json:"sku" validate:"required,min=1,max=100"`
	Price        float64 `json:"price" validate:"required,price"`
	ComparePrice float64 `json:"compare_price" validate:"omitempty,price"`
	CostPrice    float64 `json:"cost_price" validate:"omitempty,price"`
	Weight       float64 `json:"weight" validate:"omitempty,weight"`
	Quantity     int     `json:"quantity" validate:"omitempty,min=0"`
	IsActive     bool    `json:"is_active"`
	Position     int     `json:"position" validate:"omitempty,min=0"`
}

// UpdateProductVariantRequest represents the request to update a product variant
type UpdateProductVariantRequest struct {
	Name         *string  `json:"name" validate:"omitempty,min=1,max=255"`
	SKU          *string  `json:"sku" validate:"omitempty,min=1,max=100"`
	Price        *float64 `json:"price" validate:"omitempty,price"`
	ComparePrice *float64 `json:"compare_price" validate:"omitempty,price"`
	CostPrice    *float64 `json:"cost_price" validate:"omitempty,price"`
	Weight       *float64 `json:"weight" validate:"omitempty,weight"`
	Quantity     *int     `json:"quantity" validate:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active"`
	Position     *int     `json:"position" validate:"omitempty,min=0"`
}

// ProductVariantResponse represents the response for product variant data
type ProductVariantResponse struct {
	ID           int64   `json:"id"`
	ProductID    int64   `json:"product_id"`
	Name         string  `json:"name"`
	SKU          string  `json:"sku"`
	Price        float64 `json:"price"`
	ComparePrice float64 `json:"compare_price"`
	CostPrice    float64 `json:"cost_price"`
	Weight       float64 `json:"weight"`
	Quantity     int     `json:"quantity"`
	IsActive     bool    `json:"is_active"`
	Position     int     `json:"position"`
}