-- =====================================================
-- Migration: 000004_oauth_providers.down.sql
-- Description: Rollback OAuth providers support
-- =====================================================

-- Drop oauth_providers table
DROP TABLE IF EXISTS oauth_providers;

-- Make users.password NOT NULL again (be careful with this in production)
-- Note: This will fail if there are users without passwords
ALTER TABLE users ALTER COLUMN password SET NOT NULL;

