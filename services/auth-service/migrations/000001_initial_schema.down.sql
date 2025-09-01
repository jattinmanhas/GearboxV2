-- =====================================================
-- Migration: 000001_initial_schema.down.sql
-- Description: Rollback initial authentication schema
-- =====================================================

-- Drop triggers first
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON refresh_tokens;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
