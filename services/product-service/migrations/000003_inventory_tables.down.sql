-- Drop inventory tables migration
-- This migration removes all inventory-related tables and functions

-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_update_inventory_updated_at ON inventory;
DROP TRIGGER IF EXISTS trigger_update_inventory_available_quantity ON inventory;

-- Drop functions
DROP FUNCTION IF EXISTS update_inventory_updated_at();
DROP FUNCTION IF EXISTS update_inventory_available_quantity();
DROP FUNCTION IF EXISTS cleanup_expired_reservations();

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS inventory_alerts;
DROP TABLE IF EXISTS stock_reservations;
DROP TABLE IF EXISTS inventory_movements;
DROP TABLE IF EXISTS inventory;
