-- Create inventory tables
-- This migration creates tables for inventory management, stock movements, reservations, and alerts

-- Create inventory table
CREATE TABLE inventory (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL, -- References products(id) ON DELETE CASCADE
    product_variant_id BIGINT, -- References product_variants(id) ON DELETE CASCADE
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    available_quantity INTEGER NOT NULL DEFAULT 0 CHECK (available_quantity >= 0),
    min_stock_level INTEGER NOT NULL DEFAULT 0 CHECK (min_stock_level >= 0),
    max_stock_level INTEGER CHECK (max_stock_level >= 0),
    reorder_point INTEGER NOT NULL DEFAULT 0 CHECK (reorder_point >= 0),
    last_restocked TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Ensure available_quantity is calculated correctly
    CONSTRAINT check_available_quantity CHECK (available_quantity = quantity - reserved_quantity),
    
    -- Unique constraint for product/variant combination
    UNIQUE(product_id, product_variant_id)
);

-- Create inventory_movements table
CREATE TABLE inventory_movements (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL, -- References products(id) ON DELETE CASCADE
    product_variant_id BIGINT, -- References product_variants(id) ON DELETE CASCADE
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('in', 'out', 'adjustment', 'transfer')),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    previous_quantity INTEGER NOT NULL DEFAULT 0 CHECK (previous_quantity >= 0),
    new_quantity INTEGER NOT NULL DEFAULT 0 CHECK (new_quantity >= 0),
    reference VARCHAR(255),
    reference_type VARCHAR(50),
    reason VARCHAR(255),
    notes TEXT,
    created_by BIGINT, -- References users(id) ON DELETE SET NULL
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create stock_reservations table
CREATE TABLE stock_reservations (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL, -- References products(id) ON DELETE CASCADE
    product_variant_id BIGINT, -- References product_variants(id) ON DELETE CASCADE
    order_id BIGINT NOT NULL, -- References orders(id) ON DELETE CASCADE
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create inventory_alerts table
CREATE TABLE inventory_alerts (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL, -- References products(id) ON DELETE CASCADE
    product_variant_id BIGINT, -- References product_variants(id) ON DELETE CASCADE
    alert_type VARCHAR(20) NOT NULL CHECK (alert_type IN ('low_stock', 'out_of_stock', 'reorder_point')),
    current_quantity INTEGER NOT NULL CHECK (current_quantity >= 0),
    threshold_quantity INTEGER NOT NULL CHECK (threshold_quantity >= 0),
    is_resolved BOOLEAN NOT NULL DEFAULT FALSE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_inventory_product_id ON inventory(product_id);
CREATE INDEX idx_inventory_product_variant_id ON inventory(product_variant_id);
CREATE INDEX idx_inventory_available_quantity ON inventory(available_quantity);
CREATE INDEX idx_inventory_low_stock ON inventory(available_quantity, reorder_point) WHERE available_quantity <= reorder_point;

CREATE INDEX idx_inventory_movements_product_id ON inventory_movements(product_id);
CREATE INDEX idx_inventory_movements_product_variant_id ON inventory_movements(product_variant_id);
CREATE INDEX idx_inventory_movements_movement_type ON inventory_movements(movement_type);
CREATE INDEX idx_inventory_movements_created_at ON inventory_movements(created_at);
CREATE INDEX idx_inventory_movements_reference ON inventory_movements(reference);

CREATE INDEX idx_stock_reservations_product_id ON stock_reservations(product_id);
CREATE INDEX idx_stock_reservations_product_variant_id ON stock_reservations(product_variant_id);
CREATE INDEX idx_stock_reservations_order_id ON stock_reservations(order_id);
CREATE INDEX idx_stock_reservations_expires_at ON stock_reservations(expires_at);

CREATE INDEX idx_inventory_alerts_product_id ON inventory_alerts(product_id);
CREATE INDEX idx_inventory_alerts_product_variant_id ON inventory_alerts(product_variant_id);
CREATE INDEX idx_inventory_alerts_alert_type ON inventory_alerts(alert_type);
CREATE INDEX idx_inventory_alerts_is_resolved ON inventory_alerts(is_resolved);
CREATE INDEX idx_inventory_alerts_created_at ON inventory_alerts(created_at);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_inventory_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_inventory_updated_at
    BEFORE UPDATE ON inventory
    FOR EACH ROW
    EXECUTE FUNCTION update_inventory_updated_at();

-- Create trigger to automatically update available_quantity
CREATE OR REPLACE FUNCTION update_inventory_available_quantity()
RETURNS TRIGGER AS $$
BEGIN
    NEW.available_quantity = NEW.quantity - NEW.reserved_quantity;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_inventory_available_quantity
    BEFORE INSERT OR UPDATE ON inventory
    FOR EACH ROW
    EXECUTE FUNCTION update_inventory_available_quantity();

-- Create function to cleanup expired reservations
CREATE OR REPLACE FUNCTION cleanup_expired_reservations()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM stock_reservations WHERE expires_at < NOW();
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Add comments for documentation
COMMENT ON TABLE inventory IS 'Tracks inventory levels for products and variants';
COMMENT ON TABLE inventory_movements IS 'Records all stock movements (in/out/adjustments)';
COMMENT ON TABLE stock_reservations IS 'Tracks reserved stock for pending orders';
COMMENT ON TABLE inventory_alerts IS 'Manages low stock and out of stock alerts';

COMMENT ON COLUMN inventory.quantity IS 'Total quantity in stock';
COMMENT ON COLUMN inventory.reserved_quantity IS 'Quantity reserved for pending orders';
COMMENT ON COLUMN inventory.available_quantity IS 'Available quantity (quantity - reserved_quantity)';
COMMENT ON COLUMN inventory.min_stock_level IS 'Minimum stock level before reorder';
COMMENT ON COLUMN inventory.max_stock_level IS 'Maximum stock level for reordering';
COMMENT ON COLUMN inventory.reorder_point IS 'Stock level at which to reorder';

COMMENT ON COLUMN inventory_movements.movement_type IS 'Type of movement: in, out, adjustment, transfer';
COMMENT ON COLUMN inventory_movements.reference IS 'Reference to order, purchase order, etc.';
COMMENT ON COLUMN inventory_movements.reference_type IS 'Type of reference: order, purchase, adjustment, etc.';

COMMENT ON COLUMN stock_reservations.expires_at IS 'When the reservation expires and stock is released';
COMMENT ON COLUMN stock_reservations.order_id IS 'Order this stock is reserved for';

COMMENT ON COLUMN inventory_alerts.alert_type IS 'Type of alert: low_stock, out_of_stock, reorder_point';
COMMENT ON COLUMN inventory_alerts.is_resolved IS 'Whether the alert has been resolved';
