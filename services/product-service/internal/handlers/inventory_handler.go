package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

type IInventoryHandler interface {
	// Inventory Management
	CreateInventory(w http.ResponseWriter, r *http.Request)
	GetInventoryByID(w http.ResponseWriter, r *http.Request)
	GetInventoryByProduct(w http.ResponseWriter, r *http.Request)
	UpdateInventory(w http.ResponseWriter, r *http.Request)
	DeleteInventory(w http.ResponseWriter, r *http.Request)
	ListInventory(w http.ResponseWriter, r *http.Request)
	GetInventorySummary(w http.ResponseWriter, r *http.Request)

	// Stock Movements
	RecordStockMovement(w http.ResponseWriter, r *http.Request)
	GetStockMovements(w http.ResponseWriter, r *http.Request)
	GetStockMovementByID(w http.ResponseWriter, r *http.Request)

	// Stock Reservations
	ReserveStock(w http.ResponseWriter, r *http.Request)
	ReleaseStock(w http.ResponseWriter, r *http.Request)
	GetStockReservations(w http.ResponseWriter, r *http.Request)

	// Inventory Alerts
	GetInventoryAlerts(w http.ResponseWriter, r *http.Request)
	ResolveInventoryAlert(w http.ResponseWriter, r *http.Request)
	CheckLowStockAlerts(w http.ResponseWriter, r *http.Request)

	// Bulk Operations
	BulkUpdateStock(w http.ResponseWriter, r *http.Request)
}

type inventoryHandler struct {
	inventoryService services.InventoryService
}

func NewInventoryHandler(inventoryService services.InventoryService) IInventoryHandler {
	return &inventoryHandler{
		inventoryService: inventoryService,
	}
}

// Inventory Management

// CreateInventory creates a new inventory record
func (h *inventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	inventory, err := h.inventoryService.CreateInventory(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create inventory", err)
		return
	}

	response := dto.InventoryResponse{
		ID:                inventory.ID,
		ProductID:         inventory.ProductID,
		ProductVariantID:  inventory.ProductVariantID,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		MinStockLevel:     inventory.MinStockLevel,
		MaxStockLevel:     getIntPointer(inventory.MaxStockLevel),
		ReorderPoint:      inventory.ReorderPoint,
		LastRestocked:     inventory.LastRestocked.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:         inventory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Inventory created successfully", response)
}

// GetInventoryByID retrieves inventory by ID
func (h *inventoryHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid inventory ID", err)
		return
	}

	inventory, err := h.inventoryService.GetInventoryByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Inventory not found", err)
		return
	}

	response := dto.InventoryResponse{
		ID:                inventory.ID,
		ProductID:         inventory.ProductID,
		ProductVariantID:  inventory.ProductVariantID,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		MinStockLevel:     inventory.MinStockLevel,
		MaxStockLevel:     getIntPointer(inventory.MaxStockLevel),
		ReorderPoint:      inventory.ReorderPoint,
		LastRestocked:     inventory.LastRestocked.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:         inventory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Inventory retrieved successfully", response)
}

// GetInventoryByProduct retrieves inventory by product and variant
func (h *inventoryHandler) GetInventoryByProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	var variantID *int64
	variantIDStr := r.URL.Query().Get("variant_id")
	if variantIDStr != "" {
		variantIDVal, err := strconv.ParseInt(variantIDStr, 10, 64)
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "Invalid variant ID", err)
			return
		}
		variantID = &variantIDVal
	}

	inventory, err := h.inventoryService.GetInventoryByProduct(r.Context(), productID, variantID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Inventory not found", err)
		return
	}

	response := dto.InventoryResponse{
		ID:                inventory.ID,
		ProductID:         inventory.ProductID,
		ProductVariantID:  inventory.ProductVariantID,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		MinStockLevel:     inventory.MinStockLevel,
		MaxStockLevel:     getIntPointer(inventory.MaxStockLevel),
		ReorderPoint:      inventory.ReorderPoint,
		LastRestocked:     inventory.LastRestocked.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:         inventory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Inventory retrieved successfully", response)
}

// UpdateInventory updates an existing inventory record
func (h *inventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid inventory ID", err)
		return
	}

	var req dto.UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	inventory, err := h.inventoryService.UpdateInventory(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update inventory", err)
		return
	}

	response := dto.InventoryResponse{
		ID:                inventory.ID,
		ProductID:         inventory.ProductID,
		ProductVariantID:  inventory.ProductVariantID,
		Quantity:          inventory.Quantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		AvailableQuantity: inventory.AvailableQuantity,
		MinStockLevel:     inventory.MinStockLevel,
		MaxStockLevel:     getIntPointer(inventory.MaxStockLevel),
		ReorderPoint:      inventory.ReorderPoint,
		LastRestocked:     inventory.LastRestocked.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:         inventory.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.OK(w, "Inventory updated successfully", response)
}

// DeleteInventory deletes an inventory record
func (h *inventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid inventory ID", err)
		return
	}

	err = h.inventoryService.DeleteInventory(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete inventory", err)
		return
	}

	httpx.OK(w, "Inventory deleted successfully", nil)
}

// ListInventory lists inventory with filters
func (h *inventoryHandler) ListInventory(w http.ResponseWriter, r *http.Request) {
	req := &dto.ListInventoryRequest{
		Page:  1,
		Limit: 20,
	}

	// Parse query parameters
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			req.Limit = limit
		}
	}

	if productIDStr := r.URL.Query().Get("product_id"); productIDStr != "" {
		if productID, err := strconv.ParseInt(productIDStr, 10, 64); err == nil {
			req.ProductID = &productID
		}
	}

	if variantIDStr := r.URL.Query().Get("variant_id"); variantIDStr != "" {
		if variantID, err := strconv.ParseInt(variantIDStr, 10, 64); err == nil {
			req.ProductVariantID = &variantID
		}
	}

	if lowStockStr := r.URL.Query().Get("low_stock"); lowStockStr != "" {
		if lowStock, err := strconv.ParseBool(lowStockStr); err == nil {
			req.LowStock = &lowStock
		}
	}

	if outOfStockStr := r.URL.Query().Get("out_of_stock"); outOfStockStr != "" {
		if outOfStock, err := strconv.ParseBool(outOfStockStr); err == nil {
			req.OutOfStock = &outOfStock
		}
	}

	response, err := h.inventoryService.ListInventory(r.Context(), req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to list inventory", err)
		return
	}

	httpx.OK(w, "Inventory listed successfully", response)
}

// GetInventorySummary gets inventory summary statistics
func (h *inventoryHandler) GetInventorySummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.inventoryService.GetInventorySummary(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get inventory summary", err)
		return
	}

	httpx.OK(w, "Inventory summary retrieved successfully", summary)
}

// Stock Movements

// RecordStockMovement records a stock movement
func (h *inventoryHandler) RecordStockMovement(w http.ResponseWriter, r *http.Request) {
	var req dto.StockMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	movement, err := h.inventoryService.RecordStockMovement(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to record stock movement", err)
		return
	}

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

	httpx.Created(w, "Stock movement recorded successfully", response)
}

// GetStockMovements retrieves stock movements with filters
func (h *inventoryHandler) GetStockMovements(w http.ResponseWriter, r *http.Request) {
	req := &dto.ListStockMovementsRequest{
		Page:  1,
		Limit: 20,
	}

	// Parse query parameters
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			req.Limit = limit
		}
	}

	if productIDStr := r.URL.Query().Get("product_id"); productIDStr != "" {
		if productID, err := strconv.ParseInt(productIDStr, 10, 64); err == nil {
			req.ProductID = &productID
		}
	}

	if variantIDStr := r.URL.Query().Get("variant_id"); variantIDStr != "" {
		if variantID, err := strconv.ParseInt(variantIDStr, 10, 64); err == nil {
			req.ProductVariantID = &variantID
		}
	}

	if movementType := r.URL.Query().Get("movement_type"); movementType != "" {
		req.MovementType = &movementType
	}

	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		req.StartDate = &startDate
	}

	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		req.EndDate = &endDate
	}

	response, err := h.inventoryService.GetStockMovements(r.Context(), req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get stock movements", err)
		return
	}

	httpx.OK(w, "Stock movements retrieved successfully", response)
}

// GetStockMovementByID retrieves a stock movement by ID
func (h *inventoryHandler) GetStockMovementByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid movement ID", err)
		return
	}

	movement, err := h.inventoryService.GetStockMovementByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "Stock movement not found", err)
		return
	}

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

	httpx.OK(w, "Stock movement retrieved successfully", response)
}

// Stock Reservations

// ReserveStock reserves stock for an order
func (h *inventoryHandler) ReserveStock(w http.ResponseWriter, r *http.Request) {
	var req dto.ReserveStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	reservation, err := h.inventoryService.ReserveStock(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to reserve stock", err)
		return
	}

	response := dto.ReserveStockResponse{
		ID:               reservation.ID,
		ProductID:        reservation.ProductID,
		ProductVariantID: reservation.ProductVariantID,
		OrderID:          reservation.OrderID,
		Quantity:         reservation.Quantity,
		ExpiresAt:        reservation.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:        reservation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	httpx.Created(w, "Stock reserved successfully", response)
}

// ReleaseStock releases reserved stock
func (h *inventoryHandler) ReleaseStock(w http.ResponseWriter, r *http.Request) {
	var req dto.ReleaseStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	err := h.inventoryService.ReleaseStock(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to release stock", err)
		return
	}

	httpx.OK(w, "Stock released successfully", nil)
}

// GetStockReservations gets stock reservations for an order
func (h *inventoryHandler) GetStockReservations(w http.ResponseWriter, r *http.Request) {
	orderIDStr := r.URL.Query().Get("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	reservations, err := h.inventoryService.GetStockReservations(r.Context(), orderID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get stock reservations", err)
		return
	}

	var responses []dto.ReserveStockResponse
	for _, reservation := range reservations {
		response := dto.ReserveStockResponse{
			ID:               reservation.ID,
			ProductID:        reservation.ProductID,
			ProductVariantID: reservation.ProductVariantID,
			OrderID:          reservation.OrderID,
			Quantity:         reservation.Quantity,
			ExpiresAt:        reservation.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:        reservation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		responses = append(responses, response)
	}

	httpx.OK(w, "Stock reservations retrieved successfully", responses)
}

// Inventory Alerts

// GetInventoryAlerts gets inventory alerts
func (h *inventoryHandler) GetInventoryAlerts(w http.ResponseWriter, r *http.Request) {
	var resolved *bool
	if resolvedStr := r.URL.Query().Get("resolved"); resolvedStr != "" {
		if resolvedVal, err := strconv.ParseBool(resolvedStr); err == nil {
			resolved = &resolvedVal
		}
	}

	alerts, err := h.inventoryService.GetInventoryAlerts(r.Context(), resolved)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get inventory alerts", err)
		return
	}

	var responses []dto.InventoryAlertResponse
	for _, alert := range alerts {
		response := dto.InventoryAlertResponse{
			ID:                alert.ID,
			ProductID:         alert.ProductID,
			ProductVariantID:  alert.ProductVariantID,
			AlertType:         alert.AlertType,
			CurrentQuantity:   alert.CurrentQuantity,
			ThresholdQuantity: alert.ThresholdQuantity,
			IsResolved:        alert.IsResolved,
			CreatedAt:         alert.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		if alert.ResolvedAt != nil {
			resolvedAt := alert.ResolvedAt.Format("2006-01-02T15:04:05Z07:00")
			response.ResolvedAt = &resolvedAt
		}

		responses = append(responses, response)
	}

	httpx.OK(w, "Inventory alerts retrieved successfully", responses)
}

// ResolveInventoryAlert resolves an inventory alert
func (h *inventoryHandler) ResolveInventoryAlert(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid alert ID", err)
		return
	}

	err = h.inventoryService.ResolveInventoryAlert(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to resolve inventory alert", err)
		return
	}

	httpx.OK(w, "Inventory alert resolved successfully", nil)
}

// CheckLowStockAlerts checks for low stock and creates alerts
func (h *inventoryHandler) CheckLowStockAlerts(w http.ResponseWriter, r *http.Request) {
	err := h.inventoryService.CheckLowStockAlerts(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to check low stock alerts", err)
		return
	}

	httpx.OK(w, "Low stock alerts checked successfully", nil)
}

// Bulk Operations

// BulkUpdateStock performs bulk stock updates
func (h *inventoryHandler) BulkUpdateStock(w http.ResponseWriter, r *http.Request) {
	var req dto.BulkStockUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	response, err := h.inventoryService.BulkUpdateStock(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to perform bulk stock update", err)
		return
	}

	httpx.OK(w, "Bulk stock update completed", response)
}

// Helper functions
func getIntPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func getInt64Pointer(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}
