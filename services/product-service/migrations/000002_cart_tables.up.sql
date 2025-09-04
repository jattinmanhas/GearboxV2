-- Create carts table
CREATE TABLE carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT, -- null for guest carts
    session_id VARCHAR(255) NOT NULL, -- for guest carts
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for carts
CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_session_id ON carts(session_id);
CREATE INDEX idx_carts_expires_at ON carts(expires_at);
CREATE INDEX idx_carts_created_at ON carts(created_at);

-- Create cart_items table
CREATE TABLE cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    product_variant_id BIGINT REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(cart_id, product_id, product_variant_id)
);

-- Create indexes for cart_items
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);
CREATE INDEX idx_cart_items_product_variant_id ON cart_items(product_variant_id);
CREATE INDEX idx_cart_items_created_at ON cart_items(created_at);

-- Create cart_coupons table
CREATE TABLE cart_coupons (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    coupon_code VARCHAR(50) NOT NULL,
    discount_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(cart_id, coupon_code)
);

-- Create indexes for cart_coupons
CREATE INDEX idx_cart_coupons_cart_id ON cart_coupons(cart_id);
CREATE INDEX idx_cart_coupons_coupon_code ON cart_coupons(coupon_code);

-- Create cart_shipping table
CREATE TABLE cart_shipping (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    shipping_method_id BIGINT NOT NULL,
    shipping_method VARCHAR(100) NOT NULL,
    shipping_amount DECIMAL(10,2) NOT NULL,
    estimated_days INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(cart_id)
);

-- Create indexes for cart_shipping
CREATE INDEX idx_cart_shipping_cart_id ON cart_shipping(cart_id);
CREATE INDEX idx_cart_shipping_shipping_method_id ON cart_shipping(shipping_method_id);

-- Create wishlists table
CREATE TABLE wishlists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, -- references users table (from auth service)
    name VARCHAR(100) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for wishlists
CREATE INDEX idx_wishlists_user_id ON wishlists(user_id);
CREATE INDEX idx_wishlists_is_public ON wishlists(is_public);
CREATE INDEX idx_wishlists_created_at ON wishlists(created_at);

-- Create wishlist_items table
CREATE TABLE wishlist_items (
    id BIGSERIAL PRIMARY KEY,
    wishlist_id BIGINT NOT NULL REFERENCES wishlists(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    product_variant_id BIGINT REFERENCES product_variants(id) ON DELETE CASCADE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(wishlist_id, product_id, product_variant_id)
);

-- Create indexes for wishlist_items
CREATE INDEX idx_wishlist_items_wishlist_id ON wishlist_items(wishlist_id);
CREATE INDEX idx_wishlist_items_product_id ON wishlist_items(product_id);
CREATE INDEX idx_wishlist_items_product_variant_id ON wishlist_items(product_variant_id);
CREATE INDEX idx_wishlist_items_created_at ON wishlist_items(created_at);

-- Create triggers to automatically update updated_at for cart tables
CREATE TRIGGER update_carts_updated_at BEFORE UPDATE ON carts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_cart_items_updated_at BEFORE UPDATE ON cart_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_wishlists_updated_at BEFORE UPDATE ON wishlists FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add constraint to ensure either user_id or session_id is provided
ALTER TABLE carts ADD CONSTRAINT check_cart_user_or_session 
    CHECK (user_id IS NOT NULL OR session_id IS NOT NULL);

-- Add constraint to ensure positive quantities and prices
ALTER TABLE cart_items ADD CONSTRAINT check_cart_items_positive_quantity 
    CHECK (quantity > 0);
ALTER TABLE cart_items ADD CONSTRAINT check_cart_items_positive_unit_price 
    CHECK (unit_price >= 0);
ALTER TABLE cart_items ADD CONSTRAINT check_cart_items_positive_total_price 
    CHECK (total_price >= 0);

-- Add constraint to ensure positive discount amount
ALTER TABLE cart_coupons ADD CONSTRAINT check_cart_coupons_positive_discount 
    CHECK (discount_amount >= 0);

-- Add constraint to ensure positive shipping amount and days
ALTER TABLE cart_shipping ADD CONSTRAINT check_cart_shipping_positive_amount 
    CHECK (shipping_amount >= 0);
ALTER TABLE cart_shipping ADD CONSTRAINT check_cart_shipping_positive_days 
    CHECK (estimated_days > 0);
