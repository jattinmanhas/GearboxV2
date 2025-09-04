-- Drop cart-related tables in reverse order (due to foreign key constraints)

-- Drop triggers first
DROP TRIGGER IF EXISTS update_wishlists_updated_at ON wishlists;
DROP TRIGGER IF EXISTS update_cart_items_updated_at ON cart_items;
DROP TRIGGER IF EXISTS update_carts_updated_at ON carts;

-- Drop tables
DROP TABLE IF EXISTS wishlist_items;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS cart_shipping;
DROP TABLE IF EXISTS cart_coupons;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS carts;
