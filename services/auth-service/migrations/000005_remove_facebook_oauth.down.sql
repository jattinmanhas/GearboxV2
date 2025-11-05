-- =====================================================
-- Migration: 000005_remove_facebook_oauth.down.sql
-- Description: Rollback - Restore Facebook in OAuth providers constraint
-- =====================================================

-- Drop the constraint
ALTER TABLE oauth_providers DROP CONSTRAINT IF EXISTS oauth_providers_provider_check;

-- Restore the original constraint with Facebook
ALTER TABLE oauth_providers ADD CONSTRAINT oauth_providers_provider_check 
    CHECK (provider IN ('google', 'facebook', 'github'));

