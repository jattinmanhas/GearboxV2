-- =====================================================
-- Migration: 000002_user_addresses.up.sql
-- Description: Add user address management system
-- Tables: user_addresses, user_phone_numbers
-- =====================================================

-- Create user_addresses table
CREATE TABLE IF NOT EXISTS user_addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL DEFAULT 'home' CHECK (address_type IN ('home', 'work', 'billing', 'shipping', 'other')),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    company VARCHAR(100),
    address_line_1 VARCHAR(255) NOT NULL,
    address_line_2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    is_verified BOOLEAN DEFAULT FALSE,
    is_default BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_phone_numbers table
CREATE TABLE IF NOT EXISTS user_phone_numbers (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phone_type VARCHAR(20) NOT NULL DEFAULT 'mobile' CHECK (phone_type IN ('mobile', 'home', 'work', 'fax', 'other')),
    phone_number VARCHAR(20) NOT NULL,
    country_code VARCHAR(5) NOT NULL DEFAULT '+1',
    is_verified BOOLEAN DEFAULT FALSE,
    is_primary BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_addresses_user_id ON user_addresses(user_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_addresses_address_type ON user_addresses(address_type) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_addresses_is_default ON user_addresses(is_default) WHERE is_default = TRUE AND is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_addresses_is_verified ON user_addresses(is_verified) WHERE is_verified = TRUE AND is_deleted = FALSE;

CREATE INDEX IF NOT EXISTS idx_user_phone_numbers_user_id ON user_phone_numbers(user_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_phone_numbers_phone_type ON user_phone_numbers(phone_type) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_phone_numbers_is_primary ON user_phone_numbers(is_primary) WHERE is_primary = TRUE AND is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_user_phone_numbers_is_verified ON user_phone_numbers(is_verified) WHERE is_verified = TRUE AND is_deleted = FALSE;

-- Create unique constraint to ensure only one default address per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_addresses_user_default 
    ON user_addresses(user_id) 
    WHERE is_default = TRUE AND is_deleted = FALSE;

-- Create unique constraint to ensure only one primary phone per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_phone_numbers_user_primary 
    ON user_phone_numbers(user_id) 
    WHERE is_primary = TRUE AND is_deleted = FALSE;

-- Create triggers for updated_at
CREATE TRIGGER update_user_addresses_updated_at 
    BEFORE UPDATE ON user_addresses 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_phone_numbers_updated_at 
    BEFORE UPDATE ON user_phone_numbers 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE user_addresses IS 'User saved addresses for orders and profile management';
COMMENT ON TABLE user_phone_numbers IS 'User phone numbers for contact and verification';

COMMENT ON COLUMN user_addresses.address_type IS 'Type of address: home, work, billing, shipping, other';
COMMENT ON COLUMN user_addresses.is_verified IS 'Whether the address has been verified';
COMMENT ON COLUMN user_addresses.is_default IS 'Whether this is the default address for the user';

COMMENT ON COLUMN user_phone_numbers.phone_type IS 'Type of phone: mobile, home, work, fax, other';
COMMENT ON COLUMN user_phone_numbers.is_verified IS 'Whether the phone number has been verified';
COMMENT ON COLUMN user_phone_numbers.is_primary IS 'Whether this is the primary phone for the user';
