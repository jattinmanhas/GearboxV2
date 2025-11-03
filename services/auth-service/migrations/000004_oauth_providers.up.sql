-- =====================================================
-- Migration: 000004_oauth_providers.up.sql
-- Description: Add OAuth providers support
-- Tables: oauth_providers
-- Changes: Make users.password nullable
-- =====================================================

-- Modify users table to make password nullable for OAuth users
ALTER TABLE users ALTER COLUMN password DROP NOT NULL;

-- Create oauth_providers table
CREATE TABLE IF NOT EXISTS oauth_providers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL CHECK (provider IN ('google', 'facebook', 'github')),
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one provider per user (user can link multiple providers, but only once per provider)
    UNIQUE(user_id, provider),
    -- Ensure provider user_id is unique per provider
    UNIQUE(provider, provider_user_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_oauth_providers_user_id ON oauth_providers(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_provider ON oauth_providers(provider);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_provider_user_id ON oauth_providers(provider, provider_user_id);

-- Create trigger for updated_at
CREATE TRIGGER update_oauth_providers_updated_at 
    BEFORE UPDATE ON oauth_providers 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

