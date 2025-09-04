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

	AlertType         string     `json:"alert_type" db:"alert_type"` // low_stock, out_of_stock, reorder_point
	CurrentQuantity   int        `json:"current_quantity" db:"current_quantity"`
	ThresholdQuantity int        `json:"threshold_quantity" db:"threshold_quantity"`
	IsResolved        bool       `json:"is_resolved" db:"is_resolved"`
	ResolvedAt        *time.Time `json:"resolved_at" db:"resolved_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

// InventorySummary represents inventory summary statistics
type InventorySummary struct {
	TotalProducts     int64   `json:"total_products"`
	TotalVariants     int64   `json:"total_variants"`
	TotalQuantity     int64   `json:"total_quantity"`
	TotalReserved     int64   `json:"total_reserved"`
	TotalAvailable    int64   `json:"total_available"`
	LowStockItems     int64   `json:"low_stock_items"`
	OutOfStockItems   int64   `json:"out_of_stock_items"`
	TotalValue        float64 `json:"total_value"`
	AverageStockLevel float64 `json:"average_stock_level"`
}
