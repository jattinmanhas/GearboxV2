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
	UpdateProductQuantity(ctx context.Context, id int64, quantity int) error
	GetProductsByTags(ctx context.Context, tags []string, offset, limit int) ([]*domain.Product, int64, error)

	// Product Variants
	CreateProductVariant(ctx context.Context, variant *domain.ProductVariant) error
	GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error)
	GetProductVariantsByProductID(ctx context.Context, productID int64) ([]*domain.ProductVariant, error)
	UpdateProductVariant(ctx context.Context, id int64, variant *domain.ProductVariant) error
	DeleteProductVariant(ctx context.Context, id int64) error

	// Product Categories
	AddProductToCategory(ctx context.Context, productID, categoryID int64, isPrimary bool) error
	RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error
	GetProductCategories(ctx context.Context, productID int64) ([]*domain.Category, error)
	UpdateProductCategories(ctx context.Context, productID int64, categoryIDs []int64) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// CreateProduct creates a new product
func (r *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (
			name, description, short_description, sku, price, compare_price, cost_price,
			weight, dimensions, is_active, is_digital, requires_shipping, taxable,
			track_quantity, quantity, min_quantity, max_quantity, meta_title,
			meta_description, tags, created_at, updated_at
		) VALUES (
			:name, :description, :short_description, :sku, :price, :compare_price, :cost_price,
			:weight, :dimensions, :is_active, :is_digital, :requires_shipping, :taxable,
			:track_quantity, :quantity, :min_quantity, :max_quantity, :meta_title,
			:meta_description, :tags, :created_at, :updated_at
		)`

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get product ID: %w", err)
	}

	product.ID = id
	return nil
}

// GetProductByID retrieves a product by ID
func (r *productRepository) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT * FROM products WHERE id = ?`

	var product domain.Product
	err := r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// GetProductBySKU retrieves a product by SKU
func (r *productRepository) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	query := `SELECT * FROM products WHERE sku = ?`

	var product domain.Product
	err := r.db.GetContext(ctx, &product, query, sku)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with SKU %s not found", sku)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// UpdateProduct updates an existing product
func (r *productRepository) UpdateProduct(ctx context.Context, id int64, product *domain.Product) error {
	query := `
		UPDATE products SET
			name = :name, description = :description, short_description = :short_description,
			sku = :sku, price = :price, compare_price = :compare_price, cost_price = :cost_price,
			weight = :weight, dimensions = :dimensions, is_active = :is_active,
			is_digital = :is_digital, requires_shipping = :requires_shipping, taxable = :taxable,
			track_quantity = :track_quantity, quantity = :quantity, min_quantity = :min_quantity,
			max_quantity = :max_quantity, meta_title = :meta_title, meta_description = :meta_description,
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
	query := `DELETE FROM products WHERE id = ?`

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
		LIMIT ? OFFSET ?`, whereClause, orderClause)

	args = append(args, limit, offset)

	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, total, nil
}

// GetProductsByCategory retrieves products by category
func (r *productRepository) GetProductsByCategory(ctx context.Context, categoryID int64, offset, limit int) ([]*domain.Product, int64, error) {
	// Count query
	countQuery := `
		SELECT COUNT(DISTINCT p.id) 
		FROM products p 
		INNER JOIN product_categories pc ON p.id = pc.product_id 
		WHERE pc.category_id = ?`

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

	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
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
		WHERE name LIKE ? OR description LIKE ? OR sku LIKE ? OR tags LIKE ?`

	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Search query
	searchQuery := `
		SELECT * 
		FROM products 
		WHERE name LIKE ? OR description LIKE ? OR sku LIKE ? OR tags LIKE ?
		ORDER BY 
			CASE 
				WHEN name LIKE ? THEN 1
				WHEN sku LIKE ? THEN 2
				WHEN description LIKE ? THEN 3
				ELSE 4
			END,
			name
		LIMIT ? OFFSET ?`

	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, searchQuery,
		searchTerm, searchTerm, searchTerm, searchTerm,
		searchTerm, searchTerm, searchTerm,
		limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return products, total, nil
}

// UpdateProductQuantity updates product quantity
func (r *productRepository) UpdateProductQuantity(ctx context.Context, id int64, quantity int) error {
	query := `UPDATE products SET quantity = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, quantity, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
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

	args = append(args, limit, offset)

	var products []*domain.Product
	err = r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products by tags: %w", err)
	}

	return products, total, nil
}

// Product Variant methods

// CreateProductVariant creates a new product variant
func (r *productRepository) CreateProductVariant(ctx context.Context, variant *domain.ProductVariant) error {
	query := `
		INSERT INTO product_variants (
			product_id, name, sku, price, compare_price, cost_price,
			weight, quantity, is_active, position
		) VALUES (
			:product_id, :name, :sku, :price, :compare_price, :cost_price,
			:weight, :quantity, :is_active, :position
		)`

	result, err := r.db.NamedExecContext(ctx, query, variant)
	if err != nil {
		return fmt.Errorf("failed to create product variant: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get variant ID: %w", err)
	}

	variant.ID = id
	return nil
}

// GetProductVariantByID retrieves a product variant by ID
func (r *productRepository) GetProductVariantByID(ctx context.Context, id int64) (*domain.ProductVariant, error) {
	query := `SELECT * FROM product_variants WHERE id = ?`

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
	query := `SELECT * FROM product_variants WHERE product_id = ? ORDER BY position, name`

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
			cost_price = :cost_price, weight = :weight, quantity = :quantity,
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
	query := `DELETE FROM product_variants WHERE id = ?`

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
		VALUES (?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, productID, categoryID, isPrimary)
	if err != nil {
		return fmt.Errorf("failed to add product to category: %w", err)
	}

	return nil
}

// RemoveProductFromCategory removes a product from a category
func (r *productRepository) RemoveProductFromCategory(ctx context.Context, productID, categoryID int64) error {
	query := `DELETE FROM product_categories WHERE product_id = ? AND category_id = ?`

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
		WHERE pc.product_id = ? 
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
	_, err = tx.ExecContext(ctx, "DELETE FROM product_categories WHERE product_id = ?", productID)
	if err != nil {
		return fmt.Errorf("failed to remove existing categories: %w", err)
	}

	// Add new categories
	for i, categoryID := range categoryIDs {
		isPrimary := i == 0 // First category is primary
		_, err = tx.ExecContext(ctx,
			"INSERT INTO product_categories (product_id, category_id, is_primary) VALUES (?, ?, ?)",
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

	if filter.CategoryID != nil {
		conditions = append(conditions, "id IN (SELECT product_id FROM product_categories WHERE category_id = ?)")
		args = append(args, *filter.CategoryID)
	}

	if filter.IsActive != nil {
		conditions = append(conditions, "is_active = ?")
		args = append(args, *filter.IsActive)
	}

	if filter.IsDigital != nil {
		conditions = append(conditions, "is_digital = ?")
		args = append(args, *filter.IsDigital)
	}

	if filter.MinPrice != nil {
		conditions = append(conditions, "price >= ?")
		args = append(args, *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, "price <= ?")
		args = append(args, *filter.MaxPrice)
	}

	if filter.InStock != nil && *filter.InStock {
		conditions = append(conditions, "quantity > 0")
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		conditions = append(conditions, "(name LIKE ? OR description LIKE ? OR sku LIKE ? OR tags LIKE ?)")
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if len(filter.Tags) > 0 {
		var tagConditions []string
		for _, tag := range filter.Tags {
			tagConditions = append(tagConditions, "tags LIKE ?")
			args = append(args, "%"+tag+"%")
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " OR ")+")")
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
