package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type InventoryService interface {
	// Inventory Management
	CreateInventory(ctx context.Context, req *dto.CreateInventoryRequest) (*domain.Inventory, error)
	GetInventoryByID(ctx context.Context, id int64) (*domain.Inventory, error)
	GetInventoryByProduct(ctx context.Context, productID int64, variantID *int64) (*domain.Inventory, error)
	UpdateInventory(ctx context.Context, id int64, req *dto.UpdateInventoryRequest) (*domain.Inventory, error)
	DeleteInventory(ctx context.Context, id int64) error
	ListInventory(ctx context.Context, req *dto.ListInventoryRequest) (*dto.ListInventoryResponse, error)
	GetInventorySummary(ctx context.Context) (*dto.InventorySummaryResponse, error)

	// Stock Movements
	RecordStockMovement(ctx context.Context, req *dto.StockMovementRequest) (*domain.InventoryMovement, error)
	GetStockMovements(ctx context.Context, req *dto.ListStockMovementsRequest) (*dto.ListStockMovementsResponse, error)
	GetStockMovementByID(ctx context.Context, id int64) (*domain.InventoryMovement, error)

	// Stock Reservations
	ReserveStock(ctx context.Context, req *dto.ReserveStockRequest) (*domain.StockReservation, error)
	ReleaseStock(ctx context.Context, req *dto.ReleaseStockRequest) error
	GetStockReservations(ctx context.Context, orderID int64) ([]*domain.StockReservation, error)

	// Inventory Alerts
	GetInventoryAlerts(ctx context.Context, resolved *bool) ([]*domain.InventoryAlert, error)
	ResolveInventoryAlert(ctx context.Context, alertID int64) error
	CheckLowStockAlerts(ctx context.Context) error

	// Bulk Operations
	BulkUpdateStock(ctx context.Context, req *dto.BulkStockUpdateRequest) (*dto.BulkStockUpdateResponse, error)
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	productRepo   repository.ProductRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository, productRepo repository.ProductRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

// Inventory Management

// CreateInventory creates a new inventory record
func (s *inventoryService) CreateInventory(ctx context.Context, req *dto.CreateInventoryRequest) (*domain.Inventory, error) {
	// Validate product exists
	_, err := s.productRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate product: %w", err)
	}

	// Validate variant exists if provided
	if req.ProductVariantID != nil {
		_, err := s.productRepo.GetProductVariantByID(ctx, *req.ProductVariantID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate product variant: %w", err)
		}
	}

	// Check if inventory already exists
	_, err = s.inventoryRepo.GetInventoryByProduct(ctx, req.ProductID, req.ProductVariantID)
	if err == nil {
		return nil, fmt.Errorf("inventory already exists for this product/variant combination")
	}

	now := time.Now()
	inventory := &domain.Inventory{
		ProductID:         req.ProductID,
		ProductVariantID:  req.ProductVariantID,
		Quantity:          req.Quantity,
		ReservedQuantity:  0,
		AvailableQuantity: req.Quantity,
		MinStockLevel:     req.MinStockLevel,
		MaxStockLevel:     getIntValue(req.MaxStockLevel),
		ReorderPoint:      req.ReorderPoint,
		LastRestocked:     now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err = s.inventoryRepo.CreateInventory(ctx, inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	// Record initial stock movement
	movement := &domain.InventoryMovement{
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		MovementType:     "in",
		Quantity:         req.Quantity,
		PreviousQuantity: 0,
		NewQuantity:      req.Quantity,
		Reference:        "initial_stock",
		ReferenceType:    "setup",
		Reason:           "Initial inventory setup",
		CreatedAt:        now,
	}

	err = s.inventoryRepo.RecordStockMovement(ctx, movement)
	if err != nil {
		// Log error but don't fail the inventory creation
		fmt.Printf("Warning: failed to record initial stock movement: %v\n", err)
	}

	return inventory, nil
}

// GetInventoryByID retrieves inventory by ID
func (s *inventoryService) GetInventoryByID(ctx context.Context, id int64) (*domain.Inventory, error) {
	inventory, err := s.inventoryRepo.GetInventoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return inventory, nil
}

// GetInventoryByProduct retrieves inventory by product and variant
func (s *inventoryService) GetInventoryByProduct(ctx context.Context, productID int64, variantID *int64) (*domain.Inventory, error) {
	inventory, err := s.inventoryRepo.GetInventoryByProduct(ctx, productID, variantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return inventory, nil
}

// UpdateInventory updates an existing inventory record
func (s *inventoryService) UpdateInventory(ctx context.Context, id int64, req *dto.UpdateInventoryRequest) (*domain.Inventory, error) {
	// Get existing inventory
	existing, err := s.inventoryRepo.GetInventoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing inventory: %w", err)
	}

	// Update fields that are provided
	updateInventory := *existing

	if req.Quantity != nil {
		updateInventory.Quantity = *req.Quantity
		updateInventory.AvailableQuantity = *req.Quantity - updateInventory.ReservedQuantity
	}
	if req.MinStockLevel != nil {
		updateInventory.MinStockLevel = *req.MinStockLevel
	}
	if req.MaxStockLevel != nil {
		updateInventory.MaxStockLevel = *req.MaxStockLevel
	}
	if req.ReorderPoint != nil {
		updateInventory.ReorderPoint = *req.ReorderPoint
	}

	updateInventory.UpdatedAt = time.Now()

	err = s.inventoryRepo.UpdateInventory(ctx, &updateInventory)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	return &updateInventory, nil
}

// DeleteInventory deletes an inventory record
func (s *inventoryService) DeleteInventory(ctx context.Context, id int64) error {
	// Check if inventory exists
	_, err := s.inventoryRepo.GetInventoryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	err = s.inventoryRepo.DeleteInventory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	return nil
}

// ListInventory lists inventory with filters
func (s *inventoryService) ListInventory(ctx context.Context, req *dto.ListInventoryRequest) (*dto.ListInventoryResponse, error) {
	// Convert DTO request to repository request
	repoReq := &repository.ListInventoryRequest{
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		LowStock:         req.LowStock,
		OutOfStock:       req.OutOfStock,
		Page:             req.Page,
		Limit:            req.Limit,
	}

	inventory, total, err := s.inventoryRepo.ListInventory(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}

	// Convert domain models to DTOs
	var inventoryResponses []dto.InventoryResponse
	for _, inv := range inventory {
		response := dto.InventoryResponse{
			ID:                inv.ID,
			ProductID:         inv.ProductID,
			ProductVariantID:  inv.ProductVariantID,
			Quantity:          inv.Quantity,
			ReservedQuantity:  inv.ReservedQuantity,
			AvailableQuantity: inv.AvailableQuantity,
			MinStockLevel:     inv.MinStockLevel,
			MaxStockLevel:     getIntPointer(inv.MaxStockLevel),
			ReorderPoint:      inv.ReorderPoint,
			LastRestocked:     inv.LastRestocked.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:         inv.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:         inv.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		inventoryResponses = append(inventoryResponses, response)
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.ListInventoryResponse{
		Inventory:  inventoryResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetInventorySummary gets inventory summary statistics
func (s *inventoryService) GetInventorySummary(ctx context.Context) (*dto.InventorySummaryResponse, error) {
	summary, err := s.inventoryRepo.GetInventorySummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory summary: %w", err)
	}

	return &dto.InventorySummaryResponse{
		TotalProducts:     summary.TotalProducts,
		TotalVariants:     summary.TotalVariants,
		TotalQuantity:     summary.TotalQuantity,
		TotalReserved:     summary.TotalReserved,
		TotalAvailable:    summary.TotalAvailable,
		LowStockItems:     summary.LowStockItems,
		OutOfStockItems:   summary.OutOfStockItems,
		TotalValue:        summary.TotalValue,
		AverageStockLevel: summary.AverageStockLevel,
	}, nil
}

// Stock Movements

// RecordStockMovement records a stock movement
func (s *inventoryService) RecordStockMovement(ctx context.Context, req *dto.StockMovementRequest) (*domain.InventoryMovement, error) {
	// Get current inventory
	inventory, err := s.inventoryRepo.GetInventoryByProduct(ctx, req.ProductID, req.ProductVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	// Calculate new quantity based on movement type
	var newQuantity int
	switch req.MovementType {
	case "in":
		newQuantity = inventory.Quantity + req.Quantity
	case "out":
		newQuantity = inventory.Quantity - req.Quantity
		if newQuantity < 0 {
			return nil, fmt.Errorf("insufficient stock: requested %d, available %d", req.Quantity, inventory.Quantity)
		}
	case "adjustment":
		newQuantity = req.Quantity
	case "transfer":
		newQuantity = inventory.Quantity - req.Quantity
		if newQuantity < 0 {
			return nil, fmt.Errorf("insufficient stock for transfer: requested %d, available %d", req.Quantity, inventory.Quantity)
		}
	default:
		return nil, fmt.Errorf("invalid movement type: %s", req.MovementType)
	}

	// Update inventory
	previousQuantity := inventory.Quantity
	inventory.Quantity = newQuantity
	inventory.AvailableQuantity = newQuantity - inventory.ReservedQuantity
	inventory.UpdatedAt = time.Now()

	err = s.inventoryRepo.UpdateInventory(ctx, inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Record movement
	movement := &domain.InventoryMovement{
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		MovementType:     req.MovementType,
		Quantity:         req.Quantity,
		PreviousQuantity: previousQuantity,
		NewQuantity:      newQuantity,
		Reference:        getStringValue(req.Reference),
		ReferenceType:    getStringValue(req.ReferenceType),
		Reason:           getStringValue(req.Reason),
		Notes:            getStringValue(req.Notes),
		CreatedBy:        getInt64Value(req.CreatedBy),
		CreatedAt:        time.Now(),
	}

	err = s.inventoryRepo.RecordStockMovement(ctx, movement)
	if err != nil {
		return nil, fmt.Errorf("failed to record stock movement: %w", err)
	}

	return movement, nil
}

// GetStockMovements retrieves stock movements with filters
func (s *inventoryService) GetStockMovements(ctx context.Context, req *dto.ListStockMovementsRequest) (*dto.ListStockMovementsResponse, error) {
	// Convert DTO request to repository request
	repoReq := &repository.ListStockMovementsRequest{
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		MovementType:     req.MovementType,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		Page:             req.Page,
		Limit:            req.Limit,
	}

	movements, total, err := s.inventoryRepo.GetStockMovements(ctx, repoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
	}

	// Convert domain models to DTOs
	var movementResponses []dto.StockMovementResponse
	for _, movement := range movements {
		response := dto.StockMovementResponse{
			ID:               movement.ID,
			ProductID:        movement.ProductID,
			ProductVariantID: movement.ProductVariantID,
			MovementType:     movement.MovementType,
			Quantity:         movement.Quantity,
			PreviousQuantity: movement.PreviousQuantity,
			NewQuantity:      movement.NewQuantity,
			Reference:        movement.Reference,
			ReferenceType:    movement.ReferenceType,
			Reason:           movement.Reason,
			Notes:            movement.Notes,
			CreatedBy:        getInt64Pointer(movement.CreatedBy),
			CreatedAt:        movement.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		movementResponses = append(movementResponses, response)
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.ListStockMovementsResponse{
		Movements:  movementResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetStockMovementByID retrieves a stock movement by ID
func (s *inventoryService) GetStockMovementByID(ctx context.Context, id int64) (*domain.InventoryMovement, error) {
	movement, err := s.inventoryRepo.GetStockMovementByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movement: %w", err)
	}

	return movement, nil
}

// Stock Reservations

// ReserveStock reserves stock for an order
func (s *inventoryService) ReserveStock(ctx context.Context, req *dto.ReserveStockRequest) (*domain.StockReservation, error) {
	// Parse expiration time
	expiresAt, err := time.Parse("2006-01-02T15:04:05Z07:00", req.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time format: %w", err)
	}

	// Get current inventory
	inventory, err := s.inventoryRepo.GetInventoryByProduct(ctx, req.ProductID, req.ProductVariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	// Check if enough stock is available
	if inventory.AvailableQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: requested %d, available %d", req.Quantity, inventory.AvailableQuantity)
	}

	// Create reservation
	reservation := &domain.StockReservation{
		ProductID:        req.ProductID,
		ProductVariantID: req.ProductVariantID,
		OrderID:          req.OrderID,
		Quantity:         req.Quantity,
		ExpiresAt:        expiresAt,
		CreatedAt:        time.Now(),
	}

	err = s.inventoryRepo.ReserveStock(ctx, reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Update inventory
	inventory.ReservedQuantity += req.Quantity
	inventory.AvailableQuantity -= req.Quantity
	inventory.UpdatedAt = time.Now()

	err = s.inventoryRepo.UpdateInventory(ctx, inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to update inventory after reservation: %w", err)
	}

	return reservation, nil
}

// ReleaseStock releases reserved stock
func (s *inventoryService) ReleaseStock(ctx context.Context, req *dto.ReleaseStockRequest) error {
	if req.ReservationID != nil {
		// Release by reservation ID
		err := s.inventoryRepo.ReleaseStock(ctx, *req.ReservationID)
		if err != nil {
			return fmt.Errorf("failed to release stock by reservation ID: %w", err)
		}
	} else if req.OrderID != nil {
		// Release by order ID
		err := s.inventoryRepo.ReleaseStockByOrderID(ctx, *req.OrderID)
		if err != nil {
			return fmt.Errorf("failed to release stock by order ID: %w", err)
		}
	} else {
		return fmt.Errorf("either reservation_id or order_id must be provided")
	}

	return nil
}

// GetStockReservations gets stock reservations for an order
func (s *inventoryService) GetStockReservations(ctx context.Context, orderID int64) ([]*domain.StockReservation, error) {
	reservations, err := s.inventoryRepo.GetStockReservations(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock reservations: %w", err)
	}

	return reservations, nil
}

// Inventory Alerts

// GetInventoryAlerts gets inventory alerts
func (s *inventoryService) GetInventoryAlerts(ctx context.Context, resolved *bool) ([]*domain.InventoryAlert, error) {
	alerts, err := s.inventoryRepo.GetInventoryAlerts(ctx, resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory alerts: %w", err)
	}

	return alerts, nil
}

// ResolveInventoryAlert resolves an inventory alert
func (s *inventoryService) ResolveInventoryAlert(ctx context.Context, alertID int64) error {
	err := s.inventoryRepo.ResolveInventoryAlert(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to resolve inventory alert: %w", err)
	}

	return nil
}

// CheckLowStockAlerts checks for low stock and creates alerts
func (s *inventoryService) CheckLowStockAlerts(ctx context.Context) error {
	err := s.inventoryRepo.CheckLowStockAlerts(ctx)
	if err != nil {
		return fmt.Errorf("failed to check low stock alerts: %w", err)
	}

	return nil
}

// Bulk Operations

// BulkUpdateStock performs bulk stock updates
func (s *inventoryService) BulkUpdateStock(ctx context.Context, req *dto.BulkStockUpdateRequest) (*dto.BulkStockUpdateResponse, error) {
	// Convert DTO items to repository items
	var repoUpdates []repository.StockUpdateItem
	for _, item := range req.Updates {
		repoUpdates = append(repoUpdates, repository.StockUpdateItem{
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			MovementType:     item.MovementType,
			Reason:           item.Reason,
			Notes:            item.Notes,
		})
	}

	response, err := s.inventoryRepo.BulkUpdateStock(ctx, repoUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to perform bulk stock update: %w", err)
	}

	// Convert repository response to DTO response
	dtoResponse := &dto.BulkStockUpdateResponse{
		UpdatedItems: response.UpdatedItems,
		FailedItems:  []dto.FailedStockUpdateItem{},
		Success:      response.Success,
	}

	for _, failed := range response.FailedItems {
		dtoResponse.FailedItems = append(dtoResponse.FailedItems, dto.FailedStockUpdateItem{
			ProductID:        failed.ProductID,
			ProductVariantID: failed.ProductVariantID,
			Error:            failed.Error,
		})
	}

	return dtoResponse, nil
}

// Helper functions to safely get values from pointers
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getIntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func getIntPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func getInt64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func getInt64Pointer(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}
