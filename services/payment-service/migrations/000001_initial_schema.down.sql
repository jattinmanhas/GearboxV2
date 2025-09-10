-- Drop triggers
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_payment_gateways_updated_at ON payment_gateways;
DROP TRIGGER IF EXISTS update_payment_methods_updated_at ON payment_methods;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS payment_webhooks;
DROP TABLE IF EXISTS payment_refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_gateways;
DROP TABLE IF EXISTS payment_methods;
