-- Migration to remove quantity columns from products and product_variants tables
-- This makes inventory table the single source of truth for quantity

-- Remove quantity column from products table
ALTER TABLE products DROP COLUMN IF EXISTS quantity;

-- Remove quantity column from product_variants table  
ALTER TABLE product_variants DROP COLUMN IF EXISTS quantity;

-- Add comments to clarify the new structure
COMMENT ON TABLE inventory IS 'Single source of truth for product quantities and stock management';
COMMENT ON COLUMN inventory.quantity IS 'Total quantity in stock - this is the authoritative quantity for all products and variants';
COMMENT ON COLUMN inventory.available_quantity IS 'Available quantity for sale (quantity - reserved_quantity)';
COMMENT ON COLUMN inventory.reserved_quantity IS 'Quantity reserved for pending orders';

-- Update any existing indexes that might reference the removed columns
-- (PostgreSQL will automatically drop indexes on dropped columns)
