-- Create coupons table
CREATE TABLE coupons (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('percentage', 'fixed_amount', 'free_shipping')),
    value DECIMAL(10,2) NOT NULL,
    minimum_amount DECIMAL(10,2) DEFAULT 0,
    maximum_discount DECIMAL(10,2),
    usage_limit INTEGER,
    used_count INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    starts_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for coupons
CREATE INDEX idx_coupons_code ON coupons(code);
CREATE INDEX idx_coupons_type ON coupons(type);
CREATE INDEX idx_coupons_is_active ON coupons(is_active);
CREATE INDEX idx_coupons_starts_at ON coupons(starts_at);
CREATE INDEX idx_coupons_expires_at ON coupons(expires_at);

-- Create coupon_usage table for tracking usage
CREATE TABLE coupon_usage (
    id BIGSERIAL PRIMARY KEY,
    coupon_id BIGINT NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    order_id BIGINT,
    user_id BIGINT,
    cart_id BIGINT,
    discount_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for coupon_usage
CREATE INDEX idx_coupon_usage_coupon_id ON coupon_usage(coupon_id);
CREATE INDEX idx_coupon_usage_order_id ON coupon_usage(order_id);
CREATE INDEX idx_coupon_usage_user_id ON coupon_usage(user_id);
CREATE INDEX idx_coupon_usage_cart_id ON coupon_usage(cart_id);

-- Create trigger to update updated_at timestamp for coupons
CREATE TRIGGER update_coupons_updated_at 
    BEFORE UPDATE ON coupons 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments
COMMENT ON TABLE coupons IS 'Discount coupons and promotional codes';
COMMENT ON TABLE coupon_usage IS 'Tracks coupon usage for analytics and limits';

COMMENT ON COLUMN coupons.type IS 'Coupon type: percentage, fixed_amount, or free_shipping';
COMMENT ON COLUMN coupons.value IS 'Discount value (percentage or fixed amount)';
COMMENT ON COLUMN coupons.minimum_amount IS 'Minimum cart amount required to use coupon';
COMMENT ON COLUMN coupons.maximum_discount IS 'Maximum discount amount for percentage coupons';
COMMENT ON COLUMN coupons.usage_limit IS 'Maximum number of times coupon can be used (NULL = unlimited)';
COMMENT ON COLUMN coupons.used_count IS 'Number of times coupon has been used';
COMMENT ON COLUMN coupons.starts_at IS 'When coupon becomes active';
COMMENT ON COLUMN coupons.expires_at IS 'When coupon expires (NULL = never expires)';
