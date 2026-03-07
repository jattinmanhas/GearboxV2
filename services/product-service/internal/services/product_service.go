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
	GetProductsByTags(ctx context.Context, tags []string, page, limit int) (*dto.ListProductsResponse, error)
	GetProductWithStockInfo(ctx context.Context, product *domain.Product) (*dto.ProductResponse, error)

	// Product Variants
	CreateProductVariant(ctx context.Context, req *dto.CreateProductVariantRequest) (*domain.ProductVariant, error)
	GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error)
	GetProductVariantsByProductID(ctx context.Context, productID int64) ([]*domain.ProductVariant, error)
	GetProductVariantsByProductIDWithInventory(ctx context.Context, productID int64) ([]*dto.ProductVariantResponse, error)
	GetProductVariantsByProductIDAndSKU(ctx context.Context, productID int64, sku string) ([]*domain.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, id int64, req *dto.UpdateProductVariantRequest) (*domain.ProductVariant, error)
	DeleteProductVariant(ctx context.Context, id int64) error

	// Product Categories
	AddProductToCategory(ctx context.Context, productID, categoryID int64, isPrimary bool) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error
	GetProductCategories(ctx context.Context, productID int64) ([]*domain.Category, error)
	UpdateProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error

	// Product Images
	CreateProductImage(ctx context.Context, req *dto.ProductImageRequest, productID int64) (*domain.ProductImage, error)
	GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error)
	UpdateProductImage(ctx context.Context, id int64, req *dto.ProductImageRequest) (*domain.ProductImage, error)
	DeleteProductImage(ctx context.Context, id int64) error
	SetPrimaryProductImage(ctx context.Context, productID, imageID int64) error

	// Product Analytics
	GetProductAnalytics(ctx context.Context) (*dto.ProductAnalyticsResponse, error)
	GetTopSellingProducts(ctx context.Context, limit int) ([]*dto.ProductOrderStatsResponse, error)
}

type productService struct {
	productRepo      repository.ProductRepository
	inventoryService InventoryService
}

func NewProductService(productRepo repository.ProductRepository, inventoryService InventoryService) ProductService {
	return &productService{
		productRepo:      productRepo,
		inventoryService: inventoryService,
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
			fmt.Printf("Warning: failed to add product to categories: %v\n", err)
		}
	}

	// Add product images if provided
	if len(req.Images) > 0 {
		err = s.createProductImages(ctx, product.ID, req.Images)
		if err != nil {
			// Log error but don't fail the product creation
			fmt.Printf("Warning: failed to add product images: %v\n", err)
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

	// Update images if provided
	if req.Images != nil {
		// Delete existing images
		existingImages, err := s.productRepo.GetProductImages(ctx, id)
		if err == nil {
			for _, img := range existingImages {
				err := s.productRepo.DeleteProductImage(ctx, img.ID)
				if err != nil {
					fmt.Printf("Warning: failed to delete existing product image: %v\n", err)
				}
			}
		}

		// Create new images
		if len(req.Images) > 0 {
			err = s.createProductImages(ctx, id, req.Images)
			if err != nil {
				// Log error but don't fail the product update
				fmt.Printf("Warning: failed to add product images: %v\n", err)
			}
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
		req.Limit = 10
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
		OnSale:     req.OnSale,
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
		// Calculate stock information for products without variants
		quantity := 0
		availableQuantity := 0
		isInStock := false

		// Check if product has variants
		hasVariants, err := s.productRepo.HasVariants(ctx, product.ID)
		if err != nil {
			// If we can't determine if product has variants, assume it doesn't
			hasVariants = false
		}

		if !hasVariants {
			// Product has no variants, check inventory directly
			if product.IsDigital {
				// Digital products are always in stock and don't need inventory tracking
				quantity = 999
				availableQuantity = 999
				isInStock = true
			} else {
				inventory, err := s.inventoryService.GetInventoryByProduct(ctx, product.ID, nil)
				if err == nil {
					quantity = inventory.Quantity
					availableQuantity = inventory.AvailableQuantity
					isInStock = inventory.AvailableQuantity > 0
				} else {
					// If no inventory record exists, it's out of stock
					quantity = 0
					availableQuantity = 0
					isInStock = false
				}
			}
		} else {
			// Product has variants, check if any variant is in stock
			if product.IsDigital {
				// Digital products with variants are always in stock
				quantity = 999
				availableQuantity = 999
				isInStock = true
			} else {
				variants, err := s.productRepo.GetProductVariants(ctx, product.ID)
				if err == nil {
					for _, variant := range variants {
						if variant.IsActive {
							inventory, err := s.inventoryService.GetInventoryByProduct(ctx, product.ID, &variant.ID)
							if err == nil && inventory.AvailableQuantity > 0 {
								isInStock = true
								break
							}
						}
					}
				}
			}
		}

		productResponses[i] = dto.ProductResponse{
			ID:                product.ID,
			Name:              product.Name,
			Description:       product.Description,
			ShortDesc:         product.ShortDesc,
			SKU:               product.SKU,
			Price:             product.Price,
			ComparePrice:      product.ComparePrice,
			CostPrice:         product.CostPrice,
			Weight:            product.Weight,
			Dimensions:        product.Dimensions,
			IsActive:          product.IsActive,
			IsDigital:         product.IsDigital,
			RequiresShipping:  product.RequiresShipping,
			MetaTitle:         product.MetaTitle,
			MetaDescription:   product.MetaDesc,
			Tags:              product.Tags,
			CategoryIDs:       product.CategoryIDs,
			CategoryNames:     product.CategoryNames,
			Images:            convertDomainImagesToDTO(product.Images),
			CreatedAt:         product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         product.UpdatedAt.Format(time.RFC3339),
			Quantity:          quantity,
			AvailableQuantity: availableQuantity,
			IsInStock:         isInStock,
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

// convertDomainImagesToDTO converts domain images to DTO images
func convertDomainImagesToDTO(domainImages []domain.ProductImage) []dto.ProductImageResponse {
	images := make([]dto.ProductImageResponse, len(domainImages))
	for i, img := range domainImages {
		images[i] = dto.ProductImageResponse{
			ID:        img.ID,
			ProductID: img.ProductID,
			URL:       img.URL,
			Alt:       img.Alt,
			Position:  img.Position,
			IsPrimary: img.IsPrimary,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}
	return images
}

// GetProductWithStockInfo gets a product with stock information calculated
func (s *productService) GetProductWithStockInfo(ctx context.Context, product *domain.Product) (*dto.ProductResponse, error) {
	// Calculate stock information for products without variants
	quantity := 0
	availableQuantity := 0
	isInStock := false

	// Check if product has variants
	hasVariants, err := s.productRepo.HasVariants(ctx, product.ID)
	if err != nil {
		// If we can't determine if product has variants, assume it doesn't
		hasVariants = false
	}

	if !hasVariants {
		// Product has no variants, check inventory directly
		if product.IsDigital {
			// Digital products are always in stock and don't need inventory tracking
			quantity = 999
			availableQuantity = 999
			isInStock = true
		} else {
			inventory, err := s.inventoryService.GetInventoryByProduct(ctx, product.ID, nil)
			if err == nil {
				quantity = inventory.Quantity
				availableQuantity = inventory.AvailableQuantity
				isInStock = inventory.AvailableQuantity > 0
			} else {
				// If no inventory record exists, it's out of stock
				quantity = 0
				availableQuantity = 0
				isInStock = false
			}
		}
	} else {
		// Product has variants, check if any variant is in stock
		if product.IsDigital {
			// Digital products with variants are always in stock
			quantity = 999
			availableQuantity = 999
			isInStock = true
		} else {
			variants, err := s.productRepo.GetProductVariants(ctx, product.ID)
			if err == nil {
				for _, variant := range variants {
					if variant.IsActive {
						inventory, err := s.inventoryService.GetInventoryByProduct(ctx, product.ID, &variant.ID)
						if err == nil && inventory.AvailableQuantity > 0 {
							isInStock = true
							break
						}
					}
				}
			}
		}
	}

	// Convert domain images to DTO images
	images := convertDomainImagesToDTO(product.Images)

	return &dto.ProductResponse{
		ID:                product.ID,
		Name:              product.Name,
		Description:       product.Description,
		ShortDesc:         product.ShortDesc,
		SKU:               product.SKU,
		Price:             product.Price,
		ComparePrice:      product.ComparePrice,
		CostPrice:         product.CostPrice,
		Weight:            product.Weight,
		Dimensions:        product.Dimensions,
		IsActive:          product.IsActive,
		IsDigital:         product.IsDigital,
		RequiresShipping:  product.RequiresShipping,
		MetaTitle:         product.MetaTitle,
		MetaDescription:   product.MetaDesc,
		Tags:              product.Tags,
		CategoryIDs:       product.CategoryIDs,
		CategoryNames:     product.CategoryNames,
		Images:            images,
		CreatedAt:         product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         product.UpdatedAt.Format(time.RFC3339),
		Quantity:          quantity,
		AvailableQuantity: availableQuantity,
		IsInStock:         isInStock,
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
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CategoryIDs:      product.CategoryIDs,
			CategoryNames:    product.CategoryNames,
			Images:           convertDomainImagesToDTO(product.Images),
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
		limit = 10
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
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CategoryIDs:      product.CategoryIDs,
			CategoryNames:    product.CategoryNames,
			Images:           convertDomainImagesToDTO(product.Images),
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

// GetProductsByTags retrieves products by tags
func (s *productService) GetProductsByTags(ctx context.Context, tags []string, page, limit int) (*dto.ListProductsResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
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
			MetaTitle:        product.MetaTitle,
			MetaDescription:  product.MetaDesc,
			Tags:             product.Tags,
			CategoryIDs:      product.CategoryIDs,
			CategoryNames:    product.CategoryNames,
			Images:           convertDomainImagesToDTO(product.Images),
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
	existingVariant, err := s.productRepo.GetProductVariantsByProductIDAndSKU(ctx, req.ProductID, req.SKU)
	if err == nil && existingVariant != nil {
		return nil, fmt.Errorf("variant with SKU %s already exists", req.SKU)
	}

	// Create variant domain object
	variant := &domain.ProductVariant{
		ProductID:    req.ProductID,
		Name:         req.Name,
		SKU:          req.SKU,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		CostPrice:    req.CostPrice,
		Weight:       req.Weight,
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

// GetProductVariantsByProductIDWithInventory retrieves all variants for a product with inventory data
func (s *productService) GetProductVariantsByProductIDWithInventory(ctx context.Context, productID int64) ([]*dto.ProductVariantResponse, error) {
	// Get product information to check track_quantity setting
	product, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	variants, err := s.productRepo.GetProductVariantsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	// Convert to response DTOs with inventory data
	variantResponses := make([]*dto.ProductVariantResponse, len(variants))
	for i, variant := range variants {
		var quantity int
		var availableQuantity int
		var isInStock bool

		// Handle digital products - they are always in stock
		if product.IsDigital {
			quantity = 999
			availableQuantity = 999
			isInStock = true
		} else {
			// For physical products, check inventory
			inventory, err := s.inventoryService.GetInventoryByProduct(ctx, productID, &variant.ID)
			if err != nil {
				// If no inventory record exists, it's out of stock
				quantity = 0
				availableQuantity = 0
				isInStock = false
			} else {
				// Use actual inventory data
				quantity = inventory.Quantity
				availableQuantity = inventory.AvailableQuantity
				isInStock = inventory.AvailableQuantity > 0
			}
		}

		variantResponses[i] = &dto.ProductVariantResponse{
			ID:                variant.ID,
			ProductID:         variant.ProductID,
			Name:              variant.Name,
			SKU:               variant.SKU,
			Price:             variant.Price,
			ComparePrice:      variant.ComparePrice,
			CostPrice:         variant.CostPrice,
			Weight:            variant.Weight,
			IsActive:          variant.IsActive,
			Position:          variant.Position,
			Quantity:          quantity,
			AvailableQuantity: availableQuantity,
			IsInStock:         isInStock,
		}
	}

	return variantResponses, nil
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

// GetProductVariantsByProductIDAndSKU retrieves all variants for a product and SKU
func (s *productService) GetProductVariantsByProductIDAndSKU(ctx context.Context, productID int64, sku string) ([]*domain.ProductVariant, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Get variants
	variants, err := s.productRepo.GetProductVariantsByProductIDAndSKU(ctx, productID, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	return variants, nil
}

// Product Image methods

// createProductImages is a helper method to create multiple product images
func (s *productService) createProductImages(ctx context.Context, productID int64, imageRequests []dto.ProductImageRequest) error {
	for i, imageReq := range imageRequests {
		image := &domain.ProductImage{
			ProductID: productID,
			URL:       imageReq.URL,
			Alt:       imageReq.Alt,
			Position:  imageReq.Position,
			IsPrimary: imageReq.IsPrimary,
		}

		// If position is not set, use the index
		if image.Position == 0 {
			image.Position = i + 1
		}

		// If this is the first image and no primary is set, make it primary
		if i == 0 && !imageReq.IsPrimary {
			image.IsPrimary = true
		}

		err := s.productRepo.CreateProductImage(ctx, image)
		if err != nil {
			return fmt.Errorf("failed to create product image: %w", err)
		}
	}

	return nil
}

// CreateProductImage creates a new product image
func (s *productService) CreateProductImage(ctx context.Context, req *dto.ProductImageRequest, productID int64) (*domain.ProductImage, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	image := &domain.ProductImage{
		ProductID: productID,
		URL:       req.URL,
		Alt:       req.Alt,
		Position:  req.Position,
		IsPrimary: req.IsPrimary,
	}

	err = s.productRepo.CreateProductImage(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("failed to create product image: %w", err)
	}

	return image, nil
}

// GetProductImages retrieves all images for a product
func (s *productService) GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error) {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	images, err := s.productRepo.GetProductImages(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}

	return images, nil
}

// UpdateProductImage updates an existing product image
func (s *productService) UpdateProductImage(ctx context.Context, id int64, req *dto.ProductImageRequest) (*domain.ProductImage, error) {
	image := &domain.ProductImage{
		URL:       req.URL,
		Alt:       req.Alt,
		Position:  req.Position,
		IsPrimary: req.IsPrimary,
	}

	err := s.productRepo.UpdateProductImage(ctx, id, image)
	if err != nil {
		return nil, fmt.Errorf("failed to update product image: %w", err)
	}

	// Get the updated image
	updatedImage, err := s.productRepo.GetProductImages(ctx, image.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated image: %w", err)
	}

	// Find the updated image by ID
	for _, img := range updatedImage {
		if img.ID == id {
			return img, nil
		}
	}

	return nil, fmt.Errorf("updated image not found")
}

// DeleteProductImage deletes a product image
func (s *productService) DeleteProductImage(ctx context.Context, id int64) error {
	err := s.productRepo.DeleteProductImage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product image: %w", err)
	}

	return nil
}

// SetPrimaryProductImage sets a product image as primary
func (s *productService) SetPrimaryProductImage(ctx context.Context, productID, imageID int64) error {
	// Check if product exists
	_, err := s.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	err = s.productRepo.SetPrimaryProductImage(ctx, productID, imageID)
	if err != nil {
		return fmt.Errorf("failed to set primary product image: %w", err)
	}

	return nil
}

// Product Analytics

// GetProductAnalytics retrieves product analytics
func (s *productService) GetProductAnalytics(ctx context.Context) (*dto.ProductAnalyticsResponse, error) {
	analytics, err := s.productRepo.GetProductAnalytics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get product analytics: %w", err)
	}

	return &dto.ProductAnalyticsResponse{
		TotalProducts:       analytics.TotalProducts,
		ActiveProducts:      analytics.ActiveProducts,
		InactiveProducts:    analytics.InactiveProducts,
		LowStockProducts:    analytics.LowStockProducts,
		OutOfStockProducts:  analytics.OutOfStockProducts,
		TotalCategories:     analytics.TotalCategories,
		AveragePrice:        analytics.AveragePrice,
		TotalInventoryValue: analytics.TotalInventoryValue,
	}, nil
}

// GetTopSellingProducts retrieves top selling products
func (s *productService) GetTopSellingProducts(ctx context.Context, limit int) ([]*dto.ProductOrderStatsResponse, error) {
	products, err := s.productRepo.GetTopSellingProducts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top selling products: %w", err)
	}

	responses := make([]*dto.ProductOrderStatsResponse, len(products))
	for i, product := range products {
		responses[i] = &dto.ProductOrderStatsResponse{
			ProductID:     product.ProductID,
			ProductName:   product.ProductName,
			SKU:           product.SKU,
			TotalQuantity: product.TotalQuantity,
			TotalRevenue:  product.TotalRevenue,
			OrderCount:    product.OrderCount,
			AveragePrice:  product.AveragePrice,
		}
	}

	return responses, nil
}
