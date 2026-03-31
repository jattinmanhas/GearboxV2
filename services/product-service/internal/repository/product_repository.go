package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProductByID(ctx context.Context, id int64) (*domain.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id int64, product *domain.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context, filter *domain.ProductFilter, offset, limit int) ([]*domain.Product, int64, error)
	GetProductsByCategory(ctx context.Context, categoryID int64, offset, limit int) ([]*domain.Product, int64, error)
	SearchProducts(ctx context.Context, query string, offset, limit int) ([]*domain.Product, int64, error)
	GetProductsByTags(ctx context.Context, tags []string, offset, limit int) ([]*domain.Product, int64, error)

	// Product Variants
	CreateProductVariant(ctx context.Context, variant *domain.ProductVariant) error
	GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error)
	GetProductVariantsByProductID(ctx context.Context, productID int64) ([]*domain.ProductVariant, error)
	GetProductVariants(ctx context.Context, productID int64) ([]*domain.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, id int64, variant *domain.ProductVariant) error
	DeleteProductVariant(ctx context.Context, id int64) error
	GetProductVariantsByProductIDAndSKU(ctx context.Context, productID int64, sku string) ([]*domain.ProductVariant, error)
	HasVariants(ctx context.Context, productID int64) (bool, error)

	// Product Categories
	AddProductToCategory(ctx context.Context, productID, categoryID int64, isPrimary bool) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error
	GetProductCategories(ctx context.Context, productID int64) ([]*domain.Category, error)
	UpdateProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error
	CheckCategoryHasProducts(ctx context.Context, categoryID int64) (bool, error)

	// Product Images
	CreateProductImage(ctx context.Context, image *domain.ProductImage) error
	GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error)
	UpdateProductImage(ctx context.Context, id int64, image *domain.ProductImage) error
	DeleteProductImage(ctx context.Context, id int64) error
	DeleteProductImages(ctx context.Context, productID int64) error
	SetPrimaryProductImage(ctx context.Context, productID, imageID int64) error

	// Product Analytics
	GetProductAnalytics(ctx context.Context) (*domain.ProductAnalytics, error)
	GetTopSellingProducts(ctx context.Context, limit int) ([]*domain.ProductOrderStats, error)
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (
			name, description, short_description, sku, price, compare_price, cost_price,
			weight, dimensions, is_active, is_digital, requires_shipping, meta_title,
			meta_description, tags, created_at, updated_at
		) VALUES (
			:name, :description, :short_description, :sku, :price, :compare_price, :cost_price,
			:weight, :dimensions, :is_active, :is_digital, :requires_shipping, :meta_title,
			:meta_description, :tags, :created_at, :updated_at
		)
		RETURNING id
	`

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	rows, err := r.db.NamedQueryContext(ctx, query, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&product.ID); err != nil {
			return fmt.Errorf("failed to scan product ID: %w", err)
		}
	}

	return nil
}

// GetProductByID retrieves a product by ID
func (r *productRepository) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT * FROM products WHERE id = $1`

	var product domain.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Populate category information
	categoryIDs, categoryNames, err := r.getProductCategories(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories for product %d: %w", product.ID, err)
	}
	product.CategoryIDs = categoryIDs
	product.CategoryNames = categoryNames

	// Populate image information
	images, err := r.GetProductImages(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for product %d: %w", product.ID, err)
	}
	// Convert []*ProductImage to []ProductImage
	productImages := make([]domain.ProductImage, len(images))
	for i, img := range images {
		productImages[i] = *img
	}
	product.Images = productImages

	return &product, nil
}

// GetProductBySKU retrieves a product by SKU
func (r *productRepository) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `SELECT * FROM products WHERE sku = $1`

	var product domain.Product
	err := r.db.GetContext(ctx, &product, query, sku)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with SKU %s not found", sku)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Populate category information
	categoryIDs, categoryNames, err := r.getProductCategories(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories for product %d: %w", product.ID, err)
	}
	product.CategoryIDs = categoryIDs
	product.CategoryNames = categoryNames

	// Populate image information
	images, err := r.GetProductImages(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get images for product %d: %w", product.ID, err)
	}
	// Convert []*ProductImage to []ProductImage
	productImages := make([]domain.ProductImage, len(images))
	for i, img := range images {
		productImages[i] = *img
	}
	product.Images = productImages

	return &product, nil
}

// UpdateProduct updates an existing product
func (r *productRepository) UpdateProduct(ctx context.Context, id int64, product *domain.Product) error {
	query := `
		UPDATE products SET
			name = :name, description = :description, short_description = :short_description,
			sku = :sku, price = :price, compare_price = :compare_price, cost_price = :cost_price,
			weight = :weight, dimensions = :dimensions, is_active = :is_active,
			is_digital = :is_digital, requires_shipping = :requires_shipping,
			meta_title = :meta_title, meta_description = :meta_description,
			tags = :tags, updated_at = :updated_at
		WHERE id = :id`

	product.UpdatedAt = time.Now()
	product.ID = id

	result, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

// DeleteProduct deletes a product
func (r *productRepository) DeleteProduct(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}

	return nil
}

// ListProducts retrieves products with filters
func (r *productRepository) ListProducts(ctx context.Context, filter *domain.ProductFilter, offset, limit int) ([]*domain.Product, int64, error) {
	whereClause, args := r.buildWhereClause(filter)

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)

	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// List query
	orderClause := r.buildOrderClause(filter)
	query := fmt.Sprintf(`
    SELECT * FROM products 
    %s 
    %s 
    LIMIT $%d OFFSET $%d`, whereClause, orderClause, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	// Populate category and image information for each product
	for _, product := range products {
		categoryIDs, categoryNames, err := r.getProductCategories(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get categories for product %d: %w", product.ID, err)
		}
		product.CategoryIDs = categoryIDs
		product.CategoryNames = categoryNames

		// Populate image information
		images, err := r.GetProductImages(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get images for product %d: %w", product.ID, err)
		}
		// Convert []*ProductImage to []ProductImage
		productImages := make([]domain.ProductImage, len(images))
		for i, img := range images {
			productImages[i] = *img
		}
		product.Images = productImages
	}

	return products, total, nil
}

// getProductCategories retrieves category IDs and names for a product
func (r *productRepository) getProductCategories(ctx context.Context, productID int64) ([]int64, []string, error) {
	query := `
		SELECT c.id, c.name 
		FROM categories c
		INNER JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
		ORDER BY c.name`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get product categories: %w", err)
	}
	defer rows.Close()

	var categoryIDs []int64
	var categoryNames []string

	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categoryIDs = append(categoryIDs, id)
		categoryNames = append(categoryNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categoryIDs, categoryNames, nil
}

// GetProductsByCategory retrieves products by category
func (r *productRepository) GetProductsByCategory(ctx context.Context, categoryID int64, offset, limit int) ([]*domain.Product, int64, error) {
	// Count query
	countQuery := `
		SELECT COUNT(DISTINCT p.id) 
		FROM products p 
		INNER JOIN product_categories pc ON p.id = pc.product_id 
		WHERE pc.category_id = ?`
	countQuery = r.db.Rebind(countQuery)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products by category: %w", err)
	}

	// List query
	query := `
		SELECT DISTINCT p.* 
		FROM products p 
		INNER JOIN product_categories pc ON p.id = pc.product_id 
		WHERE pc.category_id = ? 
		ORDER BY p.created_at DESC 
		LIMIT ? OFFSET ?`
	query = r.db.Rebind(query)
	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
	}

	// Populate category and image information for each product
	for _, product := range products {
		categoryIDs, categoryNames, err := r.getProductCategories(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get categories for product %d: %w", product.ID, err)
		}
		product.CategoryIDs = categoryIDs
		product.CategoryNames = categoryNames

		// Populate image information
		images, err := r.GetProductImages(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get images for product %d: %w", product.ID, err)
		}
		// Convert []*ProductImage to []ProductImage
		productImages := make([]domain.ProductImage, len(images))
		for i, img := range images {
			productImages[i] = *img
		}
		product.Images = productImages
	}

	return products, total, nil
}

// SearchProducts searches products by name, description, or SKU
func (r *productRepository) SearchProducts(ctx context.Context, query string, offset, limit int) ([]*domain.Product, int64, error) {
	searchTerm := "%" + query + "%"

	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM products 
		WHERE name ILIKE $1 OR description ILIKE $2 OR sku ILIKE $3 OR tags ILIKE $4`

	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Search query
	searchQuery := `
		SELECT * 
		FROM products 
		WHERE name ILIKE $1 OR description ILIKE $2 OR sku ILIKE $3 OR tags ILIKE $4
		ORDER BY 
			CASE 
				WHEN name ILIKE $5 THEN 1
				WHEN sku ILIKE $6 THEN 2
				WHEN description ILIKE $7 THEN 3
				ELSE 4
			END,
			name
		LIMIT $8 OFFSET $9`
	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, searchQuery,
		searchTerm, searchTerm, searchTerm, searchTerm,
		searchTerm, searchTerm, searchTerm,
		limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	// Populate category and image information for each product
	for _, product := range products {
		categoryIDs, categoryNames, err := r.getProductCategories(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get categories for product %d: %w", product.ID, err)
		}
		product.CategoryIDs = categoryIDs
		product.CategoryNames = categoryNames

		// Populate image information
		images, err := r.GetProductImages(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get images for product %d: %w", product.ID, err)
		}
		// Convert []*ProductImage to []ProductImage
		productImages := make([]domain.ProductImage, len(images))
		for i, img := range images {
			productImages[i] = *img
		}
		product.Images = productImages
	}

	return products, total, nil
}

// GetProductsByTags retrieves products by tags
func (r *productRepository) GetProductsByTags(ctx context.Context, tags []string, offset, limit int) ([]*domain.Product, int64, error) {
	if len(tags) == 0 {
		return []*domain.Product{}, 0, nil
	}

	// Build tag conditions
	var tagConditions []string
	var args []interface{}

	for _, tag := range tags {
		tagConditions = append(tagConditions, "tags LIKE ?")
		args = append(args, "%"+tag+"%")
	}

	whereClause := "WHERE " + strings.Join(tagConditions, " OR ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	countQuery = r.db.Rebind(countQuery)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products by tags: %w", err)
	}

	// List query
	query := fmt.Sprintf(`
		SELECT * FROM products 
		%s 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, whereClause)
	query = r.db.Rebind(query)
	args = append(args, limit, offset)

	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products by tags: %w", err)
	}

	// Populate category and image information for each product
	for _, product := range products {
		categoryIDs, categoryNames, err := r.getProductCategories(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get categories for product %d: %w", product.ID, err)
		}
		product.CategoryIDs = categoryIDs
		product.CategoryNames = categoryNames

		// Populate image information
		images, err := r.GetProductImages(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get images for product %d: %w", product.ID, err)
		}
		// Convert []*ProductImage to []ProductImage
		productImages := make([]domain.ProductImage, len(images))
		for i, img := range images {
			productImages[i] = *img
		}
		product.Images = productImages
	}

	return products, total, nil
}

// Product Variant methods

// CreateProductVariant creates a new product variant
func (r *productRepository) CreateProductVariant(ctx context.Context, variant *domain.ProductVariant) error {
	query := `
		INSERT INTO product_variants (
			product_id, name, sku, price, compare_price, cost_price,
			weight, is_active, position
		) VALUES (
			:product_id, :name, :sku, :price, :compare_price, :cost_price,
			:weight, :is_active, :position
		)
		RETURNING id`

	// Use NamedQueryContext to fetch the generated id
	rows, err := r.db.NamedQueryContext(ctx, query, variant)
	if err != nil {
		return fmt.Errorf("failed to create product variant: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&variant.ID); err != nil {
			return fmt.Errorf("failed to scan variant ID: %w", err)
		}
	}

	return nil
}

// GetProductVariantByID retrieves a product variant by ID
func (r *productRepository) GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error) {
	query := `SELECT * FROM product_variants WHERE id = $1`

	var variant domain.ProductVariant
	err := r.db.GetContext(ctx, &variant, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product variant with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product variant: %w", err)
	}

	return &variant, nil
}

// GetProductVariantsByProductID retrieves all variants for a product
func (r *productRepository) GetProductVariantsByProductID(ctx context.Context, productID int64) ([]*domain.ProductVariant, error) {
	query := `SELECT * FROM product_variants WHERE product_id = $1 ORDER BY position, name`

	var variants []*domain.ProductVariant
	err := r.db.SelectContext(ctx, &variants, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	return variants, nil
}

// UpdateProductVariant updates an existing product variant
func (r *productRepository) UpdateProductVariant(ctx context.Context, id int64, variant *domain.ProductVariant) error {
	query := `
		UPDATE product_variants SET
			name = :name, sku = :sku, price = :price, compare_price = :compare_price,
			cost_price = :cost_price, weight = :weight,
			is_active = :is_active, position = :position
		WHERE id = :id`

	variant.ID = id

	result, err := r.db.NamedExecContext(ctx, query, variant)
	if err != nil {
		return fmt.Errorf("failed to update product variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product variant with ID %d not found", id)
	}

	return nil
}

// DeleteProductVariant deletes a product variant
func (r *productRepository) DeleteProductVariant(ctx context.Context, id int64) error {
	query := `DELETE FROM product_variants WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product variant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product variant with ID %d not found", id)
	}

	return nil
}

// Product Category methods

// AddProductToCategory adds a product to a category
func (r *productRepository) AddProductToCategory(ctx context.Context, productID, categoryID int64, isPrimary bool) error {
	query := `
		INSERT INTO product_categories (product_id, category_id, is_primary)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, productID, categoryID, isPrimary)
	if err != nil {
		return fmt.Errorf("failed to add product to category: %w", err)
	}

	return nil
}

// RemoveProductFromCategory removes a product from a category
func (r *productRepository) RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error {
	query := `DELETE FROM product_categories WHERE product_id = $1 AND category_id = $2`

	result, err := r.db.ExecContext(ctx, query, productID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to remove product from category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product-category relationship not found")
	}

	return nil
}

// GetProductCategories retrieves all categories for a product
func (r *productRepository) GetProductCategories(ctx context.Context, productID int64) ([]*domain.Category, error) {
	query := `
		SELECT c.* 
		FROM categories c 
		INNER JOIN product_categories pc ON c.id = pc.category_id 
		WHERE pc.product_id = $1 
		ORDER BY pc.is_primary DESC, c.name`

	var categories []*domain.Category
	err := r.db.SelectContext(ctx, &categories, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product categories: %w", err)
	}

	return categories, nil
}

// UpdateProductCategories updates all categories for a product
func (r *productRepository) UpdateProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Remove existing categories
	_, err = tx.ExecContext(ctx, "DELETE FROM product_categories WHERE product_id = $1", productID)
	if err != nil {
		return fmt.Errorf("failed to remove existing categories: %w", err)
	}

	// Add new categories
	for i, categoryID := range categoryIDs {
		isPrimary := i == 0 // First category is primary
		_, err = tx.ExecContext(ctx,
			"INSERT INTO product_categories (product_id, category_id, is_primary) VALUES ($1, $2, $3)",
			productID, categoryID, isPrimary)
		if err != nil {
			return fmt.Errorf("failed to add category: %w", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Helper methods

// buildWhereClause builds the WHERE clause for product filtering
func (r *productRepository) buildWhereClause(filter *domain.ProductFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("id IN (SELECT product_id FROM product_categories WHERE category_id = $%d)", argIndex))
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.IsDigital != nil {
		conditions = append(conditions, fmt.Sprintf("is_digital = $%d", argIndex))
		args = append(args, *filter.IsDigital)
		argIndex++
	}

	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *filter.MaxPrice)
		argIndex++
	}

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR sku ILIKE $%d OR tags ILIKE $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3))
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
		argIndex += 4
	}

	if len(filter.Tags) > 0 {
		var tagConditions []string
		for _, tag := range filter.Tags {
			tagConditions = append(tagConditions, fmt.Sprintf("tags ILIKE $%d", argIndex))
			args = append(args, "%"+tag+"%")
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " OR ")+")")
	}

	if filter.InStock != nil {
		if *filter.InStock {
			// Filter for products that are in stock:
			// 1. Digital products (always in stock)
			// 2. Products with variants that have at least one in-stock variant
			// 3. Products without variants that have inventory with available_quantity > 0
			conditions = append(conditions, `(
				is_digital = true OR
				id IN (
					SELECT DISTINCT pv.product_id FROM product_variants pv
					INNER JOIN inventory i ON pv.id = i.product_variant_id
					WHERE pv.is_active = true AND i.available_quantity > 0
				) OR (
					id NOT IN (SELECT DISTINCT product_id FROM product_variants) 
					AND id IN (
						SELECT product_id FROM inventory 
						WHERE product_variant_id IS NULL 
						AND available_quantity > 0
					)
				)
			)`)
		} else {
			// Filter for products that are out of stock:
			// 1. Non-digital products with variants that have no in-stock variants
			// 2. Non-digital products without variants that have no inventory or available_quantity = 0
			conditions = append(conditions, `(
				is_digital = false AND (
					id NOT IN (
						SELECT DISTINCT pv.product_id FROM product_variants pv
						INNER JOIN inventory i ON pv.id = i.product_variant_id
						WHERE pv.is_active = true AND i.available_quantity > 0
					) AND (
						id IN (SELECT DISTINCT product_id FROM product_variants) OR
						id NOT IN (
							SELECT product_id FROM inventory 
							WHERE product_variant_id IS NULL 
							AND available_quantity > 0
						)
					)
				)
			)`)
		}
	}

	if filter.OnSale != nil {
		if *filter.OnSale {
			// Filter for products that are on sale (compare_price > price)
			conditions = append(conditions, "compare_price > price AND compare_price > 0")
		} else {
			// Filter for products that are not on sale
			conditions = append(conditions, "(compare_price <= price OR compare_price = 0 OR compare_price IS NULL)")
		}
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}

// buildOrderClause builds the ORDER BY clause for product sorting
func (r *productRepository) buildOrderClause(filter *domain.ProductFilter) string {
	sortBy := "created_at"
	sortOrder := "DESC"

	if filter.SortBy != "" {
		validSortFields := map[string]string{
			"name":       "name",
			"price":      "price",
			"created_at": "created_at",
			"updated_at": "updated_at",
			"sku":        "sku",
		}
		if field, ok := validSortFields[filter.SortBy]; ok {
			sortBy = field
		}
	}

	if filter.SortOrder != "" {
		if filter.SortOrder == "asc" {
			sortOrder = "ASC"
		}
	}

	return fmt.Sprintf("ORDER BY %s %s", sortBy, sortOrder)
}

// GetProductVariantsByProductIDAndSKU retrieves all variants for a product and SKU
func (r *productRepository) GetProductVariantsByProductIDAndSKU(ctx context.Context, productID int64, sku string) ([]*domain.ProductVariant, error) {
	query := `SELECT * FROM product_variants WHERE product_id = $1 AND sku = $2`

	var variants []*domain.ProductVariant
	err := r.db.SelectContext(ctx, &variants, query, productID, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	return variants, nil
}

// GetProductVariants is an alias for GetProductVariantsByProductID for consistency
func (r *productRepository) GetProductVariants(ctx context.Context, productID int64) ([]*domain.ProductVariant, error) {
	return r.GetProductVariantsByProductID(ctx, productID)
}

// HasVariants checks if a product has any variants
func (r *productRepository) HasVariants(ctx context.Context, productID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM product_variants WHERE product_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, query, productID)
	if err != nil {
		return false, fmt.Errorf("failed to check if product has variants: %w", err)
	}
	return count > 0, nil
}

func (r *productRepository) CheckCategoryHasProducts(ctx context.Context, categoryID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM product_categories WHERE category_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, query, categoryID)
	if err != nil {
		return false, fmt.Errorf("failed to check category has products: %w", err)
	}
	return count > 0, nil
}

// Product Image methods

// CreateProductImage creates a new product image
func (r *productRepository) CreateProductImage(ctx context.Context, image *domain.ProductImage) error {
	query := `
		INSERT INTO product_images (product_id, url, alt, public_id, position, is_primary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query,
		image.ProductID, image.URL, image.Alt, image.PublicID, image.Position, image.IsPrimary,
	).Scan(&image.ID, &image.CreatedAt, &image.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product image: %w", err)
	}

	return nil
}

// GetProductImages retrieves all images for a product
func (r *productRepository) GetProductImages(ctx context.Context, productID int64) ([]*domain.ProductImage, error) {
	query := `
		SELECT id, product_id, url, alt, public_id, position, is_primary, created_at, updated_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY position ASC, created_at
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return []*domain.ProductImage{}, fmt.Errorf("failed to query product images: %w", err)
	}
	defer rows.Close()

	var images []*domain.ProductImage
	for rows.Next() {
		var image domain.ProductImage
		err := rows.Scan(
			&image.ID, &image.ProductID, &image.URL, &image.Alt, &image.PublicID,
			&image.Position, &image.IsPrimary, &image.CreatedAt, &image.UpdatedAt,
		)
		if err != nil {
			return []*domain.ProductImage{}, fmt.Errorf("failed to scan product image: %w", err)
		}
		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return []*domain.ProductImage{}, nil
	}

	return images, nil
}

// UpdateProductImage updates an existing product image
func (r *productRepository) UpdateProductImage(ctx context.Context, id int64, image *domain.ProductImage) error {
	query := `
		UPDATE product_images SET
			url = :url, alt = :alt, position = :position, is_primary = :is_primary
		WHERE id = :id`

	image.ID = id

	result, err := r.db.NamedExecContext(ctx, query, image)
	if err != nil {
		return fmt.Errorf("failed to update product image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product image with ID %d not found", id)
	}

	return nil
}

// DeleteProductImage deletes a product image
func (r *productRepository) DeleteProductImage(ctx context.Context, id int64) error {
	query := `DELETE FROM product_images WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product image: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product image with ID %d not found", id)
	}

	return nil
}

// DeleteProductImages deletes all images for a product
func (r *productRepository) DeleteProductImages(ctx context.Context, productID int64) error {
	query := `DELETE FROM product_images WHERE product_id = $1`

	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product images: %w", err)
	}

	return nil
}

// SetPrimaryProductImage sets a product image as primary (and unsets others)
func (r *productRepository) SetPrimaryProductImage(ctx context.Context, productID, imageID int64) error {
	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Unset all primary images for this product
	_, err = tx.ExecContext(ctx, "UPDATE product_images SET is_primary = false WHERE product_id = $1", productID)
	if err != nil {
		return fmt.Errorf("failed to unset primary images: %w", err)
	}

	// Set the specified image as primary
	_, err = tx.ExecContext(ctx, "UPDATE product_images SET is_primary = true WHERE id = $1 AND product_id = $2", imageID, productID)
	if err != nil {
		return fmt.Errorf("failed to set primary image: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Product Analytics

// GetProductAnalytics retrieves product analytics
func (r *productRepository) GetProductAnalytics(ctx context.Context) (*domain.ProductAnalytics, error) {
	analytics := &domain.ProductAnalytics{}

	// Get total products count
	err := r.db.GetContext(ctx, &analytics.TotalProducts, `SELECT COUNT(*) FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total products count: %w", err)
	}

	// Get active/inactive products count
	var statusCounts struct {
		Active   int64 `db:"active"`
		Inactive int64 `db:"inactive"`
	}

	err = r.db.GetContext(ctx, &statusCounts, `
		SELECT 
			COUNT(CASE WHEN is_active = true THEN 1 END) as active,
			COUNT(CASE WHEN is_active = false THEN 1 END) as inactive
		FROM products`)
	if err != nil {
		return nil, fmt.Errorf("failed to get active/inactive products count: %w", err)
	}
	analytics.ActiveProducts = statusCounts.Active
	analytics.InactiveProducts = statusCounts.Inactive

	// Get low stock and out of stock products count
	var stockCounts struct {
		LowStock   int64 `db:"low_stock"`
		OutOfStock int64 `db:"out_of_stock"`
	}
	err = r.db.GetContext(ctx, &stockCounts, `
		SELECT 
			COUNT(CASE WHEN total_qty <= 10 AND total_qty > 0 THEN 1 END) as low_stock,
			COUNT(CASE WHEN total_qty = 0 THEN 1 END) as out_of_stock
		FROM (
			SELECT 
				p.id,
				COALESCE(SUM(i.available_quantity), 0) as total_qty
			FROM products p
			LEFT JOIN inventory i ON p.id = i.product_id
			GROUP BY p.id
		) as stock_stats`)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock products count: %w", err)
	}
	analytics.LowStockProducts = stockCounts.LowStock
	analytics.OutOfStockProducts = stockCounts.OutOfStock

	// Get total categories count
	err = r.db.GetContext(ctx, &analytics.TotalCategories, `SELECT COUNT(*) FROM categories`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total categories count: %w", err)
	}

	// Get average price
	err = r.db.GetContext(ctx, &analytics.AveragePrice, `SELECT COALESCE(AVG(price), 0) FROM products WHERE is_active = true`)
	if err != nil {
		return nil, fmt.Errorf("failed to get average price: %w", err)
	}

	// Get total inventory value
	err = r.db.GetContext(ctx, &analytics.TotalInventoryValue, `
		SELECT COALESCE(SUM(COALESCE(pv.price, p.price) * i.available_quantity), 0)
		FROM inventory i
		INNER JOIN products p ON i.product_id = p.id
		LEFT JOIN product_variants pv ON i.product_variant_id = pv.id
		WHERE p.is_active = true`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total inventory value: %w", err)
	}

	return analytics, nil
}

// GetTopSellingProducts retrieves top selling products
func (r *productRepository) GetTopSellingProducts(ctx context.Context, limit int) ([]*domain.ProductOrderStats, error) {
	query := `
		SELECT 
			p.id as product_id,
			p.name as product_name,
			p.sku,
			COALESCE(SUM(oi.quantity), 0) as total_quantity,
			COALESCE(SUM(oi.total_price), 0) as total_revenue,
			COUNT(DISTINCT o.id) as order_count,
			COALESCE(AVG(oi.unit_price), 0) as average_price
		FROM products p
		LEFT JOIN order_items oi ON p.id = oi.product_id
		LEFT JOIN orders o ON oi.order_id = o.id AND o.status IN ('confirmed', 'processing', 'shipped', 'delivered')
		WHERE p.is_active = true
		GROUP BY p.id, p.name, p.sku
		ORDER BY total_revenue DESC, total_quantity DESC
		LIMIT $1`

	var products []*domain.ProductOrderStats
	err := r.db.SelectContext(ctx, &products, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top selling products: %w", err)
	}

	return products, nil
}
