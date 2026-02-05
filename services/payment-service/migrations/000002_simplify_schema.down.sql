-- Rollback: Restore removed columns and data

-- Step 1: Re-add the 'config' column to payment_gateways
ALTER TABLE payment_gateways ADD COLUMN IF NOT EXISTS config TEXT DEFAULT '{}';

-- Step 2: Re-add webhook_url column
ALTER TABLE payment_gateways ADD COLUMN IF NOT EXISTS webhook_url TEXT;

-- Step 3: Re-insert PayPal and Razorpay gateways
INSERT INTO payment_gateways (name, code, is_active, is_test_mode, config, sort_order)
VALUES 
    ('PayPal', 'paypal', false, true, '{}', 2),
    ('Razorpay', 'razorpay', false, true, '{}', 3)
ON CONFLICT (code) DO NOTHING;

-- Step 4: Re-insert removed payment methods
INSERT INTO payment_methods (name, code, type, is_active, sort_order, description)
VALUES 
    ('PayPal', 'paypal', 'paypal', false, 3, 'Pay with PayPal')
ON CONFLICT (code) DO NOTHING;
