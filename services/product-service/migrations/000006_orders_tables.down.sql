-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS order_refunds;
DROP TABLE IF EXISTS order_fulfillment;
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_addresses;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- Drop function
DROP FUNCTION IF EXISTS generate_order_number();
