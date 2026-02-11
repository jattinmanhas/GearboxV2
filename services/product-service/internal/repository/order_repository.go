package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	// Order CRUD
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id int64) (*domain.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	UpdateOrder(ctx context.Context, id int64, order *domain.Order) error
	DeleteOrder(ctx context.Context, id int64) error
	ListOrders(ctx context.Context, filter *domain.OrderFilter, offset, limit int) ([]*domain.Order, int64, error)

	// Order Items
	CreateOrderItems(ctx context.Context, items []*domain.OrderItem) error
	GetOrderItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error)
	UpdateOrderItem(ctx context.Context, id int64, item *domain.OrderItem) error
	DeleteOrderItem(ctx context.Context, id int64) error

	// Order Addresses
	CreateOrderAddresses(ctx context.Context, addresses []*domain.OrderAddress) error
	GetOrderAddresses(ctx context.Context, orderID int64) ([]*domain.OrderAddress, error)
	UpdateOrderAddress(ctx context.Context, id int64, address *domain.OrderAddress) error
	DeleteOrderAddress(ctx context.Context, id int64) error

	// Order Status History
	CreateOrderStatusHistory(ctx context.Context, history *domain.OrderStatusHistory) error
	GetOrderStatusHistory(ctx context.Context, orderID int64) ([]*domain.OrderStatusHistory, error)

	// Order Fulfillment
	CreateOrderFulfillment(ctx context.Context, fulfillment *domain.OrderFulfillment) error
	GetOrderFulfillment(ctx context.Context, orderID int64) (*domain.OrderFulfillment, error)
	UpdateOrderFulfillment(ctx context.Context, id int64, fulfillment *domain.OrderFulfillment) error

	// Order Refunds
	CreateOrderRefund(ctx context.Context, refund *domain.OrderRefund) error
	GetOrderRefunds(ctx context.Context, orderID int64) ([]*domain.OrderRefund, error)
	UpdateOrderRefund(ctx context.Context, id int64, refund *domain.OrderRefund) error

	// Order Analytics
	GetOrderAnalytics(ctx context.Context) (*domain.OrderAnalytics, error)
	GetOrderAnalyticsByDateRange(ctx context.Context, startDate, endDate time.Time) (*domain.OrderAnalytics, error)
	GetTopSellingProducts(ctx context.Context, limit int) ([]*domain.ProductOrderStats, error)
	GetOrderConversionRate(ctx context.Context) (float64, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Order CRUD

// CreateOrder creates a new order
func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	// Generate order number if not provided
	if order.OrderNumber == "" {
		var orderNumber string
		err := r.db.GetContext(ctx, &orderNumber, "SELECT generate_order_number()")
		if err != nil {
			return fmt.Errorf("failed to generate order number: %w", err)
		}
		order.OrderNumber = orderNumber
	}

	query := `
		INSERT INTO orders (order_number, user_id, status, payment_status, fulfillment_status,
			subtotal, discount_amount, total_amount, currency,
			notes, internal_notes, created_at, updated_at, confirmed_at, shipped_at, delivered_at, cancelled_at)
		VALUES (:order_number, :user_id, :status, :payment_status, :fulfillment_status,
			:subtotal, :discount_amount, :total_amount, :currency,
			:notes, :internal_notes, :created_at, :updated_at, :confirmed_at, :shipped_at, :delivered_at, :cancelled_at)
		RETURNING id`

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	rows, err := r.db.NamedQueryContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&order.ID); err != nil {
			return fmt.Errorf("failed to get order ID: %w", err)
		}
	}

	return nil
}

// GetOrderByID retrieves an order by ID
func (r *orderRepository) GetOrderByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `SELECT * FROM orders WHERE id = $1`

	var order domain.Order
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetOrderByNumber retrieves an order by order number
func (r *orderRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `SELECT * FROM orders WHERE order_number = $1`

	var order domain.Order
	err := r.db.GetContext(ctx, &order, query, orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with number %s not found", orderNumber)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// UpdateOrder updates an existing order
func (r *orderRepository) UpdateOrder(ctx context.Context, id int64, order *domain.Order) error {
	query := `
		UPDATE orders SET 
			status = :status, payment_status = :payment_status, fulfillment_status = :fulfillment_status,
			subtotal = :subtotal,
			discount_amount = :discount_amount, total_amount = :total_amount, currency = :currency,
			notes = :notes, internal_notes = :internal_notes, updated_at = :updated_at,
			confirmed_at = :confirmed_at, shipped_at = :shipped_at, delivered_at = :delivered_at, cancelled_at = :cancelled_at
		WHERE id = :id`

	order.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}

// DeleteOrder deletes an order
func (r *orderRepository) DeleteOrder(ctx context.Context, id int64) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	return nil
}

// ListOrders retrieves orders with filters
func (r *orderRepository) ListOrders(ctx context.Context, filter *domain.OrderFilter, offset, limit int) ([]*domain.Order, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.UserID != nil {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.PaymentStatus != "" {
		whereClause += fmt.Sprintf(" AND payment_status = $%d", argIndex)
		args = append(args, filter.PaymentStatus)
		argIndex++
	}

	if filter.FulfillmentStatus != "" {
		whereClause += fmt.Sprintf(" AND fulfillment_status = $%d", argIndex)
		args = append(args, filter.FulfillmentStatus)
		argIndex++
	}

	if filter.DateFrom != nil {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filter.DateTo)
		argIndex++
	}

	if filter.MinAmount != nil {
		whereClause += fmt.Sprintf(" AND total_amount >= $%d", argIndex)
		args = append(args, *filter.MinAmount)
		argIndex++
	}

	if filter.MaxAmount != nil {
		whereClause += fmt.Sprintf(" AND total_amount <= $%d", argIndex)
		args = append(args, *filter.MaxAmount)
		argIndex++
	}

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (order_number ILIKE $%d OR notes ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orders %s", whereClause)
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %w", err)
	}

	// Determine sort field and order
	sortBy := "created_at"
	if filter.SortBy != "" {
		// Validate sort field to prevent SQL injection
		validSortFields := map[string]bool{
			"created_at":   true,
			"total_amount": true,
			"order_number": true,
			"status":       true,
			"updated_at":   true,
		}
		if validSortFields[filter.SortBy] {
			sortBy = filter.SortBy
		}
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT * FROM orders %s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`, whereClause, sortBy, sortOrder, argIndex, argIndex+1)

	args = append(args, limit, offset)

	var orders []*domain.Order
	err = r.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, total, nil
}

// Order Items

// CreateOrderItems creates order items
func (r *orderRepository) CreateOrderItems(ctx context.Context, items []*domain.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	query := `
		INSERT INTO order_items (order_id, product_id, product_variant_id, product_name, product_sku,
			quantity, unit_price, total_price, discount_amount, is_digital, requires_shipping)
		VALUES (:order_id, :product_id, :product_variant_id, :product_name, :product_sku,
			:quantity, :unit_price, :total_price, :discount_amount, :is_digital, :requires_shipping)`

	_, err := r.db.NamedExecContext(ctx, query, items)
	if err != nil {
		return fmt.Errorf("failed to create order items: %w", err)
	}

	return nil
}

// GetOrderItems retrieves order items
func (r *orderRepository) GetOrderItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error) {
	query := `SELECT * FROM order_items WHERE order_id = $1 ORDER BY id`

	var items []*domain.OrderItem
	err := r.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return items, nil
}

// UpdateOrderItem updates an order item
func (r *orderRepository) UpdateOrderItem(ctx context.Context, id int64, item *domain.OrderItem) error {
	query := `
		UPDATE order_items SET 
			quantity = :quantity, unit_price = :unit_price, total_price = :total_price,
			discount_amount = :discount_amount
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update order item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order item with ID %d not found", id)
	}

	return nil
}

// DeleteOrderItem deletes an order item
func (r *orderRepository) DeleteOrderItem(ctx context.Context, id int64) error {
	query := `DELETE FROM order_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order item with ID %d not found", id)
	}

	return nil
}

// Order Addresses

// CreateOrderAddresses creates order addresses
func (r *orderRepository) CreateOrderAddresses(ctx context.Context, addresses []*domain.OrderAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	query := `
		INSERT INTO order_addresses (order_id, type, first_name, last_name, company, address1, address2,
			city, state, country, postal_code, phone, email)
		VALUES (:order_id, :type, :first_name, :last_name, :company, :address1, :address2,
			:city, :state, :country, :postal_code, :phone, :email)`

	_, err := r.db.NamedExecContext(ctx, query, addresses)
	if err != nil {
		return fmt.Errorf("failed to create order addresses: %w", err)
	}

	return nil
}

// GetOrderAddresses retrieves order addresses
func (r *orderRepository) GetOrderAddresses(ctx context.Context, orderID int64) ([]*domain.OrderAddress, error) {
	query := `SELECT * FROM order_addresses WHERE order_id = $1 ORDER BY type`

	var addresses []*domain.OrderAddress
	err := r.db.SelectContext(ctx, &addresses, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order addresses: %w", err)
	}

	return addresses, nil
}

// UpdateOrderAddress updates an order address
func (r *orderRepository) UpdateOrderAddress(ctx context.Context, id int64, address *domain.OrderAddress) error {
	query := `
		UPDATE order_addresses SET 
			first_name = :first_name, last_name = :last_name, company = :company,
			address1 = :address1, address2 = :address2, city = :city, state = :state,
			country = :country, postal_code = :postal_code, phone = :phone, email = :email
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, address)
	if err != nil {
		return fmt.Errorf("failed to update order address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order address with ID %d not found", id)
	}

	return nil
}

// DeleteOrderAddress deletes an order address
func (r *orderRepository) DeleteOrderAddress(ctx context.Context, id int64) error {
	query := `DELETE FROM order_addresses WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order address: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order address with ID %d not found", id)
	}

	return nil
}

// Order Status History

// CreateOrderStatusHistory creates order status history
func (r *orderRepository) CreateOrderStatusHistory(ctx context.Context, history *domain.OrderStatusHistory) error {
	query := `
		INSERT INTO order_status_history (order_id, status, previous_status, notes, created_by, created_at)
		VALUES (:order_id, :status, :previous_status, :notes, :created_by, :created_at)`

	history.CreatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, history)
	if err != nil {
		return fmt.Errorf("failed to create order status history: %w", err)
	}

	return nil
}

// GetOrderStatusHistory retrieves order status history
func (r *orderRepository) GetOrderStatusHistory(ctx context.Context, orderID int64) ([]*domain.OrderStatusHistory, error) {
	query := `SELECT * FROM order_status_history WHERE order_id = $1 ORDER BY created_at ASC`

	var history []*domain.OrderStatusHistory
	err := r.db.SelectContext(ctx, &history, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status history: %w", err)
	}

	return history, nil
}

// Order Fulfillment

// CreateOrderFulfillment creates order fulfillment
func (r *orderRepository) CreateOrderFulfillment(ctx context.Context, fulfillment *domain.OrderFulfillment) error {
	query := `
		INSERT INTO order_fulfillment (order_id, tracking_number, carrier, service, status,
			shipped_at, delivered_at, estimated_delivery, notes, created_at, updated_at)
		VALUES (:order_id, :tracking_number, :carrier, :service, :status,
			:shipped_at, :delivered_at, :estimated_delivery, :notes, :created_at, :updated_at)`

	fulfillment.CreatedAt = time.Now()
	fulfillment.UpdatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, fulfillment)
	if err != nil {
		return fmt.Errorf("failed to create order fulfillment: %w", err)
	}

	return nil
}

// GetOrderFulfillment retrieves order fulfillment
func (r *orderRepository) GetOrderFulfillment(ctx context.Context, orderID int64) (*domain.OrderFulfillment, error) {
	query := `SELECT * FROM order_fulfillment WHERE order_id = $1`

	var fulfillment domain.OrderFulfillment
	err := r.db.GetContext(ctx, &fulfillment, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No fulfillment record
		}
		return nil, fmt.Errorf("failed to get order fulfillment: %w", err)
	}

	return &fulfillment, nil
}

// UpdateOrderFulfillment updates order fulfillment
func (r *orderRepository) UpdateOrderFulfillment(ctx context.Context, id int64, fulfillment *domain.OrderFulfillment) error {
	query := `
		UPDATE order_fulfillment SET 
			tracking_number = :tracking_number, carrier = :carrier, service = :service, status = :status,
			shipped_at = :shipped_at, delivered_at = :delivered_at, estimated_delivery = :estimated_delivery,
			notes = :notes, updated_at = :updated_at
		WHERE id = :id`

	fulfillment.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, fulfillment)
	if err != nil {
		return fmt.Errorf("failed to update order fulfillment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order fulfillment with ID %d not found", id)
	}

	return nil
}

// Order Refunds

// CreateOrderRefund creates order refund
func (r *orderRepository) CreateOrderRefund(ctx context.Context, refund *domain.OrderRefund) error {
	query := `
		INSERT INTO order_refunds (order_id, amount, reason, status, processed_at, created_by, created_at)
		VALUES (:order_id, :amount, :reason, :status, :processed_at, :created_by, :created_at)`

	refund.CreatedAt = time.Now()

	_, err := r.db.NamedExecContext(ctx, query, refund)
	if err != nil {
		return fmt.Errorf("failed to create order refund: %w", err)
	}

	return nil
}

// GetOrderRefunds retrieves order refunds
func (r *orderRepository) GetOrderRefunds(ctx context.Context, orderID int64) ([]*domain.OrderRefund, error) {
	query := `SELECT * FROM order_refunds WHERE order_id = $1 ORDER BY created_at DESC`

	var refunds []*domain.OrderRefund
	err := r.db.SelectContext(ctx, &refunds, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order refunds: %w", err)
	}

	return refunds, nil
}

// UpdateOrderRefund updates order refund
func (r *orderRepository) UpdateOrderRefund(ctx context.Context, id int64, refund *domain.OrderRefund) error {
	query := `
		UPDATE order_refunds SET 
			amount = :amount, reason = :reason, status = :status, processed_at = :processed_at
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, refund)
	if err != nil {
		return fmt.Errorf("failed to update order refund: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order refund with ID %d not found", id)
	}

	return nil
}

// Order Analytics

// GetOrderAnalytics retrieves order analytics
func (r *orderRepository) GetOrderAnalytics(ctx context.Context) (*domain.OrderAnalytics, error) {
	// Get total orders count
	var totalOrders int64
	err := r.db.GetContext(ctx, &totalOrders, `SELECT COUNT(*) FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders count: %w", err)
	}

	// Get orders by status
	var pendingOrders, confirmedOrders, processingOrders, shippedOrders, deliveredOrders, cancelledOrders int64

	statusQuery := `
		SELECT 
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as confirmed,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing,
			COUNT(CASE WHEN status = 'shipped' THEN 1 END) as shipped,
			COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled
		FROM orders`

	var counts struct {
		Pending    int64 `db:"pending"`
		Confirmed  int64 `db:"confirmed"`
		Processing int64 `db:"processing"`
		Shipped    int64 `db:"shipped"`
		Delivered  int64 `db:"delivered"`
		Cancelled  int64 `db:"cancelled"`
	}

	err = r.db.GetContext(ctx, &counts, statusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}

	pendingOrders = counts.Pending
	confirmedOrders = counts.Confirmed
	processingOrders = counts.Processing
	shippedOrders = counts.Shipped
	deliveredOrders = counts.Delivered
	cancelledOrders = counts.Cancelled

	// Get total revenue
	var totalRevenue float64
	err = r.db.GetContext(ctx, &totalRevenue, `
		SELECT COALESCE(SUM(total_amount), 0) FROM orders 
		WHERE status IN ('confirmed', 'processing', 'shipped', 'delivered')`)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}

	// Calculate average order value
	var averageOrderValue float64
	if totalOrders > 0 {
		averageOrderValue = totalRevenue / float64(totalOrders)
	}

	// Get time-based metrics
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var newOrdersToday, newOrdersThisWeek, newOrdersThisMonth int64

	err = r.db.GetContext(ctx, &newOrdersToday, `SELECT COUNT(*) FROM orders WHERE created_at >= $1`, today)
	if err != nil {
		newOrdersToday = 0
	}

	err = r.db.GetContext(ctx, &newOrdersThisWeek, `SELECT COUNT(*) FROM orders WHERE created_at >= $1`, weekStart)
	if err != nil {
		newOrdersThisWeek = 0
	}

	err = r.db.GetContext(ctx, &newOrdersThisMonth, `SELECT COUNT(*) FROM orders WHERE created_at >= $1`, monthStart)
	if err != nil {
		newOrdersThisMonth = 0
	}

	// Calculate conversion rate (placeholder - would need visitor data)
	var conversionRate float64
	if totalOrders > 0 {
		conversionRate = 5.0 // Placeholder
	}

	analytics := &domain.OrderAnalytics{
		TotalOrders:        totalOrders,
		PendingOrders:      pendingOrders,
		ConfirmedOrders:    confirmedOrders,
		ProcessingOrders:   processingOrders,
		ShippedOrders:      shippedOrders,
		DeliveredOrders:    deliveredOrders,
		CancelledOrders:    cancelledOrders,
		TotalRevenue:       totalRevenue,
		AverageOrderValue:  averageOrderValue,
		ConversionRate:     conversionRate,
		NewOrdersToday:     newOrdersToday,
		NewOrdersThisWeek:  newOrdersThisWeek,
		NewOrdersThisMonth: newOrdersThisMonth,
	}

	return analytics, nil
}

// GetOrderAnalyticsByDateRange retrieves order analytics by date range
func (r *orderRepository) GetOrderAnalyticsByDateRange(ctx context.Context, startDate, endDate time.Time) (*domain.OrderAnalytics, error) {
	// Get total orders count in range
	var totalOrders int64
	err := r.db.GetContext(ctx, &totalOrders, `SELECT COUNT(*) FROM orders WHERE created_at BETWEEN $1 AND $2`, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get total orders count by range: %w", err)
	}

	// Get orders by status in range
	var pendingOrders, confirmedOrders, processingOrders, shippedOrders, deliveredOrders, cancelledOrders int64

	statusQuery := `
		SELECT 
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as confirmed,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing,
			COUNT(CASE WHEN status = 'shipped' THEN 1 END) as shipped,
			COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled
		FROM orders
		WHERE created_at BETWEEN $1 AND $2`

	var counts struct {
		Pending    int64 `db:"pending"`
		Confirmed  int64 `db:"confirmed"`
		Processing int64 `db:"processing"`
		Shipped    int64 `db:"shipped"`
		Delivered  int64 `db:"delivered"`
		Cancelled  int64 `db:"cancelled"`
	}

	err = r.db.GetContext(ctx, &counts, statusQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status in range: %w", err)
	}

	pendingOrders = counts.Pending
	confirmedOrders = counts.Confirmed
	processingOrders = counts.Processing
	shippedOrders = counts.Shipped
	deliveredOrders = counts.Delivered
	cancelledOrders = counts.Cancelled

	// Get total revenue in range
	var totalRevenue float64
	err = r.db.GetContext(ctx, &totalRevenue, `
		SELECT COALESCE(SUM(total_amount), 0) FROM orders 
		WHERE status IN ('confirmed', 'processing', 'shipped', 'delivered')
		AND created_at BETWEEN $1 AND $2`, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue in range: %w", err)
	}

	// Calculate average order value
	var averageOrderValue float64
	if totalOrders > 0 {
		averageOrderValue = totalRevenue / float64(totalOrders)
	}

	today := time.Now().Truncate(24 * time.Hour)

	var newOrdersToday int64

	err = r.db.GetContext(ctx, &newOrdersToday, `SELECT COUNT(*) FROM orders WHERE created_at >= $1`, today)
	if err != nil {
		newOrdersToday = 0
	}

	// Calculate conversion rate placeholder
	conversionRate := 5.0

	return &domain.OrderAnalytics{
		TotalOrders:        totalOrders,
		PendingOrders:      pendingOrders,
		ConfirmedOrders:    confirmedOrders,
		ProcessingOrders:   processingOrders,
		ShippedOrders:      shippedOrders,
		DeliveredOrders:    deliveredOrders,
		CancelledOrders:    cancelledOrders,
		TotalRevenue:       totalRevenue,
		AverageOrderValue:  averageOrderValue,
		ConversionRate:     conversionRate,
		NewOrdersToday:     newOrdersToday,
		NewOrdersThisWeek:  0, // Omitted for brevity in range calc
		NewOrdersThisMonth: 0, // Omitted for brevity in range calc
	}, nil
}

// GetTopSellingProducts retrieves top selling products
func (r *orderRepository) GetTopSellingProducts(ctx context.Context, limit int) ([]*domain.ProductOrderStats, error) {
	query := `
		SELECT 
			oi.product_id,
			p.name as product_name,
			p.sku,
			SUM(oi.quantity) as total_quantity,
			SUM(oi.total_price) as total_revenue,
			COUNT(DISTINCT oi.order_id) as order_count,
			AVG(oi.unit_price) as average_price
		FROM order_items oi
		INNER JOIN products p ON oi.product_id = p.id
		INNER JOIN orders o ON oi.order_id = o.id
		WHERE o.status IN ('confirmed', 'processing', 'shipped', 'delivered')
		GROUP BY oi.product_id, p.name, p.sku
		ORDER BY total_quantity DESC
		LIMIT $1`

	var stats []*domain.ProductOrderStats
	err := r.db.SelectContext(ctx, &stats, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top selling products: %w", err)
	}

	return stats, nil
}

// GetOrderConversionRate calculates order conversion rate
func (r *orderRepository) GetOrderConversionRate(ctx context.Context) (float64, error) {
	// This would need visitor data to calculate actual conversion rate
	// For now, returning a placeholder
	return 5.0, nil
}
