package dto

// Inventory Management DTOs

// CreateInventoryRequest represents the request to create inventory
type CreateInventoryRequest struct {
	ProductID        int64  `json:"product_id" validate:"required"`
	ProductVariantID *int64 `json:"product_variant_id" validate:"omitempty"`
	Quantity         int    `json:"quantity" validate:"required,min=0"`
	MinStockLevel    int    `json:"min_stock_level" validate:"min=0"`
	MaxStockLevel    *int   `json:"max_stock_level" validate:"omitempty,min=0"`
	ReorderPoint     int    `json:"reorder_point" validate:"min=0"`
}

// UpdateInventoryRequest represents the request to update inventory
type UpdateInventoryRequest struct {
	Quantity      *int `json:"quantity" validate:"omitempty,min=0"`
	MinStockLevel *int `json:"min_stock_level" validate:"omitempty,min=0"`
	MaxStockLevel *int `json:"max_stock_level" validate:"omitempty,min=0"`
	ReorderPoint  *int `json:"reorder_point" validate:"omitempty,min=0"`
}

// InventoryResponse represents the response for inventory data
type InventoryResponse struct {
	ID                int64  `json:"id"`
	ProductID         int64  `json:"product_id"`
	ProductVariantID  *int64 `json:"product_variant_id"`
	Quantity          int    `json:"quantity"`
	ReservedQuantity  int    `json:"reserved_quantity"`
	AvailableQuantity int    `json:"available_quantity"`
	MinStockLevel     int    `json:"min_stock_level"`
	MaxStockLevel     *int   `json:"max_stock_level"`
	ReorderPoint      int    `json:"reorder_point"`
	LastRestocked     string `json:"last_restocked"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// StockMovementRequest represents the request to record stock movement
type StockMovementRequest struct {
	ProductID        int64   `json:"product_id" validate:"required"`
	ProductVariantID *int64  `json:"product_variant_id" validate:"omitempty"`
	MovementType     string  `json:"movement_type" validate:"required,oneof=in out adjustment transfer"`
	Quantity         int     `json:"quantity" validate:"required"`
	Reference        *string `json:"reference" validate:"omitempty,max=255"`
	ReferenceType    *string `json:"reference_type" validate:"omitempty,max=50"`
	Reason           *string `json:"reason" validate:"omitempty,max=255"`
	Notes            *string `json:"notes" validate:"omitempty,max=1000"`
	CreatedBy        *int64  `json:"created_by" validate:"omitempty"`
}

// StockMovementResponse represents the response for stock movement data
type StockMovementResponse struct {
	ID               int64  `json:"id"`
	ProductID        int64  `json:"product_id"`
	ProductVariantID *int64 `json:"product_variant_id"`
	MovementType     string `json:"movement_type"`
	Quantity         int    `json:"quantity"`
	PreviousQuantity int    `json:"previous_quantity"`
	NewQuantity      int    `json:"new_quantity"`
	Reference        string `json:"reference"`
	ReferenceType    string `json:"reference_type"`
	Reason           string `json:"reason"`
	Notes            string `json:"notes"`
	CreatedBy        *int64 `json:"created_by"`
	CreatedAt        string `json:"created_at"`
}

// ReserveStockRequest represents the request to reserve stock
type ReserveStockRequest struct {
	ProductID        int64  `json:"product_id" validate:"required"`
	ProductVariantID *int64 `json:"product_variant_id" validate:"omitempty"`
	OrderID          int64  `json:"order_id" validate:"required"`
	Quantity         int    `json:"quantity" validate:"required,min=1"`
	ExpiresAt        string `json:"expires_at" validate:"required"`
}

// ReserveStockResponse represents the response for stock reservation
type ReserveStockResponse struct {
	ID               int64  `json:"id"`
	ProductID        int64  `json:"product_id"`
	ProductVariantID *int64 `json:"product_variant_id"`
	OrderID          int64  `json:"order_id"`
	Quantity         int    `json:"quantity"`
	ExpiresAt        string `json:"expires_at"`
	CreatedAt        string `json:"created_at"`
}

// ReleaseStockRequest represents the request to release reserved stock
type ReleaseStockRequest struct {
	ReservationID *int64 `json:"reservation_id" validate:"omitempty"`
	OrderID       *int64 `json:"order_id" validate:"omitempty"`
}

// InventoryAlertResponse represents the response for inventory alerts
type InventoryAlertResponse struct {
	ID                int64   `json:"id"`
	ProductID         int64   `json:"product_id"`
	ProductVariantID  *int64  `json:"product_variant_id"`
	AlertType         string  `json:"alert_type"`
	CurrentQuantity   int     `json:"current_quantity"`
	ThresholdQuantity int     `json:"threshold_quantity"`
	IsResolved        bool    `json:"is_resolved"`
	ResolvedAt        *string `json:"resolved_at"`
	CreatedAt         string  `json:"created_at"`
}

// ListInventoryRequest represents the request to list inventory
type ListInventoryRequest struct {
	ProductID        *int64 `json:"product_id" validate:"omitempty"`
	ProductVariantID *int64 `json:"product_variant_id" validate:"omitempty"`
	LowStock         *bool  `json:"low_stock" validate:"omitempty"`
	OutOfStock       *bool  `json:"out_of_stock" validate:"omitempty"`
	Page             int    `json:"page" validate:"min=1"`
	Limit            int    `json:"limit" validate:"min=1,max=100"`
}

// ListInventoryResponse represents the response for listing inventory
type ListInventoryResponse struct {
	Inventory  []InventoryResponse `json:"inventory"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

// ListStockMovementsRequest represents the request to list stock movements
type ListStockMovementsRequest struct {
	ProductID        *int64  `json:"product_id" validate:"omitempty"`
	ProductVariantID *int64  `json:"product_variant_id" validate:"omitempty"`
	MovementType     *string `json:"movement_type" validate:"omitempty,oneof=in out adjustment transfer"`
	StartDate        *string `json:"start_date" validate:"omitempty"`
	EndDate          *string `json:"end_date" validate:"omitempty"`
	Page             int     `json:"page" validate:"min=1"`
	Limit            int     `json:"limit" validate:"min=1,max=100"`
}

// ListStockMovementsResponse represents the response for listing stock movements
type ListStockMovementsResponse struct {
	Movements  []StockMovementResponse `json:"movements"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

// InventorySummaryResponse represents the response for inventory summary
type InventorySummaryResponse struct {
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

// BulkStockUpdateRequest represents the request for bulk stock updates
type BulkStockUpdateRequest struct {
	Updates []StockUpdateItem `json:"updates" validate:"required,min=1"`
}

// StockUpdateItem represents a single stock update item
type StockUpdateItem struct {
	ProductID        int64  `json:"product_id" validate:"required"`
	ProductVariantID *int64 `json:"product_variant_id" validate:"omitempty"`
	Quantity         int    `json:"quantity" validate:"required"`
	MovementType     string `json:"movement_type" validate:"required,oneof=in out adjustment"`
	Reason           string `json:"reason" validate:"required,max=255"`
	Notes            string `json:"notes" validate:"omitempty,max=1000"`
}

// BulkStockUpdateResponse represents the response for bulk stock updates
type BulkStockUpdateResponse struct {
	UpdatedItems int64                   `json:"updated_items"`
	FailedItems  []FailedStockUpdateItem `json:"failed_items"`
	Success      bool                    `json:"success"`
}

// FailedStockUpdateItem represents a failed stock update item
type FailedStockUpdateItem struct {
	ProductID        int64  `json:"product_id"`
	ProductVariantID *int64 `json:"product_variant_id"`
	Error            string `json:"error"`
}
