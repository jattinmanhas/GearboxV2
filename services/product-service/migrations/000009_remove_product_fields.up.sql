-- Remove deprecated product fields
-- taxable: removed as part of tax removal update
-- track_quantity: stock tracking is handled by the inventory table
-- min_quantity / max_quantity: these are already in the inventory table (min_stock_level / max_stock_level)
ALTER TABLE products
    DROP COLUMN IF EXISTS taxable,
    DROP COLUMN IF EXISTS track_quantity,
    DROP COLUMN IF EXISTS min_quantity,
    DROP COLUMN IF EXISTS max_quantity;
