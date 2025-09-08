-- =====================================================
-- Migration: 000002_user_addresses.down.sql
-- Description: Remove user address management system
-- =====================================================

-- Drop triggers first
DROP TRIGGER IF EXISTS update_user_addresses_updated_at ON user_addresses;
DROP TRIGGER IF EXISTS update_user_phone_numbers_updated_at ON user_phone_numbers;

-- Drop tables
DROP TABLE IF EXISTS user_phone_numbers;
DROP TABLE IF EXISTS user_addresses;
