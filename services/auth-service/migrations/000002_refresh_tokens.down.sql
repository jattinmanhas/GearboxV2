-- Drop refresh_tokens table and related objects
DROP TRIGGER IF EXISTS update_refresh_tokens_updated_at ON refresh_tokens;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS refresh_tokens;
