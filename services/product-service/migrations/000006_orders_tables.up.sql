-- Create orders table
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded')),
    payment_status VARCHAR(30) NOT NULL DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded', 'partially_refunded')),
    fulfillment_status VARCHAR(20) NOT NULL DEFAULT 'unfulfilled' CHECK (fulfillment_status IN ('unfulfilled', 'partial', 'fulfilled')),
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    shipping_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    notes TEXT,
    internal_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    shipped_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE
);

-- Create order_items table
CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_variant_id BIGINT REFERENCES product_variants(id) ON DELETE RESTRICT,
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10,2) NOT NULL CHECK (unit_price >= 0),
    total_price DECIMAL(10,2) NOT NULL CHECK (total_price >= 0),
    tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_digital BOOLEAN NOT NULL DEFAULT false,
    requires_shipping BOOLEAN NOT NULL DEFAULT true
);

-- Create order_addresses table
CREATE TABLE order_addresses (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('shipping', 'billing')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    company VARCHAR(100),
    address1 VARCHAR(255) NOT NULL,
    address2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255)
);

-- Create order_status_history table
CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    previous_status VARCHAR(20),
    notes TEXT,
    created_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create order_fulfillment table
CREATE TABLE order_fulfillment (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    tracking_number VARCHAR(100),
    carrier VARCHAR(100),
    service VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'shipped', 'delivered', 'failed')),
    shipped_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    estimated_delivery TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create order_refunds table
CREATE TABLE order_refunds (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processed', 'failed')),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for orders table
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_fulfillment_status ON orders(fulfillment_status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_order_number ON orders(order_number);

-- Create indexes for order_items table
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_items_product_variant_id ON order_items(product_variant_id);

-- Create indexes for order_addresses table
CREATE INDEX idx_order_addresses_order_id ON order_addresses(order_id);
CREATE INDEX idx_order_addresses_type ON order_addresses(type);

-- Create indexes for order_status_history table
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_order_status_history_status ON order_status_history(status);
CREATE INDEX idx_order_status_history_created_at ON order_status_history(created_at);

-- Create indexes for order_fulfillment table
CREATE INDEX idx_order_fulfillment_order_id ON order_fulfillment(order_id);
CREATE INDEX idx_order_fulfillment_tracking_number ON order_fulfillment(tracking_number);
CREATE INDEX idx_order_fulfillment_status ON order_fulfillment(status);

-- Create indexes for order_refunds table
CREATE INDEX idx_order_refunds_order_id ON order_refunds(order_id);
CREATE INDEX idx_order_refunds_status ON order_refunds(status);
CREATE INDEX idx_order_refunds_created_at ON order_refunds(created_at);

-- Create trigger to update updated_at timestamp for orders
CREATE TRIGGER update_orders_updated_at 
    BEFORE UPDATE ON orders 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create trigger to update updated_at timestamp for order_fulfillment
CREATE TRIGGER update_order_fulfillment_updated_at 
    BEFORE UPDATE ON order_fulfillment 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to generate order numbers
CREATE OR REPLACE FUNCTION generate_order_number() RETURNS TEXT AS $$
DECLARE
    new_number TEXT;
    counter INTEGER;
BEGIN
    -- Get current date in YYYYMMDD format
    new_number := TO_CHAR(NOW(), 'YYYYMMDD');
    
    -- Get count of orders created today
    SELECT COUNT(*) + 1 INTO counter 
    FROM orders 
    WHERE DATE(created_at) = CURRENT_DATE;
    
    -- Append counter with leading zeros
    new_number := new_number || LPAD(counter::TEXT, 4, '0');
    
    RETURN new_number;
END;
$$ LANGUAGE plpgsql;

-- Add comments
COMMENT ON TABLE orders IS 'Customer orders and their status';
COMMENT ON TABLE order_items IS 'Individual items within orders';
COMMENT ON TABLE order_addresses IS 'Shipping and billing addresses for orders';
COMMENT ON TABLE order_status_history IS 'History of order status changes';
COMMENT ON TABLE order_fulfillment IS 'Order fulfillment and tracking information';
COMMENT ON TABLE order_refunds IS 'Order refund transactions';

COMMENT ON COLUMN orders.status IS 'Order status: pending, confirmed, processing, shipped, delivered, cancelled, refunded';
COMMENT ON COLUMN orders.payment_status IS 'Payment status: pending, paid, failed, refunded, partially_refunded';
COMMENT ON COLUMN orders.fulfillment_status IS 'Fulfillment status: unfulfilled, partial, fulfilled';
COMMENT ON COLUMN order_addresses.type IS 'Address type: shipping or billing';
COMMENT ON COLUMN order_fulfillment.status IS 'Fulfillment status: pending, shipped, delivered, failed';
COMMENT ON COLUMN order_refunds.status IS 'Refund status: pending, processed, failed';
