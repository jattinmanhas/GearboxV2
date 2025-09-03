package services

import (
	"context"
	"fmt"

	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*domain.Product, error)
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context, req *dto.ListProductsRequest) (*dto.ListProductsResponse, error)
	GetProductsByCategory(ctx context.Context, categoryID int64, page, limit int) (*dto.ListProductsResponse, error)
	SearchProducts(ctx context.Context, query string, page, limit int) (*dto.ListProductsResponse, error)
	UpdateProductQuantity(ctx context.Context, id int64, quantity int) error
	GetProductsByTags(ctx context.Context, tags []string, page, limit int) (*dto.ListProductsResponse, error)

	// Product Variants
	CreateProductVariant(ctx context.Context, req *dto.CreateProductVariantRequest) (*domain.ProductVariant, error)
	GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error)
	GetProductVariantsByProductID(ctx context.Context, productID int64) ([]*domain.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, id int64, req *dto.UpdateProductVariantRequest) (*domain.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, id int64) error

	// Product Categories
	AddProductToCategory(ctx context.Context, productID, categoryID int64, isPrimary bool) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error
	GetProductCategories(ctx context.Context, productID int64) ([]*domain.Category, error)
	UpdateProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*domain.Product, error) {
	// Check if SKU already exists
	existingProduct, err := s.productRepo.GetProductBySKU(ctx, req.SKU)
	if err == nil && existingProduct != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
	}

	// Create product domain object
	product := &domain.Product{
		Name:             req.Name,
		Description:      req.Description,
		ShortDesc:        req.ShortDesc,
		SKU:              req.SKU,
		Price:            req.Price,
		ComparePrice:     req.ComparePrice,
		CostPrice:        req.CostPrice,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		IsActive:         req.IsActive,
		IsDigital:        req.IsDigital,
		RequiresShipping: req.RequiresShipping,
		Taxable:          req.Taxable,
		TrackQuantity:    req.TrackQuantity,
		Quantity:         req.Quantity,
		MinQuantity:      req.MinQuantity,
		MaxQuantity:      req.MaxQuantity,
		MetaTitle:        req.MetaTitle,
		MetaDesc:         req.MetaDescription,
		Tags:             req.Tags,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Create product in repository
	err = s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Add product to categories if provided
	if len(req.CategoryIDs) > 0 {
		err = s.productRepo.UpdateProductCategories(ctx, product.ID, req.CategoryIDs)
		if err != nil {
			// Log error but don't fail the product creation
			// In a real application, you might want to handle this differently
			fmt.Printf("Warning: failed to add product to categories: %v\n", err)
		}
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (s *productService) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// GetProductBySKU retrieves a product by SKU
func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	product, err := s.productRepo.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*domain.Product, error) {
	// Get existing product
	existingProduct, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing product: %w", err)
	}

	// Check SKU uniqueness if SKU is being updated
	if req.SKU != nil && *req.SKU != existingProduct.SKU {
		skuProduct, err := s.productRepo.GetProductBySKU(ctx, *req.SKU)
		if err == nil && skuProduct != nil && skuProduct.ID != id {
			return nil, fmt.Errorf("product with SKU %s already exists", *req.SKU)
		}
	}

	// Update fields that are provided
	updateProduct := *existingProduct

	if req.Name != nil {
		updateProduct.Name = *req.Name
	}
	if req.Description != nil {
		updateProduct.Description = *req.Description
	}
	if req.ShortDesc != nil {
		updateProduct.ShortDesc = *req.ShortDesc
	}
	if req.SKU != nil {
		updateProduct.SKU = *req.SKU
	}
	if req.Price != nil {
		updateProduct.Price = *req.Price
	}
	if req.ComparePrice != nil {
		updateProduct.ComparePrice = *req.ComparePrice
	}
	if req.CostPrice != nil {
		updateProduct.CostPrice = *req.CostPrice
	}
	if req.Weight != nil {
		updateProduct.Weight = *req.Weight
	}
	if req.Dimensions != nil {
		updateProduct.Dimensions = *req.Dimensions
	}
	if req.IsActive != nil {
		updateProduct.IsActive = *req.IsActive
	}
	if req.IsDigital != nil {
		updateProduct.IsDigital = *req.IsDigital
	}
	if req.RequiresShipping != nil {
		updateProduct.RequiresShipping = *req.RequiresShipping
	}
	if req.Taxable != nil {
		updateProduct.Taxable = *req.Taxable
	}
	if req.TrackQuantity != nil {
		updateProduct.TrackQuantity = *req.TrackQuantity
	}
	if req.Quantity != nil {
		updateProduct.Quantity = *req.Quantity
	}
	if req.MinQuantity != nil {
		updateProduct.MinQuantity = *req.MinQuantity
	}
	if req.MaxQuantity != nil {
		updateProduct.MaxQuantity = *req.MaxQuantity
	}
	if req.MetaTitle != nil {
		updateProduct.MetaTitle = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		updateProduct.MetaDesc = *req.MetaDescription
	}
	if req.Tags != nil {
		updateProduct.Tags = *req.Tags
	}

	updateProduct.UpdatedAt = time.Now()

	// Update product in repository
	err = s.productRepo.UpdateProduct(ctx, id, &updateProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Update categories if provided
	if req.CategoryIDs != nil {
		err = s.productRepo.UpdateProductCategories(ctx, id, req.CategoryIDs)
		if err != nil {
			// Log error but don't fail the product update
			fmt.Printf("Warning: failed to update product categories: %v\n", err)
		}
	}

	return &updateProduct, nil
}

// DeleteProduct deletes a product
func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Delete product
	err = s.productRepo.DeleteProduct(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// ListProducts retrieves products with filters
func (s *productService) ListProducts(ctx context.Context, req *dto.ListProductsRequest) (*dto.ListProductsResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	offset := (req.Page - 1) * req.Limit

	// Build filter
	filter := &domain.ProductFilter{
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
		IsDigital:  req.IsDigital,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		InStock:    req.InStock,
		Search:     req.Search,
		Tags:       req.Tags,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
	}

	// Get products from repository
	products, total, err := s.productRepo.ListProducts(ctx, filter, offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Convert to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID:               product.ID,
			Name:             product.Name,
			Description:      product.Description,
			ShortDesc:        product.ShortDesc,
			SKU:              product.SKU,
			Price:            product.Price,
			ComparePrice:     product.ComparePrice,
			CostPrice:        product.CostPrice,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsDigital:        product.IsDigital,
			RequiresShipping: product.RequiresShipping,
			Taxable:          product.Taxable,
			TrackQuantity:    product.TrackQuantity,
			Quantity:         product.Quantity,
			MinQuantity:      product.MinQuantity,
			MaxQuantity:      product.MaxQuantity,
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CreatedAt:        product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        product.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetProductsByCategory retrieves products by category
func (s *productService) GetProductsByCategory(ctx context.Context, categoryID int64, page, limit int) (*dto.ListProductsResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Get products from repository
	products, total, err := s.productRepo.GetProductsByCategory(ctx, categoryID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	// Convert to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID:               product.ID,
			Name:             product.Name,
			Description:      product.Description,
			ShortDesc:        product.ShortDesc,
			SKU:              product.SKU,
			Price:            product.Price,
			ComparePrice:     product.ComparePrice,
			CostPrice:        product.CostPrice,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsDigital:        product.IsDigital,
			RequiresShipping: product.RequiresShipping,
			Taxable:          product.Taxable,
			TrackQuantity:    product.TrackQuantity,
			Quantity:         product.Quantity,
			MinQuantity:      product.MinQuantity,
			MaxQuantity:      product.MaxQuantity,
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CreatedAt:        product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        product.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// SearchProducts searches products by query
func (s *productService) SearchProducts(ctx context.Context, query string, page, limit int) (*dto.ListProductsResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Get products from repository
	products, total, err := s.productRepo.SearchProducts(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	// Convert to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID:               product.ID,
			Name:             product.Name,
			Description:      product.Description,
			ShortDesc:        product.ShortDesc,
			SKU:              product.SKU,
			Price:            product.Price,
			ComparePrice:     product.ComparePrice,
			CostPrice:        product.CostPrice,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsDigital:        product.IsDigital,
			RequiresShipping: product.RequiresShipping,
			Taxable:          product.Taxable,
			TrackQuantity:    product.TrackQuantity,
			Quantity:         product.Quantity,
			MinQuantity:      product.MinQuantity,
			MaxQuantity:      product.MaxQuantity,
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CreatedAt:        product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        product.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateProductQuantity updates product quantity
func (s *productService) UpdateProductQuantity(ctx context.Context, id int64, quantity int) error {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Update quantity
	err = s.productRepo.UpdateProductQuantity(ctx, id, quantity)
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
	}

	return nil
}

// GetProductsByTags retrieves products by tags
func (s *productService) GetProductsByTags(ctx context.Context, tags []string, page, limit int) (*dto.ListProductsResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Get products from repository
	products, total, err := s.productRepo.GetProductsByTags(ctx, tags, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by tags: %w", err)
	}

	// Convert to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = dto.ProductResponse{
			ID:               product.ID,
			Name:             product.Name,
			Description:      product.Description,
			ShortDesc:        product.ShortDesc,
			SKU:              product.SKU,
			Price:            product.Price,
			ComparePrice:     product.ComparePrice,
			CostPrice:        product.CostPrice,
			Weight:           product.Weight,
			Dimensions:       product.Dimensions,
			IsActive:         product.IsActive,
			IsDigital:        product.IsDigital,
			RequiresShipping: product.RequiresShipping,
			Taxable:          product.Taxable,
			TrackQuantity:    product.TrackQuantity,
			Quantity:         product.Quantity,
			MinQuantity:      product.MinQuantity,
			MaxQuantity:      product.MaxQuantity,
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CreatedAt:        product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        product.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.ListProductsResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Product Variant methods

// CreateProductVariant creates a new product variant
func (s *productService) CreateProductVariant(ctx context.Context, req *dto.CreateProductVariantRequest) (*domain.ProductVariant, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Check if SKU already exists
	existingVariant, err := s.productRepo.GetProductVariantByID(ctx, 0) // This will fail, but we need to check SKU uniqueness
	// For now, we'll skip SKU uniqueness check for variants - in a real app, you'd implement this

	// Create variant domain object
	variant := &domain.ProductVariant{
		ProductID:    req.ProductID,
		Name:         req.Name,
		SKU:          req.SKU,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		CostPrice:    req.CostPrice,
		Weight:       req.Weight,
		Quantity:     req.Quantity,
		IsActive:     req.IsActive,
		Position:     req.Position,
	}

	// Create variant in repository
	err = s.productRepo.CreateProductVariant(ctx, variant)
	if err != nil {
		return nil, fmt.Errorf("failed to create product variant: %w", err)
	}

	_ = existingVariant // Suppress unused variable warning

	return variant, nil
}

// GetProductVariantByID retrieves a product variant by ID
func (s *productService) GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error) {
	variant, err := s.productRepo.GetProductVariantByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variant: %w", err)
	}

	return variant, nil
}

// GetProductVariantsByProductID retrieves all variants for a product
func (s *productService) GetProductVariantsByProductID(ctx context.Context, productID int64) ([]*domain.ProductVariant, error) {
	variants, err := s.productRepo.GetProductVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	return variants, nil
}

// UpdateProductVariant updates an existing product variant
func (s *productService) UpdateProductVariant(ctx context.Context, id int64, req *dto.UpdateProductVariantRequest) (*domain.ProductVariant, error) {
	// Get existing variant
	existingVariant, err := s.productRepo.GetProductVariantByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing variant: %w", err)
	}

	// Update fields that are provided
	updateVariant := *existingVariant

	if req.Name != nil {
		updateVariant.Name = *req.Name
	}
	if req.SKU != nil {
		updateVariant.SKU = *req.SKU
	}
	if req.Price != nil {
		updateVariant.Price = *req.Price
	}
	if req.ComparePrice != nil {
		updateVariant.ComparePrice = *req.ComparePrice
	}
	if req.CostPrice != nil {
		updateVariant.CostPrice = *req.CostPrice
	}
	if req.Weight != nil {
		updateVariant.Weight = *req.Weight
	}
	if req.Quantity != nil {
		updateVariant.Quantity = *req.Quantity
	}
	if req.IsActive != nil {
		updateVariant.IsActive = *req.IsActive
	}
	if req.Position != nil {
		updateVariant.Position = *req.Position
	}

	// Update variant in repository
	err = s.productRepo.UpdateProductVariant(ctx, id, &updateVariant)
	if err != nil {
		return nil, fmt.Errorf("failed to update product variant: %w", err)
	}

	return &updateVariant, nil
}

// DeleteProductVariant deletes a product variant
func (s *productService) DeleteProductVariant(ctx context.Context, id int64) error {
	// Check if variant exists
	_, err := s.productRepo.GetProductVariantByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get variant: %w", err)
	}

	// Delete variant
	err = s.productRepo.DeleteProductVariant(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product variant: %w", err)
	}

	return nil
}

// Product Category methods

// AddProductToCategory adds a product to a category
func (s *productService) AddProductToCategory(ctx context.Context, productID, categoryID int64, isPrimary bool) error {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Add product to category
	err = s.productRepo.AddProductToCategory(ctx, productID, categoryID, isPrimary)
	if err != nil {
		return fmt.Errorf("failed to add product to category: %w", err)
	}

	return nil
}

// RemoveProductFromCategory removes a product from a category
func (s *productService) RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Remove product from category
	err = s.productRepo.RemoveProductFromCategory(ctx, productID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to remove product from category: %w", err)
	}

	return nil
}

// GetProductCategories retrieves all categories for a product
func (s *productService) GetProductCategories(ctx context.Context, productID int64) ([]*domain.Category, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Get categories
	categories, err := s.productRepo.GetProductCategories(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product categories: %w", err)
	}

	return categories, nil
}

// UpdateProductCategories updates all categories for a product
func (s *productService) UpdateProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Update categories
	err = s.productRepo.UpdateProductCategories(ctx, productID, categoryIDs)
	if err != nil {
		return fmt.Errorf("failed to update product categories: %w", err)
	}

	return nil
}
