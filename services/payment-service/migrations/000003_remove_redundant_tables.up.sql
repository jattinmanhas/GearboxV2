-- Migration: Remove redundant payment tables
-- Drops payment_gateways and payment_methods as the service is now Stripe-only

-- Step 1: Drop foreign key from payments table
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_payment_method_id_fkey;

-- Step 2: Drop the redundant tables
DROP TABLE IF EXISTS payment_gateways;
DROP TABLE IF EXISTS payment_methods;

-- Step 3: Modify the payments table
-- Change payment_method_id (BIGINT) to payment_method (VARCHAR)
-- We rename it first to avoid conflicts, then add the new column, migration data, then drop old
ALTER TABLE payments RENAME COLUMN payment_method_id TO old_payment_method_id;
ALTER TABLE payments ADD COLUMN payment_method VARCHAR(50) NOT NULL DEFAULT 'card';

-- Step 4: Drop the old column
ALTER TABLE payments DROP COLUMN old_payment_method_id;

-- Step 5: Update gateway_id to be more consistent (already VARCHAR)
-- No changes needed to gateway_id column type, but we ensure it's used as a code now.
