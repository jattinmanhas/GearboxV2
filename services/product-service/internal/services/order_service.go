package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
)

type OrderService interface {
	// Order Management
	CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*domain.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	UpdateOrder(ctx context.Context, id int64, req *dto.UpdateOrderRequest) (*domain.Order, error)
	DeleteOrder(ctx context.Context, id int64) error
	ListOrders(ctx context.Context, req *dto.ListOrdersRequest) (*dto.ListOrdersResponse, error)

	// Order Items
	GetOrderItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error)
	UpdateOrderItem(ctx context.Context, id int64, item *domain.OrderItem) error
	DeleteOrderItem(ctx context.Context, id int64) error

	// Order Addresses
	GetOrderAddresses(ctx context.Context, orderID int64) ([]*domain.OrderAddress, error)
	UpdateOrderAddress(ctx context.Context, id int64, address *domain.OrderAddress) error

	// Order Status Management
	UpdateOrderStatus(ctx context.Context, orderID int64, req *dto.UpdateOrderStatusRequest) error
	GetOrderStatusHistory(ctx context.Context, orderID int64) ([]*domain.OrderStatusHistory, error)

	// Order Fulfillment
	CreateOrderFulfillment(ctx context.Context, orderID int64, req *dto.CreateOrderFulfillmentRequest) (*domain.OrderFulfillment, error)
	GetOrderFulfillment(ctx context.Context, orderID int64) (*domain.OrderFulfillment, error)
	UpdateOrderFulfillment(ctx context.Context, id int64, req *dto.UpdateOrderFulfillmentRequest) (*domain.OrderFulfillment, error)

	// Order Refunds
	CreateOrderRefund(ctx context.Context, orderID int64, req *dto.CreateOrderRefundRequest) (*domain.OrderRefund, error)
	GetOrderRefunds(ctx context.Context, orderID int64) ([]*domain.OrderRefund, error)

	// Order Analytics
	GetOrderAnalytics(ctx context.Context) (*dto.OrderAnalyticsResponse, error)
	GetOrderAnalyticsByDateRange(ctx context.Context, startDate, endDate time.Time) (*dto.OrderAnalyticsResponse, error)
	GetTopSellingProducts(ctx context.Context, limit int) ([]*dto.ProductOrderStatsResponse, error)

	// Cart Integration
	CreateOrderFromCart(ctx context.Context, userID int64, cartID int64, req *dto.CreateOrderFromCartRequest) (*domain.Order, error)
}

type orderService struct {
	orderRepo        repository.OrderRepository
	productRepo      repository.ProductRepository
	inventoryService InventoryService
	cartService      CartService
	couponService    CouponService
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository, inventoryService InventoryService, cartService CartService, couponService CouponService) OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		productRepo:      productRepo,
		inventoryService: inventoryService,
		cartService:      cartService,
		couponService:    couponService,
	}
}

// Order Management

// CreateOrder creates a new order
func (s *orderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*domain.Order, error) {
	// Validate and check inventory for all items
	for _, item := range req.Items {
		// Check if product exists
		_, err := s.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found: %w", item.ProductID, err)
		}

		// Check inventory availability
		available, err := s.inventoryService.CheckStockAvailability(ctx, item.ProductID, item.ProductVariantID, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to check stock for product %d: %w", item.ProductID, err)
		}
		if !available {
			return nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}
	}

	// Calculate order totals
	subtotal := 0.0
	items := make([]*domain.OrderItem, len(req.Items))

	for i, item := range req.Items {
		product, _ := s.productRepo.GetProductByID(ctx, item.ProductID)

		totalPrice := item.UnitPrice * float64(item.Quantity)
		subtotal += totalPrice

		productName := ""
		productSKU := ""
		isDigital := false
		requiresShipping := true

		if product != nil {
			productName = product.Name
			productSKU = product.SKU
			isDigital = !product.RequiresShipping
			requiresShipping = product.RequiresShipping
		}

		items[i] = &domain.OrderItem{
			OrderID:          0, // Will be set after order creation
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			ProductName:      productName,
			ProductSKU:       productSKU,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			TotalPrice:       totalPrice,
			IsDigital:        isDigital,
			RequiresShipping: requiresShipping,
		}
	}

	// Apply coupon discounts if provided
	discountAmount := 0.0
	if len(req.ApplyCoupons) > 0 {
		for _, couponCode := range req.ApplyCoupons {
			// Validate coupon
			validateReq := &dto.ValidateCouponRequest{
				CouponCode: couponCode,
				CartAmount: subtotal,
				UserID:     &req.UserID,
			}

			validateResp, err := s.couponService.ValidateCoupon(ctx, validateReq)
			if err == nil && validateResp.Valid {
				discountAmount += validateResp.DiscountAmount
			}
		}
	}

	// Calculate tax (placeholder - would integrate with tax service)
	taxAmount := subtotal * 0.08 // 8% tax rate

	// Calculate shipping (placeholder - would integrate with shipping service)
	shippingAmount := 10.0 // Fixed shipping cost

	totalAmount := subtotal + taxAmount + shippingAmount - discountAmount

	// Create order
	order := &domain.Order{
		UserID:            req.UserID,
		Status:            "pending",
		PaymentStatus:     "pending",
		FulfillmentStatus: "unfulfilled",
		Subtotal:          subtotal,
		TaxAmount:         taxAmount,
		ShippingAmount:    shippingAmount,
		DiscountAmount:    discountAmount,
		TotalAmount:       totalAmount,
		Currency:          req.Currency,
		Notes:             req.Notes,
		InternalNotes:     req.InternalNotes,
	}

	// Create order in database
	err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Set order ID for items
	for _, item := range items {
		item.OrderID = order.ID
	}

	// Create order items
	err = s.orderRepo.CreateOrderItems(ctx, items)
	if err != nil {
		// Rollback order creation
		s.orderRepo.DeleteOrder(ctx, order.ID)
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// Create addresses
	addresses := []*domain.OrderAddress{
		{
			OrderID:    order.ID,
			Type:       "shipping",
			FirstName:  req.ShippingAddress.FirstName,
			LastName:   req.ShippingAddress.LastName,
			Company:    req.ShippingAddress.Company,
			Address1:   req.ShippingAddress.Address1,
			Address2:   req.ShippingAddress.Address2,
			City:       req.ShippingAddress.City,
			State:      req.ShippingAddress.State,
			Country:    req.ShippingAddress.Country,
			PostalCode: req.ShippingAddress.PostalCode,
			Phone:      req.ShippingAddress.Phone,
			Email:      req.ShippingAddress.Email,
		},
	}

	// Add billing address if different from shipping
	if req.BillingAddress != nil {
		addresses = append(addresses, &domain.OrderAddress{
			OrderID:    order.ID,
			Type:       "billing",
			FirstName:  req.BillingAddress.FirstName,
			LastName:   req.BillingAddress.LastName,
			Company:    req.BillingAddress.Company,
			Address1:   req.BillingAddress.Address1,
			Address2:   req.BillingAddress.Address2,
			City:       req.BillingAddress.City,
			State:      req.BillingAddress.State,
			Country:    req.BillingAddress.Country,
			PostalCode: req.BillingAddress.PostalCode,
			Phone:      req.BillingAddress.Phone,
			Email:      req.BillingAddress.Email,
		})
	} else {
		// Use shipping address as billing address
		addresses = append(addresses, &domain.OrderAddress{
			OrderID:    order.ID,
			Type:       "billing",
			FirstName:  req.ShippingAddress.FirstName,
			LastName:   req.ShippingAddress.LastName,
			Company:    req.ShippingAddress.Company,
			Address1:   req.ShippingAddress.Address1,
			Address2:   req.ShippingAddress.Address2,
			City:       req.ShippingAddress.City,
			State:      req.ShippingAddress.State,
			Country:    req.ShippingAddress.Country,
			PostalCode: req.ShippingAddress.PostalCode,
			Phone:      req.ShippingAddress.Phone,
			Email:      req.ShippingAddress.Email,
		})
	}

	err = s.orderRepo.CreateOrderAddresses(ctx, addresses)
	if err != nil {
		// Rollback order creation
		s.orderRepo.DeleteOrder(ctx, order.ID)
		return nil, fmt.Errorf("failed to create order addresses: %w", err)
	}

	// Create initial status history
	statusHistory := &domain.OrderStatusHistory{
		OrderID:        order.ID,
		Status:         "pending",
		PreviousStatus: "",
		Notes:          "Order created",
		CreatedBy:      req.UserID,
	}

	err = s.orderRepo.CreateOrderStatusHistory(ctx, statusHistory)
	if err != nil {
		// Log error but don't fail the order creation
		fmt.Printf("Warning: failed to create status history: %v\n", err)
	}

	// Reserve inventory
	for _, item := range items {
		reserveReq := &dto.ReserveStockRequest{
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			OrderID:          order.ID,
			Quantity:         item.Quantity,
			ExpiresAt:        time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z07:00"),
		}

		_, err = s.inventoryService.ReserveStock(ctx, reserveReq)
		if err != nil {
			// Log error but don't fail the order creation
			fmt.Printf("Warning: failed to reserve stock for product %d: %v\n", item.ProductID, err)
		}
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (s *orderService) GetOrderByID(ctx context.Context, id int64) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// GetOrderByNumber retrieves an order by order number
func (s *orderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderByNumber(ctx, orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// UpdateOrder updates an existing order
func (s *orderService) UpdateOrder(ctx context.Context, id int64, req *dto.UpdateOrderRequest) (*domain.Order, error) {
	// Get existing order
	existingOrder, err := s.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing order: %w", err)
	}

	// Update fields that are provided
	updateOrder := *existingOrder

	if req.Status != nil {
		updateOrder.Status = *req.Status
	}
	if req.PaymentStatus != nil {
		updateOrder.PaymentStatus = *req.PaymentStatus
	}
	if req.FulfillmentStatus != nil {
		updateOrder.FulfillmentStatus = *req.FulfillmentStatus
	}
	if req.Notes != nil {
		updateOrder.Notes = *req.Notes
	}
	if req.InternalNotes != nil {
		updateOrder.InternalNotes = *req.InternalNotes
	}

	updateOrder.UpdatedAt = time.Now()

	// Update order in repository
	err = s.orderRepo.UpdateOrder(ctx, id, &updateOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return &updateOrder, nil
}

// DeleteOrder deletes an order
func (s *orderService) DeleteOrder(ctx context.Context, id int64) error {
	// Check if order can be deleted (only pending orders)
	order, err := s.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if order.Status != "pending" {
		return fmt.Errorf("only pending orders can be deleted")
	}

	// Release reserved inventory
	_, err = s.orderRepo.GetOrderItems(ctx, id)
	if err == nil {
		releaseReq := &dto.ReleaseStockRequest{
			OrderID: &id,
		}
		s.inventoryService.ReleaseStock(ctx, releaseReq)
	}

	// Delete order
	err = s.orderRepo.DeleteOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	return nil
}

// ListOrders retrieves orders with filters
func (s *orderService) ListOrders(ctx context.Context, req *dto.ListOrdersRequest) (*dto.ListOrdersResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	offset := (req.Page - 1) * req.Limit

	// Build filter
	filter := &domain.OrderFilter{
		UserID:            req.UserID,
		Status:            req.Status,
		PaymentStatus:     req.PaymentStatus,
		FulfillmentStatus: req.FulfillmentStatus,
		DateFrom:          req.DateFrom,
		DateTo:            req.DateTo,
		MinAmount:         req.MinAmount,
		MaxAmount:         req.MaxAmount,
		Search:            req.Search,
	}

	// Get orders from repository
	orders, total, err := s.orderRepo.ListOrders(ctx, filter, offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Convert to response DTOs
	orderResponses := make([]dto.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = dto.OrderResponse{
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
	}

	// Calculate total pages
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &dto.ListOrdersResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Order Items

// GetOrderItems retrieves order items
func (s *orderService) GetOrderItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error) {
	items, err := s.orderRepo.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return items, nil
}

// UpdateOrderItem updates an order item
func (s *orderService) UpdateOrderItem(ctx context.Context, id int64, item *domain.OrderItem) error {
	err := s.orderRepo.UpdateOrderItem(ctx, id, item)
	if err != nil {
		return fmt.Errorf("failed to update order item: %w", err)
	}

	return nil
}

// DeleteOrderItem deletes an order item
func (s *orderService) DeleteOrderItem(ctx context.Context, id int64) error {
	err := s.orderRepo.DeleteOrderItem(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	return nil
}

// Order Addresses

// GetOrderAddresses retrieves order addresses
func (s *orderService) GetOrderAddresses(ctx context.Context, orderID int64) ([]*domain.OrderAddress, error) {
	addresses, err := s.orderRepo.GetOrderAddresses(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order addresses: %w", err)
	}

	return addresses, nil
}

// UpdateOrderAddress updates an order address
func (s *orderService) UpdateOrderAddress(ctx context.Context, id int64, address *domain.OrderAddress) error {
	err := s.orderRepo.UpdateOrderAddress(ctx, id, address)
	if err != nil {
		return fmt.Errorf("failed to update order address: %w", err)
	}

	return nil
}

// Order Status Management

// UpdateOrderStatus updates order status
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID int64, req *dto.UpdateOrderStatusRequest) error {
	// Get existing order
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, req.Status) {
		return fmt.Errorf("invalid status transition from %s to %s", order.Status, req.Status)
	}

	// Update order status
	updateReq := &dto.UpdateOrderRequest{
		Status: &req.Status,
	}

	_, err = s.UpdateOrder(ctx, orderID, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Create status history
	statusHistory := &domain.OrderStatusHistory{
		OrderID:        orderID,
		Status:         req.Status,
		PreviousStatus: order.Status,
		Notes:          req.Notes,
		CreatedBy:      0, // System update
	}

	err = s.orderRepo.CreateOrderStatusHistory(ctx, statusHistory)
	if err != nil {
		// Log error but don't fail the status update
		fmt.Printf("Warning: failed to create status history: %v\n", err)
	}

	// Handle status-specific actions
	switch req.Status {
	case "confirmed":
		// Reserve inventory permanently
		items, err := s.orderRepo.GetOrderItems(ctx, orderID)
		if err == nil {
			for _, item := range items {
				reference := fmt.Sprintf("order_%d", orderID)
				reason := "order_confirmed"
				movementReq := &dto.StockMovementRequest{
					ProductID:        item.ProductID,
					ProductVariantID: item.ProductVariantID,
					MovementType:     "out",
					Quantity:         item.Quantity,
					Reference:        &reference,
					Reason:           &reason,
				}
				s.inventoryService.RecordStockMovement(ctx, movementReq)
			}
		}
	case "cancelled":
		// Release reserved inventory
		_, err = s.orderRepo.GetOrderItems(ctx, orderID)
		if err == nil {
			releaseReq := &dto.ReleaseStockRequest{
				OrderID: &orderID,
			}
			s.inventoryService.ReleaseStock(ctx, releaseReq)
		}
	}

	return nil
}

// GetOrderStatusHistory retrieves order status history
func (s *orderService) GetOrderStatusHistory(ctx context.Context, orderID int64) ([]*domain.OrderStatusHistory, error) {
	history, err := s.orderRepo.GetOrderStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status history: %w", err)
	}

	return history, nil
}

// Order Fulfillment

// CreateOrderFulfillment creates order fulfillment
func (s *orderService) CreateOrderFulfillment(ctx context.Context, orderID int64, req *dto.CreateOrderFulfillmentRequest) (*domain.OrderFulfillment, error) {
	fulfillment := &domain.OrderFulfillment{
		OrderID:           orderID,
		TrackingNumber:    req.TrackingNumber,
		Carrier:           req.Carrier,
		Service:           req.Service,
		Status:            "pending",
		EstimatedDelivery: req.EstimatedDelivery,
		Notes:             req.Notes,
	}

	err := s.orderRepo.CreateOrderFulfillment(ctx, fulfillment)
	if err != nil {
		return nil, fmt.Errorf("failed to create order fulfillment: %w", err)
	}

	return fulfillment, nil
}

// GetOrderFulfillment retrieves order fulfillment
func (s *orderService) GetOrderFulfillment(ctx context.Context, orderID int64) (*domain.OrderFulfillment, error) {
	fulfillment, err := s.orderRepo.GetOrderFulfillment(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order fulfillment: %w", err)
	}

	return fulfillment, nil
}

// UpdateOrderFulfillment updates order fulfillment
func (s *orderService) UpdateOrderFulfillment(ctx context.Context, id int64, req *dto.UpdateOrderFulfillmentRequest) (*domain.OrderFulfillment, error) {
	// Get existing fulfillment
	fulfillment, err := s.orderRepo.GetOrderFulfillment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing fulfillment: %w", err)
	}

	// Update fields that are provided
	if req.TrackingNumber != nil {
		fulfillment.TrackingNumber = *req.TrackingNumber
	}
	if req.Carrier != nil {
		fulfillment.Carrier = *req.Carrier
	}
	if req.Service != nil {
		fulfillment.Service = *req.Service
	}
	if req.Status != nil {
		fulfillment.Status = *req.Status
	}
	if req.EstimatedDelivery != nil {
		fulfillment.EstimatedDelivery = req.EstimatedDelivery
	}
	if req.Notes != nil {
		fulfillment.Notes = *req.Notes
	}

	// Update fulfillment
	err = s.orderRepo.UpdateOrderFulfillment(ctx, id, fulfillment)
	if err != nil {
		return nil, fmt.Errorf("failed to update order fulfillment: %w", err)
	}

	return fulfillment, nil
}

// Order Refunds

// CreateOrderRefund creates order refund
func (s *orderService) CreateOrderRefund(ctx context.Context, orderID int64, req *dto.CreateOrderRefundRequest) (*domain.OrderRefund, error) {
	refund := &domain.OrderRefund{
		OrderID:   orderID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    "pending",
		CreatedBy: req.CreatedBy,
	}

	err := s.orderRepo.CreateOrderRefund(ctx, refund)
	if err != nil {
		return nil, fmt.Errorf("failed to create order refund: %w", err)
	}

	return refund, nil
}

// GetOrderRefunds retrieves order refunds
func (s *orderService) GetOrderRefunds(ctx context.Context, orderID int64) ([]*domain.OrderRefund, error) {
	refunds, err := s.orderRepo.GetOrderRefunds(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order refunds: %w", err)
	}

	return refunds, nil
}

// Order Analytics

// GetOrderAnalytics retrieves order analytics
func (s *orderService) GetOrderAnalytics(ctx context.Context) (*dto.OrderAnalyticsResponse, error) {
	analytics, err := s.orderRepo.GetOrderAnalytics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get order analytics: %w", err)
	}

	return &dto.OrderAnalyticsResponse{
		TotalOrders:        analytics.TotalOrders,
		PendingOrders:      analytics.PendingOrders,
		ConfirmedOrders:    analytics.ConfirmedOrders,
		ProcessingOrders:   analytics.ProcessingOrders,
		ShippedOrders:      analytics.ShippedOrders,
		DeliveredOrders:    analytics.DeliveredOrders,
		CancelledOrders:    analytics.CancelledOrders,
		TotalRevenue:       analytics.TotalRevenue,
		AverageOrderValue:  analytics.AverageOrderValue,
		ConversionRate:     analytics.ConversionRate,
		NewOrdersToday:     analytics.NewOrdersToday,
		NewOrdersThisWeek:  analytics.NewOrdersThisWeek,
		NewOrdersThisMonth: analytics.NewOrdersThisMonth,
	}, nil
}

// GetOrderAnalyticsByDateRange retrieves order analytics by date range
func (s *orderService) GetOrderAnalyticsByDateRange(ctx context.Context, startDate, endDate time.Time) (*dto.OrderAnalyticsResponse, error) {
	analytics, err := s.orderRepo.GetOrderAnalyticsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get order analytics by date range: %w", err)
	}

	return &dto.OrderAnalyticsResponse{
		TotalOrders:        analytics.TotalOrders,
		PendingOrders:      analytics.PendingOrders,
		ConfirmedOrders:    analytics.ConfirmedOrders,
		ProcessingOrders:   analytics.ProcessingOrders,
		ShippedOrders:      analytics.ShippedOrders,
		DeliveredOrders:    analytics.DeliveredOrders,
		CancelledOrders:    analytics.CancelledOrders,
		TotalRevenue:       analytics.TotalRevenue,
		AverageOrderValue:  analytics.AverageOrderValue,
		ConversionRate:     analytics.ConversionRate,
		NewOrdersToday:     0, // Not applicable for date range
		NewOrdersThisWeek:  0, // Not applicable for date range
		NewOrdersThisMonth: 0, // Not applicable for date range
	}, nil
}

// GetTopSellingProducts retrieves top selling products
func (s *orderService) GetTopSellingProducts(ctx context.Context, limit int) ([]*dto.ProductOrderStatsResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	stats, err := s.orderRepo.GetTopSellingProducts(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top selling products: %w", err)
	}

	responses := make([]*dto.ProductOrderStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = &dto.ProductOrderStatsResponse{
			ProductID:     stat.ProductID,
			ProductName:   stat.ProductName,
			SKU:           stat.SKU,
			TotalQuantity: stat.TotalQuantity,
			TotalRevenue:  stat.TotalRevenue,
			OrderCount:    stat.OrderCount,
			AveragePrice:  stat.AveragePrice,
		}
	}

	return responses, nil
}

// Cart Integration

// CreateOrderFromCart creates an order from a cart
func (s *orderService) CreateOrderFromCart(ctx context.Context, userID int64, cartID int64, req *dto.CreateOrderFromCartRequest) (*domain.Order, error) {
	// Get cart items
	cartItems, err := s.cartService.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Convert cart items to order items
	orderItems := make([]dto.CreateOrderItemRequest, len(cartItems))
	for i, item := range cartItems {
		orderItems[i] = dto.CreateOrderItemRequest{
			ProductID:        item.ProductID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
		}
	}

	// Create order request
	orderReq := &dto.CreateOrderRequest{
		UserID:          userID,
		CartID:          &cartID,
		Items:           orderItems,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		PaymentMethodID: req.PaymentMethodID,
		Currency:        req.Currency,
		Notes:           req.Notes,
		ApplyCoupons:    req.ApplyCoupons,
	}

	// Create order
	order, err := s.CreateOrder(ctx, orderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create order from cart: %w", err)
	}

	// Clear cart after successful order creation
	err = s.cartService.ClearCart(ctx, cartID)
	if err != nil {
		// Log error but don't fail the order creation
		fmt.Printf("Warning: failed to clear cart after order creation: %v\n", err)
	}

	return order, nil
}

// Helper functions

// isValidStatusTransition validates if a status transition is allowed
func (s *orderService) isValidStatusTransition(currentStatus, newStatus string) bool {
	validTransitions := map[string][]string{
		"pending":    {"confirmed", "cancelled"},
		"confirmed":  {"processing", "cancelled"},
		"processing": {"shipped", "cancelled"},
		"shipped":    {"delivered", "cancelled"},
		"delivered":  {"refunded"},
		"cancelled":  {},
		"refunded":   {},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}

	return false
}
