-- Migration: Simplify payment service schema
-- Removes unused fields and gateways, keeps only Stripe

-- Step 1: Remove the 'config' column from payment_gateways (API keys should only be in env vars)
ALTER TABLE payment_gateways DROP COLUMN IF EXISTS config;

-- Step 2: Remove webhook_url from payment_gateways (configured in Stripe dashboard)
ALTER TABLE payment_gateways DROP COLUMN IF EXISTS webhook_url;

-- Step 3: Remove PayPal and Razorpay gateways (keeping only Stripe)
DELETE FROM payment_gateways WHERE code IN ('paypal', 'razorpay');

-- Step 4: Remove unused payment methods (keeping only card and stripe-related)
DELETE FROM payment_methods WHERE code NOT IN ('card', 'stripe');

-- Step 5: Drop the payment_webhooks table if not needed (raw webhook logs)
-- Uncomment if you want to remove webhook logging entirely:
-- DROP TABLE IF EXISTS payment_webhooks;
