-- =====================================================
-- Migration: 000001_initial_schema.up.sql
-- Description: Create initial authentication schema
-- Tables: users, roles, refresh_tokens
-- =====================================================

-- Enable UUID extension for future use
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create roles table first (needed for foreign key)
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100), -- Made optional
    avatar TEXT,
    gender VARCHAR(20) CHECK (gender = '' OR gender IN ('male', 'female', 'other', 'prefer_not_to_say') OR gender IS NULL), -- Fixed constraint to allow empty strings
    date_of_birth DATE, -- Made optional
    role_id BIGINT DEFAULT 1, -- Default to user role
    is_active BOOLEAN DEFAULT TRUE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create refresh_tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address INET,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active) WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_roles_is_active ON roles(is_active) WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(refresh_token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_is_revoked ON refresh_tokens(is_revoked);
CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_user_token ON refresh_tokens(user_id, refresh_token);

-- Insert default roles
INSERT INTO roles (id, name, description, is_active) VALUES
    (1, 'user', 'Basic authenticated user with limited access', TRUE),
    (2, 'editor', 'Content editor with create/edit/moderate permissions', TRUE),
    (3, 'admin', 'Full system administrator with complete access', TRUE)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    is_active = EXCLUDED.is_active;

-- Set default role for existing users (if any)
UPDATE users SET role_id = 1 WHERE role_id IS NULL;

-- Add foreign key constraint for role_id
ALTER TABLE users ADD CONSTRAINT fk_users_role_id 
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT;

-- Create function to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at 
    BEFORE UPDATE ON roles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refresh_tokens_updated_at 
    BEFORE UPDATE ON refresh_tokens 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
