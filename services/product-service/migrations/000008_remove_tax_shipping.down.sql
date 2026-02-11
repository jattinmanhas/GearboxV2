-- Restore tax_amount and shipping_amount to orders table
ALTER TABLE orders ADD COLUMN tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN shipping_amount DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Restore tax_amount to order_items table
ALTER TABLE order_items ADD COLUMN tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Restore shipping_amount to cart_shipping table
ALTER TABLE cart_shipping ADD COLUMN shipping_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE cart_shipping ADD CONSTRAINT check_cart_shipping_positive_amount CHECK (shipping_amount >= 0);
