-- =====================================================
-- Migration: 000006_password_reset_tokens.down.sql
-- Description: Drop password reset tokens table
-- =====================================================

-- Drop trigger
DROP TRIGGER IF EXISTS update_password_reset_tokens_updated_at ON password_reset_tokens;

-- Drop indexes
DROP INDEX IF EXISTS idx_password_reset_tokens_used;
DROP INDEX IF EXISTS idx_password_reset_tokens_expires_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;

-- Drop table
DROP TABLE IF EXISTS password_reset_tokens;

