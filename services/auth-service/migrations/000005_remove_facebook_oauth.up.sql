-- =====================================================
-- Migration: 000005_remove_facebook_oauth.up.sql
-- Description: Remove Facebook from OAuth providers constraint
-- =====================================================

-- Remove Facebook from the provider CHECK constraint
-- First, drop the existing constraint
ALTER TABLE oauth_providers DROP CONSTRAINT IF EXISTS oauth_providers_provider_check;

-- Add the new constraint without Facebook
ALTER TABLE oauth_providers ADD CONSTRAINT oauth_providers_provider_check 
    CHECK (provider IN ('google', 'github'));

