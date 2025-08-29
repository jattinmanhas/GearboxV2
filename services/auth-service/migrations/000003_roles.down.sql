-- Remove role_id column from users table
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- Drop roles table
DROP TABLE IF EXISTS roles;
