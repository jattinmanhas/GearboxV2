package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type InventoryRepository interface {
	// Inventory Management
	CreateInventory(ctx context.Context, inventory *domain.Inventory) error
	GetInventoryByID(ctx context.Context, id int64) (*domain.Inventory, error)
	GetInventoryByProduct(ctx context.Context, productID int64, variantID *int64) (*domain.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *domain.Inventory) error
	DeleteInventory(ctx context.Context, id int64) error
	ListInventory(ctx context.Context, req *ListInventoryRequest) ([]*domain.Inventory, int64, error)
	GetInventorySummary(ctx context.Context) (*domain.InventorySummary, error)

	// Stock Movements
	RecordStockMovement(ctx context.Context, movement *domain.InventoryMovement) error
	GetStockMovements(ctx context.Context, req *ListStockMovementsRequest) ([]*domain.InventoryMovement, int64, error)
	GetStockMovementByID(ctx context.Context, id int64) (*domain.InventoryMovement, error)

	// Stock Reservations
	ReserveStock(ctx context.Context, reservation *domain.StockReservation) error
	ReleaseStock(ctx context.Context, reservationID int64) error
	ReleaseStockByOrderID(ctx context.Context, orderID int64) error
	GetStockReservations(ctx context.Context, orderID int64) ([]*domain.StockReservation, error)
	GetExpiredReservations(ctx context.Context, before time.Time) ([]*domain.StockReservation, error)
	CleanupExpiredReservations(ctx context.Context) error

	// Inventory Alerts
	CreateInventoryAlert(ctx context.Context, alert *domain.InventoryAlert) error
	GetInventoryAlerts(ctx context.Context, resolved *bool) ([]*domain.InventoryAlert, error)
	ResolveInventoryAlert(ctx context.Context, alertID int64) error
	CheckLowStockAlerts(ctx context.Context) error

	// Bulk Operations
	BulkUpdateStock(ctx context.Context, updates []StockUpdateItem) (*BulkStockUpdateResponse, error)
}

type inventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Inventory Management

// CreateInventory creates a new inventory record
func (r *inventoryRepository) CreateInventory(ctx context.Context, inventory *domain.Inventory) error {
	query := `
		INSERT INTO inventory (
			product_id, product_variant_id, quantity, reserved_quantity, available_quantity,
			min_stock_level, max_stock_level, reorder_point, last_restocked, created_at, updated_at
		) VALUES (
			:product_id, :product_variant_id, :quantity, :reserved_quantity, :available_quantity,
			:min_stock_level, :max_stock_level, :reorder_point, :last_restocked, :created_at, :updated_at
		) RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, inventory)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&inventory.ID); err != nil {
			return fmt.Errorf("failed to get inventory ID: %w", err)
		}
	}

	return nil
}

// GetInventoryByID retrieves inventory by ID
func (r *inventoryRepository) GetInventoryByID(ctx context.Context, id int64) (*domain.Inventory, error) {
	var inventory domain.Inventory
	query := `
		SELECT id, product_id, product_variant_id, quantity, reserved_quantity, available_quantity,
			   min_stock_level, max_stock_level, reorder_point, last_restocked, created_at, updated_at
		FROM inventory WHERE id = $1`

	err := r.db.GetContext(ctx, &inventory, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("inventory with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return &inventory, nil
}

// GetInventoryByProduct retrieves inventory by product and variant
func (r *inventoryRepository) GetInventoryByProduct(ctx context.Context, productID int64, variantID *int64) (*domain.Inventory, error) {
	var inventory domain.Inventory
	var query string
	var args []interface{}

	if variantID != nil {
		query = `
			SELECT id, product_id, product_variant_id, quantity, reserved_quantity, available_quantity,
				   min_stock_level, max_stock_level, reorder_point, last_restocked, created_at, updated_at
			FROM inventory WHERE product_id = $1 AND product_variant_id = $2`
		args = []interface{}{productID, *variantID}
	} else {
		query = `
			SELECT id, product_id, product_variant_id, quantity, reserved_quantity, available_quantity,
				   min_stock_level, max_stock_level, reorder_point, last_restocked, created_at, updated_at
			FROM inventory WHERE product_id = $1 AND product_variant_id IS NULL`
		args = []interface{}{productID}
	}

	err := r.db.GetContext(ctx, &inventory, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return &inventory, nil
}

// UpdateInventory updates an existing inventory record
func (r *inventoryRepository) UpdateInventory(ctx context.Context, inventory *domain.Inventory) error {
	query := `
		UPDATE inventory SET
			quantity = :quantity, reserved_quantity = :reserved_quantity, available_quantity = :available_quantity,
			min_stock_level = :min_stock_level, max_stock_level = :max_stock_level, reorder_point = :reorder_point,
			last_restocked = :last_restocked, updated_at = :updated_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, inventory)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory with ID %d not found", inventory.ID)
	}

	return nil
}

// DeleteInventory deletes an inventory record
func (r *inventoryRepository) DeleteInventory(ctx context.Context, id int64) error {
	query := `DELETE FROM inventory WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory with ID %d not found", id)
	}

	return nil
}

// ListInventory lists inventory with filters
func (r *inventoryRepository) ListInventory(ctx context.Context, req *ListInventoryRequest) ([]*domain.Inventory, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if req.ProductID != nil {
		whereClause += fmt.Sprintf(" AND product_id = $%d", argIndex)
		args = append(args, *req.ProductID)
		argIndex++
	}

	if req.ProductVariantID != nil {
		whereClause += fmt.Sprintf(" AND product_variant_id = $%d", argIndex)
		args = append(args, *req.ProductVariantID)
		argIndex++
	}

	if req.LowStock != nil && *req.LowStock {
		whereClause += " AND available_quantity <= reorder_point"
	}

	if req.OutOfStock != nil && *req.OutOfStock {
		whereClause += " AND available_quantity = 0"
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM inventory %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, product_id, product_variant_id, quantity, reserved_quantity, available_quantity,
			   min_stock_level, max_stock_level, reorder_point, last_restocked, created_at, updated_at
		FROM inventory %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Limit, offset)

	var inventory []*domain.Inventory
	err = r.db.SelectContext(ctx, &inventory, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list inventory: %w", err)
	}

	return inventory, total, nil
}

// GetInventorySummary gets inventory summary statistics
func (r *inventoryRepository) GetInventorySummary(ctx context.Context) (*domain.InventorySummary, error) {
	query := `
		SELECT 
			COUNT(DISTINCT product_id) as total_products,
			COUNT(DISTINCT CASE WHEN product_variant_id IS NOT NULL THEN product_variant_id END) as total_variants,
			COALESCE(SUM(quantity), 0) as total_quantity,
			COALESCE(SUM(reserved_quantity), 0) as total_reserved,
			COALESCE(SUM(available_quantity), 0) as total_available,
			COUNT(CASE WHEN available_quantity <= reorder_point AND available_quantity > 0 THEN 1 END) as low_stock_items,
			COUNT(CASE WHEN available_quantity = 0 THEN 1 END) as out_of_stock_items
		FROM inventory`

	var summary domain.InventorySummary
	err := r.db.GetContext(ctx, &summary, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory summary: %w", err)
	}

	// Calculate average stock level
	if summary.TotalProducts > 0 {
		summary.AverageStockLevel = float64(summary.TotalQuantity) / float64(summary.TotalProducts)
	}

	return &summary, nil
}

// Stock Movements

// RecordStockMovement records a stock movement
func (r *inventoryRepository) RecordStockMovement(ctx context.Context, movement *domain.InventoryMovement) error {
	query := `
		INSERT INTO inventory_movements (
			product_id, product_variant_id, movement_type, quantity, previous_quantity, new_quantity,
			reference, reference_type, reason, notes, created_by, created_at
		) VALUES (
			:product_id, :product_variant_id, :movement_type, :quantity, :previous_quantity, :new_quantity,
			:reference, :reference_type, :reason, :notes, :created_by, :created_at
		) RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, movement)
	if err != nil {
		return fmt.Errorf("failed to record stock movement: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&movement.ID); err != nil {
			return fmt.Errorf("failed to get movement ID: %w", err)
		}
	}

	return nil
}

// GetStockMovements retrieves stock movements with filters
func (r *inventoryRepository) GetStockMovements(ctx context.Context, req *ListStockMovementsRequest) ([]*domain.InventoryMovement, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if req.ProductID != nil {
		whereClause += fmt.Sprintf(" AND product_id = $%d", argIndex)
		args = append(args, *req.ProductID)
		argIndex++
	}

	if req.ProductVariantID != nil {
		whereClause += fmt.Sprintf(" AND product_variant_id = $%d", argIndex)
		args = append(args, *req.ProductVariantID)
		argIndex++
	}

	if req.MovementType != nil {
		whereClause += fmt.Sprintf(" AND movement_type = $%d", argIndex)
		args = append(args, *req.MovementType)
		argIndex++
	}

	if req.StartDate != nil {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *req.StartDate)
		argIndex++
	}

	if req.EndDate != nil {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *req.EndDate)
		argIndex++
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM inventory_movements %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stock movements: %w", err)
	}

	// Get paginated results
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, product_id, product_variant_id, movement_type, quantity, previous_quantity, new_quantity,
			   reference, reference_type, reason, notes, created_by, created_at
		FROM inventory_movements %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, req.Limit, offset)

	var movements []*domain.InventoryMovement
	err = r.db.SelectContext(ctx, &movements, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list stock movements: %w", err)
	}

	return movements, total, nil
}

// GetStockMovementByID retrieves a stock movement by ID
func (r *inventoryRepository) GetStockMovementByID(ctx context.Context, id int64) (*domain.InventoryMovement, error) {
	var movement domain.InventoryMovement
	query := `
		SELECT id, product_id, product_variant_id, movement_type, quantity, previous_quantity, new_quantity,
			   reference, reference_type, reason, notes, created_by, created_at
		FROM inventory_movements WHERE id = $1`

	err := r.db.GetContext(ctx, &movement, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stock movement with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return &movement, nil
}

// Stock Reservations

// ReserveStock reserves stock for an order
func (r *inventoryRepository) ReserveStock(ctx context.Context, reservation *domain.StockReservation) error {
	query := `
		INSERT INTO stock_reservations (product_id, product_variant_id, order_id, quantity, expires_at, created_at)
		VALUES (:product_id, :product_variant_id, :order_id, :quantity, :expires_at, :created_at)
		RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, reservation)
	if err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&reservation.ID); err != nil {
			return fmt.Errorf("failed to get reservation ID: %w", err)
		}
	}

	return nil
}

// ReleaseStock releases reserved stock by reservation ID
func (r *inventoryRepository) ReleaseStock(ctx context.Context, reservationID int64) error {
	query := `DELETE FROM stock_reservations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, reservationID)
	if err != nil {
		return fmt.Errorf("failed to release stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stock reservation with ID %d not found", reservationID)
	}

	return nil
}

// ReleaseStockByOrderID releases all reserved stock for an order
func (r *inventoryRepository) ReleaseStockByOrderID(ctx context.Context, orderID int64) error {
	query := `DELETE FROM stock_reservations WHERE order_id = $1`

	_, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to release stock by order ID: %w", err)
	}

	return nil
}

// GetStockReservations gets stock reservations for an order
func (r *inventoryRepository) GetStockReservations(ctx context.Context, orderID int64) ([]*domain.StockReservation, error) {
	var reservations []*domain.StockReservation
	query := `
		SELECT id, product_id, product_variant_id, order_id, quantity, expires_at, created_at
		FROM stock_reservations WHERE order_id = $1
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &reservations, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock reservations: %w", err)
	}

	return reservations, nil
}

// GetExpiredReservations gets expired stock reservations
func (r *inventoryRepository) GetExpiredReservations(ctx context.Context, before time.Time) ([]*domain.StockReservation, error) {
	var reservations []*domain.StockReservation
	query := `
		SELECT id, product_id, product_variant_id, order_id, quantity, expires_at, created_at
		FROM stock_reservations WHERE expires_at < $1
		ORDER BY expires_at ASC`

	err := r.db.SelectContext(ctx, &reservations, query, before)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired reservations: %w", err)
	}

	return reservations, nil
}

// CleanupExpiredReservations removes expired stock reservations
func (r *inventoryRepository) CleanupExpiredReservations(ctx context.Context) error {
	query := `DELETE FROM stock_reservations WHERE expires_at < NOW()`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired reservations: %w", err)
	}

	return nil
}

// Inventory Alerts

// CreateInventoryAlert creates an inventory alert
func (r *inventoryRepository) CreateInventoryAlert(ctx context.Context, alert *domain.InventoryAlert) error {
	query := `
		INSERT INTO inventory_alerts (
			product_id, product_variant_id, alert_type, current_quantity, threshold_quantity,
			is_resolved, resolved_at, created_at
		) VALUES (
			:product_id, :product_variant_id, :alert_type, :current_quantity, :threshold_quantity,
			:is_resolved, :resolved_at, :created_at
		) RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, alert)
	if err != nil {
		return fmt.Errorf("failed to create inventory alert: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&alert.ID); err != nil {
			return fmt.Errorf("failed to get alert ID: %w", err)
		}
	}

	return nil
}

// GetInventoryAlerts gets inventory alerts
func (r *inventoryRepository) GetInventoryAlerts(ctx context.Context, resolved *bool) ([]*domain.InventoryAlert, error) {
	var alerts []*domain.InventoryAlert
	var query string
	var args []interface{}

	if resolved != nil {
		query = `
			SELECT id, product_id, product_variant_id, alert_type, current_quantity, threshold_quantity,
				   is_resolved, resolved_at, created_at
			FROM inventory_alerts WHERE is_resolved = $1
			ORDER BY created_at DESC`
		args = []interface{}{*resolved}
	} else {
		query = `
			SELECT id, product_id, product_variant_id, alert_type, current_quantity, threshold_quantity,
				   is_resolved, resolved_at, created_at
			FROM inventory_alerts
			ORDER BY created_at DESC`
		args = []interface{}{}
	}

	err := r.db.SelectContext(ctx, &alerts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory alerts: %w", err)
	}

	return alerts, nil
}

// ResolveInventoryAlert resolves an inventory alert
func (r *inventoryRepository) ResolveInventoryAlert(ctx context.Context, alertID int64) error {
	now := time.Now()
	query := `
		UPDATE inventory_alerts 
		SET is_resolved = true, resolved_at = $1, updated_at = $1
		WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, now, alertID)
	if err != nil {
		return fmt.Errorf("failed to resolve inventory alert: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory alert with ID %d not found", alertID)
	}

	return nil
}

// CheckLowStockAlerts checks for low stock and creates alerts
func (r *inventoryRepository) CheckLowStockAlerts(ctx context.Context) error {
	// This would typically be called by a background job
	// For now, we'll implement a simple check
	query := `
		INSERT INTO inventory_alerts (product_id, product_variant_id, alert_type, current_quantity, threshold_quantity, is_resolved, created_at)
		SELECT 
			i.product_id, i.product_variant_id, 'low_stock', i.available_quantity, i.reorder_point, false, NOW()
		FROM inventory i
		WHERE i.available_quantity <= i.reorder_point 
		AND i.available_quantity > 0
		AND NOT EXISTS (
			SELECT 1 FROM inventory_alerts ia 
			WHERE ia.product_id = i.product_id 
			AND ia.product_variant_id = i.product_variant_id 
			AND ia.alert_type = 'low_stock' 
			AND ia.is_resolved = false
		)`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to check low stock alerts: %w", err)
	}

	return nil
}

// Bulk Operations

// BulkUpdateStock performs bulk stock updates
func (r *inventoryRepository) BulkUpdateStock(ctx context.Context, updates []StockUpdateItem) (*BulkStockUpdateResponse, error) {
	// This would typically use a transaction for atomicity
	// For now, we'll implement a simple version
	response := &BulkStockUpdateResponse{
		UpdatedItems: 0,
		FailedItems:  []FailedStockUpdateItem{},
		Success:      true,
	}

	// TODO: Implement bulk update logic with proper transaction handling
	// This is a placeholder implementation

	return response, nil
}

// Additional types for repository
type ListInventoryRequest struct {
	ProductID        *int64
	ProductVariantID *int64
	LowStock         *bool
	OutOfStock       *bool
	Page             int
	Limit            int
}

type ListStockMovementsRequest struct {
	ProductID        *int64
	ProductVariantID *int64
	MovementType     *string
	StartDate        *string
	EndDate          *string
	Page             int
	Limit            int
}

type StockUpdateItem struct {
	ProductID        int64
	ProductVariantID *int64
	Quantity         int
	MovementType     string
	Reason           string
	Notes            string
}

type BulkStockUpdateResponse struct {
	UpdatedItems int64
	FailedItems  []FailedStockUpdateItem
	Success      bool
}

type FailedStockUpdateItem struct {
	ProductID        int64
	ProductVariantID *int64
	Error            string
}
