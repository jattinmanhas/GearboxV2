-- Rollback migration to restore quantity columns

-- Add quantity column back to products table
ALTER TABLE products ADD COLUMN quantity INTEGER NOT NULL DEFAULT 0;

-- Add quantity column back to product_variants table
ALTER TABLE product_variants ADD COLUMN quantity INTEGER NOT NULL DEFAULT 0;

-- Note: This rollback will not restore the original data values
-- You would need to manually populate these columns from inventory data if needed
