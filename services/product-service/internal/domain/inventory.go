package domain

import (
	"time"
)

// Inventory represents inventory tracking for products
type Inventory struct {
	ID                int64     `json:"id" db:"id"`
	ProductID         int64     `json:"product_id" db:"product_id"`
	ProductVariantID  *int64    `json:"product_variant_id" db:"product_variant_id"`
	Quantity          int       `json:"quantity" db:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity" db:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity" db:"available_quantity"`
	MinStockLevel     int       `json:"min_stock_level" db:"min_stock_level"`
	MaxStockLevel     int       `json:"max_stock_level" db:"max_stock_level"`
	ReorderPoint      int       `json:"reorder_point" db:"reorder_point"`
	LastRestocked     time.Time `json:"last_restocked" db:"last_restocked"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	// Product information for display
	ProductName string  `json:"product_name" db:"product_name"`
	ProductSKU  string  `json:"product_sku" db:"product_sku"`
	VariantName *string `json:"variant_name" db:"variant_name"`
	VariantSKU  *string `json:"variant_sku" db:"variant_sku"`
}

// InventoryMovement represents inventory movements (stock in/out)
type InventoryMovement struct {
	ID               int64     `json:"id" db:"id"`
	ProductID        int64     `json:"product_id" db:"product_id"`
	ProductVariantID *int64    `json:"product_variant_id" db:"product_variant_id"`
	MovementType     string    `json:"movement_type" db:"movement_type"` // in, out, adjustment, transfer
	Quantity         int       `json:"quantity" db:"quantity"`
	PreviousQuantity int       `json:"previous_quantity" db:"previous_quantity"`
	NewQuantity      int       `json:"new_quantity" db:"new_quantity"`
	Reference        string    `json:"reference" db:"reference"` // order_id, purchase_order_id, etc.
	ReferenceType    string    `json:"reference_type" db:"reference_type"`
	Reason           string    `json:"reason" db:"reason"`
	Notes            string    `json:"notes" db:"notes"`
	CreatedBy        int64     `json:"created_by" db:"created_by"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// StockReservation represents reserved stock for pending orders
type StockReservation struct {
	ID               int64     `json:"id" db:"id"`
	ProductID        int64     `json:"product_id" db:"product_id"`
	ProductVariantID *int64    `json:"product_variant_id" db:"product_variant_id"`
	OrderID          int64     `json:"order_id" db:"order_id"`
	Quantity         int       `json:"quantity" db:"quantity"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// InventoryAlert represents low stock alerts
type InventoryAlert struct {
	ID               int64  `json:"id" db:"id"`
	ProductID        int64  `json:"product_id" db:"product_id"`
	ProductVariantID *int64 `json:"product_variant_id" db:"product_variant_id"`

	AlertType         string     `json:"alert_type" db:"alert_type"`
	CurrentQuantity   int        `json:"current_quantity" db:"current_quantity"`
	ThresholdQuantity int        `json:"threshold_quantity" db:"threshold_quantity"`
	IsResolved        bool       `json:"is_resolved" db:"is_resolved"`
	ResolvedAt        *time.Time `json:"resolved_at" db:"resolved_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`

	// Product information
	ProductName string  `json:"product_name" db:"product_name"`
	ProductSKU  string  `json:"product_sku" db:"product_sku"`
	VariantName *string `json:"variant_name" db:"variant_name"`
	VariantSKU  *string `json:"variant_sku" db:"variant_sku"`
}

// InventorySummary represents inventory summary statistics
type InventorySummary struct {
	TotalProducts     int64   `db:"total_products" json:"total_products"`
	TotalVariants     int64   `db:"total_variants" json:"total_variants"`
	TotalQuantity     int64   `db:"total_quantity" json:"total_quantity"`
	TotalReserved     int64   `db:"total_reserved" json:"total_reserved"`
	TotalAvailable    int64   `db:"total_available" json:"total_available"`
	LowStockItems     int64   `db:"low_stock_items" json:"low_stock_items"`
	OutOfStockItems   int64   `db:"out_of_stock_items" json:"out_of_stock_items"`
	TotalValue        float64 `db:"total_value" json:"total_value"`
	AverageStockLevel float64 `db:"average_stock_level" json:"average_stock_level"`
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
