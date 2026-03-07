-- Restore deprecated product fields (rollback)
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS taxable BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS track_quantity BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS min_quantity INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS max_quantity INTEGER;
