package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
	"github.com/jattinmanhas/GearboxV2/services/shared/middleware"
)

type OrderHandler interface {
	// Order Management
	CreateOrder(w http.ResponseWriter, r *http.Request)
	GetOrder(w http.ResponseWriter, r *http.Request)
	GetOrderByNumber(w http.ResponseWriter, r *http.Request)
	UpdateOrder(w http.ResponseWriter, r *http.Request)
	DeleteOrder(w http.ResponseWriter, r *http.Request)
	ListOrders(w http.ResponseWriter, r *http.Request)

	// Order Items
	GetOrderItems(w http.ResponseWriter, r *http.Request)
	UpdateOrderItem(w http.ResponseWriter, r *http.Request)
	DeleteOrderItem(w http.ResponseWriter, r *http.Request)

	// Order Addresses
	GetOrderAddresses(w http.ResponseWriter, r *http.Request)
	UpdateOrderAddress(w http.ResponseWriter, r *http.Request)

	// Order Status Management
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request)
	GetOrderStatusHistory(w http.ResponseWriter, r *http.Request)

	// Order Fulfillment
	CreateOrderFulfillment(w http.ResponseWriter, r *http.Request)
	GetOrderFulfillment(w http.ResponseWriter, r *http.Request)
	UpdateOrderFulfillment(w http.ResponseWriter, r *http.Request)

	// Order Refunds
	CreateOrderRefund(w http.ResponseWriter, r *http.Request)
	GetOrderRefunds(w http.ResponseWriter, r *http.Request)

	// Order Analytics
	GetOrderAnalytics(w http.ResponseWriter, r *http.Request)
	GetOrderAnalyticsByDateRange(w http.ResponseWriter, r *http.Request)
	GetTopSellingProducts(w http.ResponseWriter, r *http.Request)

	// Cart Integration
	CreateOrderFromCart(w http.ResponseWriter, r *http.Request)

	// Payment Integration
	CreateOrderPayment(w http.ResponseWriter, r *http.Request)
	ProcessOrderPayment(w http.ResponseWriter, r *http.Request)
	GetOrderPayment(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

// Order Management

// CreateOrder handles POST /api/v1/orders
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		httpx.Error(w, http.StatusUnauthorized, "user ID not found in context", nil)
		return
	}

	// Extract auth token from context (set by middleware)
	authHeader, ok := middleware.ExtractAuthTokenFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "auth token not found in context", nil)
		return
	}

	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Set user ID from context (override any user_id in request body)
	req.UserID = int64(userID)

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	// Add auth token to context for auth client calls
	ctx := context.WithValue(r.Context(), "auth_token", authHeader)
	order, err := h.orderService.CreateOrder(ctx, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create order", err)
		return
	}

	response := dto.OrderResponse{
		ID:                order.ID,
		OrderNumber:       order.OrderNumber,
		UserID:            order.UserID,
		Status:            order.Status,
		PaymentStatus:     order.PaymentStatus,
		FulfillmentStatus: order.FulfillmentStatus,
		Subtotal:          order.Subtotal,
		TaxAmount:         order.TaxAmount,
		ShippingAmount:    order.ShippingAmount,
		DiscountAmount:    order.DiscountAmount,
		TotalAmount:       order.TotalAmount,
		Currency:          order.Currency,
		Notes:             order.Notes,
		InternalNotes:     order.InternalNotes,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		ConfirmedAt:       order.ConfirmedAt,
		ShippedAt:         order.ShippedAt,
		DeliveredAt:       order.DeliveredAt,
		CancelledAt:       order.CancelledAt,
	}

	httpx.Created(w, "Order created successfully", response)
}

// GetOrder handles GET /api/v1/orders/{id}
func (h *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	// Get enriched order details
	response, err := h.orderService.GetOrderDetails(r.Context(), id)
	if err != nil {
		if err.Error() == "order with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to get order details", err)
		}
		return
	}

	httpx.OK(w, "Order retrieved successfully", response)
}

// GetOrderByNumber handles GET /api/v1/orders/number/{orderNumber}
func (h *orderHandler) GetOrderByNumber(w http.ResponseWriter, r *http.Request) {
	orderNumber := chi.URLParam(r, "orderNumber")

	response, err := h.orderService.GetOrderDetailsByNumber(r.Context(), orderNumber)
	if err != nil {
		if err.Error() == "order with number not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to get order details", err)
		}
		return
	}

	httpx.OK(w, "Order retrieved successfully", response)
}

// UpdateOrder handles PUT /api/v1/orders/{id}
func (h *orderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req dto.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	order, err := h.orderService.UpdateOrder(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "order with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to update order", err)
		}
		return
	}

	response := dto.OrderResponse{
		ID:                order.ID,
		OrderNumber:       order.OrderNumber,
		UserID:            order.UserID,
		Status:            order.Status,
		PaymentStatus:     order.PaymentStatus,
		FulfillmentStatus: order.FulfillmentStatus,
		Subtotal:          order.Subtotal,
		TaxAmount:         order.TaxAmount,
		ShippingAmount:    order.ShippingAmount,
		DiscountAmount:    order.DiscountAmount,
		TotalAmount:       order.TotalAmount,
		Currency:          order.Currency,
		Notes:             order.Notes,
		InternalNotes:     order.InternalNotes,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		ConfirmedAt:       order.ConfirmedAt,
		ShippedAt:         order.ShippedAt,
		DeliveredAt:       order.DeliveredAt,
		CancelledAt:       order.CancelledAt,
	}

	httpx.OK(w, "Order updated successfully", response)
}

// DeleteOrder handles DELETE /api/v1/orders/{id}
func (h *orderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	err = h.orderService.DeleteOrder(r.Context(), id)
	if err != nil {
		if err.Error() == "order with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else if err.Error() == "only pending orders can be deleted" {
			httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to delete order", err)
		}
		return
	}

	httpx.OK(w, "Order deleted successfully", nil)
}

// ListOrders handles GET /api/v1/orders
func (h *orderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "created_at" // Default sort field
	}
	sortOrder := r.URL.Query().Get("order")
	if sortOrder == "" {
		sortOrder = "desc" // Default sort order
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc" // Validate sort order
	}

	req := &dto.ListOrdersRequest{
		UserID:            getInt64Param(r, "user_id"),
		Status:            r.URL.Query().Get("status"),
		PaymentStatus:     r.URL.Query().Get("payment_status"),
		FulfillmentStatus: r.URL.Query().Get("fulfillment_status"),
		DateFrom:          getTimeParam(r, "date_from"),
		DateTo:            getTimeParam(r, "date_to"),
		MinAmount:         getFloat64Param(r, "min_amount"),
		MaxAmount:         getFloat64Param(r, "max_amount"),
		Search:            r.URL.Query().Get("search"),
		Page:              getIntParam(r, "page", 1),
		Limit:             getIntParam(r, "limit", 10),
		Sort:              sortBy,
		Order:             sortOrder,
	}

	response, err := h.orderService.ListOrders(r.Context(), req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to list orders", err)
		return
	}

	httpx.OK(w, "Orders retrieved successfully", response)
}

// Order Items

// GetOrderItems handles GET /api/v1/orders/{id}/items
func (h *orderHandler) GetOrderItems(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	items, err := h.orderService.GetOrderItems(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order items", err)
		return
	}

	// Convert to response DTOs
	itemResponses := make([]dto.OrderItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.OrderItemResponse{
			ID:               item.ID,
			OrderID:          item.OrderID,
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			ProductName:      item.ProductName,
			ProductSKU:       item.ProductSKU,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			TotalPrice:       item.TotalPrice,
			TaxAmount:        item.TaxAmount,
			DiscountAmount:   item.DiscountAmount,
			IsDigital:        item.IsDigital,
			RequiresShipping: item.RequiresShipping,
		}
	}

	httpx.OK(w, "Order items retrieved successfully", itemResponses)
}

// UpdateOrderItem handles PUT /api/v1/orders/items/{id}
func (h *orderHandler) UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order item ID", err)
		return
	}

	var item domain.OrderItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = h.orderService.UpdateOrderItem(r.Context(), id, &item)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update order item", err)
		return
	}

	httpx.OK(w, "Order item updated successfully", nil)
}

// DeleteOrderItem handles DELETE /api/v1/orders/items/{id}
func (h *orderHandler) DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order item ID", err)
		return
	}

	err = h.orderService.DeleteOrderItem(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete order item", err)
		return
	}

	httpx.OK(w, "Order item deleted successfully", nil)
}

// Order Addresses

// GetOrderAddresses handles GET /api/v1/orders/{id}/addresses
func (h *orderHandler) GetOrderAddresses(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	addresses, err := h.orderService.GetOrderAddresses(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order addresses", err)
		return
	}

	// Convert to response DTOs
	addressResponses := make([]dto.OrderAddressResponse, len(addresses))
	for i, address := range addresses {
		addressResponses[i] = dto.OrderAddressResponse{
			ID:           address.ID,
			OrderID:      address.OrderID,
			Type:         address.Type,
			FirstName:    address.FirstName,
			LastName:     address.LastName,
			Company:      address.Company,
			AddressLine1: address.Address1,
			AddressLine2: address.Address2,
			City:         address.City,
			State:        address.State,
			Country:      address.Country,
			PostalCode:   address.PostalCode,
			Phone:        address.Phone,
			Email:        address.Email,
		}
	}

	httpx.OK(w, "Order addresses retrieved successfully", addressResponses)
}

// UpdateOrderAddress handles PUT /api/v1/orders/addresses/{id}
func (h *orderHandler) UpdateOrderAddress(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order address ID", err)
		return
	}

	var address domain.OrderAddress
	if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err = h.orderService.UpdateOrderAddress(r.Context(), id, &address)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update order address", err)
		return
	}

	httpx.OK(w, "Order address updated successfully", nil)
}

// Order Status Management

// UpdateOrderStatus handles PUT /api/v1/orders/{id}/status
func (h *orderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	err = h.orderService.UpdateOrderStatus(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "order with ID not found" {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
		} else if err.Error() == "invalid status transition" {
			httpx.Error(w, http.StatusBadRequest, err.Error(), nil)
		} else {
			httpx.Error(w, http.StatusInternalServerError, "Failed to update order status", err)
		}
		return
	}

	httpx.OK(w, "Order status updated successfully", nil)
}

// GetOrderStatusHistory handles GET /api/v1/orders/{id}/status-history
func (h *orderHandler) GetOrderStatusHistory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	history, err := h.orderService.GetOrderStatusHistory(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order status history", err)
		return
	}

	// Convert to response DTOs
	historyResponses := make([]dto.OrderStatusHistoryResponse, len(history))
	for i, h := range history {
		historyResponses[i] = dto.OrderStatusHistoryResponse{
			ID:             h.ID,
			OrderID:        h.OrderID,
			Status:         h.Status,
			PreviousStatus: h.PreviousStatus,
			Notes:          h.Notes,
			CreatedBy:      h.CreatedBy,
			CreatedAt:      h.CreatedAt,
		}
	}

	httpx.OK(w, "Order status history retrieved successfully", historyResponses)
}

// Order Fulfillment

// CreateOrderFulfillment handles POST /api/v1/orders/{id}/fulfillment
func (h *orderHandler) CreateOrderFulfillment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req dto.CreateOrderFulfillmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	fulfillment, err := h.orderService.CreateOrderFulfillment(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create order fulfillment", err)
		return
	}

	response := dto.OrderFulfillmentResponse{
		ID:                fulfillment.ID,
		OrderID:           fulfillment.OrderID,
		TrackingNumber:    fulfillment.TrackingNumber,
		Carrier:           fulfillment.Carrier,
		Service:           fulfillment.Service,
		Status:            fulfillment.Status,
		ShippedAt:         fulfillment.ShippedAt,
		DeliveredAt:       fulfillment.DeliveredAt,
		EstimatedDelivery: fulfillment.EstimatedDelivery,
		Notes:             fulfillment.Notes,
		CreatedAt:         fulfillment.CreatedAt,
		UpdatedAt:         fulfillment.UpdatedAt,
	}

	httpx.Created(w, "Order fulfillment created successfully", response)
}

// GetOrderFulfillment handles GET /api/v1/orders/{id}/fulfillment
func (h *orderHandler) GetOrderFulfillment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	fulfillment, err := h.orderService.GetOrderFulfillment(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order fulfillment", err)
		return
	}

	if fulfillment == nil {
		httpx.Error(w, http.StatusNotFound, "Order fulfillment not found", nil)
		return
	}

	response := dto.OrderFulfillmentResponse{
		ID:                fulfillment.ID,
		OrderID:           fulfillment.OrderID,
		TrackingNumber:    fulfillment.TrackingNumber,
		Carrier:           fulfillment.Carrier,
		Service:           fulfillment.Service,
		Status:            fulfillment.Status,
		ShippedAt:         fulfillment.ShippedAt,
		DeliveredAt:       fulfillment.DeliveredAt,
		EstimatedDelivery: fulfillment.EstimatedDelivery,
		Notes:             fulfillment.Notes,
		CreatedAt:         fulfillment.CreatedAt,
		UpdatedAt:         fulfillment.UpdatedAt,
	}

	httpx.OK(w, "Order fulfillment retrieved successfully", response)
}

// UpdateOrderFulfillment handles PUT /api/v1/orders/fulfillment/{id}
func (h *orderHandler) UpdateOrderFulfillment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid fulfillment ID", err)
		return
	}

	var req dto.UpdateOrderFulfillmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	fulfillment, err := h.orderService.UpdateOrderFulfillment(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update order fulfillment", err)
		return
	}

	response := dto.OrderFulfillmentResponse{
		ID:                fulfillment.ID,
		OrderID:           fulfillment.OrderID,
		TrackingNumber:    fulfillment.TrackingNumber,
		Carrier:           fulfillment.Carrier,
		Service:           fulfillment.Service,
		Status:            fulfillment.Status,
		ShippedAt:         fulfillment.ShippedAt,
		DeliveredAt:       fulfillment.DeliveredAt,
		EstimatedDelivery: fulfillment.EstimatedDelivery,
		Notes:             fulfillment.Notes,
		CreatedAt:         fulfillment.CreatedAt,
		UpdatedAt:         fulfillment.UpdatedAt,
	}

	httpx.OK(w, "Order fulfillment updated successfully", response)
}

// Order Refunds

// CreateOrderRefund handles POST /api/v1/orders/{id}/refunds
func (h *orderHandler) CreateOrderRefund(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	var req dto.CreateOrderRefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	refund, err := h.orderService.CreateOrderRefund(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create order refund", err)
		return
	}

	response := dto.OrderRefundResponse{
		ID:          refund.ID,
		OrderID:     refund.OrderID,
		Amount:      refund.Amount,
		Reason:      refund.Reason,
		Status:      refund.Status,
		ProcessedAt: refund.ProcessedAt,
		CreatedBy:   refund.CreatedBy,
		CreatedAt:   refund.CreatedAt,
	}

	httpx.Created(w, "Order refund created successfully", response)
}

// GetOrderRefunds handles GET /api/v1/orders/{id}/refunds
func (h *orderHandler) GetOrderRefunds(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	refunds, err := h.orderService.GetOrderRefunds(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order refunds", err)
		return
	}

	// Convert to response DTOs
	refundResponses := make([]dto.OrderRefundResponse, len(refunds))
	for i, refund := range refunds {
		refundResponses[i] = dto.OrderRefundResponse{
			ID:          refund.ID,
			OrderID:     refund.OrderID,
			Amount:      refund.Amount,
			Reason:      refund.Reason,
			Status:      refund.Status,
			ProcessedAt: refund.ProcessedAt,
			CreatedBy:   refund.CreatedBy,
			CreatedAt:   refund.CreatedAt,
		}
	}

	httpx.OK(w, "Order refunds retrieved successfully", refundResponses)
}

// Order Analytics

// GetOrderAnalytics handles GET /api/v1/orders/analytics
func (h *orderHandler) GetOrderAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.orderService.GetOrderAnalytics(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order analytics", err)
		return
	}

	httpx.OK(w, "Order analytics retrieved successfully", analytics)
}

// GetOrderAnalyticsByDateRange handles GET /api/v1/orders/analytics/date-range
func (h *orderHandler) GetOrderAnalyticsByDateRange(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		httpx.Error(w, http.StatusBadRequest, "start_date and end_date are required", nil)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD", err)
		return
	}

	// Add time to end date to include the full day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	analytics, err := h.orderService.GetOrderAnalyticsByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get order analytics by date range", err)
		return
	}

	httpx.OK(w, "Order analytics by date range retrieved successfully", analytics)
}

// GetTopSellingProducts handles GET /api/v1/orders/analytics/top-products
func (h *orderHandler) GetTopSellingProducts(w http.ResponseWriter, r *http.Request) {
	// Parse limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			httpx.Error(w, http.StatusBadRequest, "Invalid limit. Must be between 1 and 100", err)
			return
		}
	}

	products, err := h.orderService.GetTopSellingProducts(r.Context(), limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get top selling products", err)
		return
	}

	httpx.OK(w, "Top selling products retrieved successfully", products)
}

// Cart Integration

// CreateOrderFromCart handles POST /api/v1/orders/from-cart
func (h *orderHandler) CreateOrderFromCart(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		httpx.Error(w, http.StatusUnauthorized, "user ID not found in context", nil)
		return
	}

	// Parse cart ID from query parameters
	cartIDStr := r.URL.Query().Get("cart_id")
	if cartIDStr == "" {
		httpx.Error(w, http.StatusBadRequest, "cart_id is required", nil)
		return
	}

	cartID, err := strconv.ParseInt(cartIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid cart_id", err)
		return
	}

	// Extract auth token from context (set by middleware)
	authHeader, ok := middleware.ExtractAuthTokenFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "auth token not found in context", nil)
		return
	}

	var req dto.CreateOrderFromCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	// Add auth token to context for auth client calls
	ctx := context.WithValue(r.Context(), "auth_token", authHeader)
	order, err := h.orderService.CreateOrderFromCart(ctx, int64(userID), cartID, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create order from cart", err)
		return
	}

	response := dto.OrderResponse{
		ID:                order.ID,
		OrderNumber:       order.OrderNumber,
		UserID:            order.UserID,
		Status:            order.Status,
		PaymentStatus:     order.PaymentStatus,
		FulfillmentStatus: order.FulfillmentStatus,
		Subtotal:          order.Subtotal,
		TaxAmount:         order.TaxAmount,
		ShippingAmount:    order.ShippingAmount,
		DiscountAmount:    order.DiscountAmount,
		TotalAmount:       order.TotalAmount,
		Currency:          order.Currency,
		Notes:             order.Notes,
		InternalNotes:     order.InternalNotes,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
		ConfirmedAt:       order.ConfirmedAt,
		ShippedAt:         order.ShippedAt,
		DeliveredAt:       order.DeliveredAt,
		CancelledAt:       order.CancelledAt,
	}

	httpx.Created(w, "Order created from cart successfully", response)
}

// Helper functions

func getFloat64Param(r *http.Request, key string) *float64 {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return &floatValue
	}
	return nil
}

func getTimeParam(r *http.Request, key string) *time.Time {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	if timeValue, err := time.Parse("2006-01-02", value); err == nil {
		return &timeValue
	}
	return nil
}

// Payment Integration Handlers

// CreateOrderPayment creates a payment for an order
func (h *orderHandler) CreateOrderPayment(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpx.BadRequest(w, "Invalid order ID", nil)
		return
	}

	var req dto.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body", nil)
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		httpx.BadRequest(w, "Validation failed", map[string]interface{}{
			"errors": err,
		})
		return
	}

	// Create payment
	payment, err := h.orderService.CreatePaymentForOrder(r.Context(), orderID, &req)
	if err != nil {
		httpx.InternalServerError(w, "Failed to create payment", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	httpx.Created(w, "Payment created successfully", map[string]interface{}{
		"payment": payment,
	})
}

// ProcessOrderPayment processes a payment for an order
func (h *orderHandler) ProcessOrderPayment(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpx.BadRequest(w, "Invalid order ID", nil)
		return
	}

	var req dto.ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.BadRequest(w, "Invalid request body", nil)
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		httpx.BadRequest(w, "Validation failed", map[string]interface{}{
			"errors": err,
		})
		return
	}

	// Process payment
	payment, err := h.orderService.ProcessOrderPayment(r.Context(), orderID, &req)
	if err != nil {
		httpx.InternalServerError(w, "Failed to process payment", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	httpx.OK(w, "Payment processed successfully", map[string]interface{}{
		"payment": payment,
	})
}

// GetOrderPayment retrieves payment information for an order
func (h *orderHandler) GetOrderPayment(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpx.BadRequest(w, "Invalid order ID", nil)
		return
	}

	// Get payment
	payment, err := h.orderService.GetOrderPayment(r.Context(), orderID)
	if err != nil {
		httpx.InternalServerError(w, "Failed to get payment", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if payment == nil {
		httpx.NotFound(w, "Payment not found", nil)
		return
	}

	httpx.OK(w, "Payment retrieved successfully", map[string]interface{}{
		"payment": payment,
	})
}
