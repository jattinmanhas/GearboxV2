BEGIN;

TRUNCATE TABLE
    order_refunds,
    order_fulfillment,
    order_status_history,
    order_addresses,
    order_items,
    orders,
    coupon_usage,
    coupons,
    inventory_alerts,
    stock_reservations,
    inventory_movements,
    inventory,
    wishlist_items,
    wishlists,
    cart_shipping,
    cart_coupons,
    cart_items,
    carts,
    product_categories,
    product_images,
    product_attribute_values,
    product_attributes,
    product_variants,
    products,
    categories
RESTART IDENTITY CASCADE;

INSERT INTO categories (
    id, name, description, slug, parent_id, is_active, sort_order, image_url, image_public_id,
    meta_title, meta_description, created_at, updated_at
) VALUES
    (1, 'Apparel', 'Everyday apparel and performance basics.', 'apparel', NULL, TRUE, 1, 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=800&h=500&fit=crop', 'demo/categories/apparel', 'Gearbox Apparel', 'Performance-first clothing and outerwear.', '2026-01-02T09:00:00Z', '2026-03-20T09:00:00Z'),
    (2, 'Footwear', 'Shoes for training, commuting, and travel.', 'footwear', NULL, TRUE, 2, 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800&h=500&fit=crop', 'demo/categories/footwear', 'Gearbox Footwear', 'Sneakers and trail-ready styles.', '2026-01-02T09:10:00Z', '2026-03-20T09:10:00Z'),
    (3, 'Accessories', 'Useful add-ons for work and travel.', 'accessories', NULL, TRUE, 3, 'https://images.unsplash.com/photo-1523381210434-271e8be1f52b?w=800&h=500&fit=crop', 'demo/categories/accessories', 'Gearbox Accessories', 'Bags, headphones, and utility gear.', '2026-01-02T09:20:00Z', '2026-03-20T09:20:00Z'),
    (4, 'Home Office', 'Lighting and desk essentials.', 'home-office', NULL, TRUE, 4, 'https://images.unsplash.com/photo-1497366754035-f200968a6e72?w=800&h=500&fit=crop', 'demo/categories/home-office', 'Gearbox Home Office', 'Tools for focused workspaces.', '2026-01-02T09:30:00Z', '2026-03-20T09:30:00Z'),
    (5, 'Digital Goods', 'Downloadable guides and digital kits.', 'digital-goods', NULL, TRUE, 5, 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=800&h=500&fit=crop', 'demo/categories/digital-goods', 'Gearbox Digital Goods', 'Templates, guides, and quick-start kits.', '2026-01-02T09:40:00Z', '2026-03-20T09:40:00Z'),
    (6, 'Jackets', 'Layering and outerwear.', 'jackets', 1, TRUE, 1, 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=800&h=500&fit=crop', 'demo/categories/jackets', 'Gearbox Jackets', 'Mid-layers and overshirts.', '2026-01-02T09:50:00Z', '2026-03-20T09:50:00Z'),
    (7, 'Bags', 'Backpacks and travel organization.', 'bags', 3, TRUE, 1, 'https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=800&h=500&fit=crop', 'demo/categories/bags', 'Gearbox Bags', 'Carry systems for daily use.', '2026-01-02T10:00:00Z', '2026-03-20T10:00:00Z'),
    (8, 'Audio', 'Portable listening gear.', 'audio', 3, TRUE, 2, 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=800&h=500&fit=crop', 'demo/categories/audio', 'Gearbox Audio', 'Wireless audio and commuting essentials.', '2026-01-02T10:10:00Z', '2026-03-20T10:10:00Z');

INSERT INTO products (
    id, name, description, short_description, sku, price, compare_price, cost_price, weight, dimensions,
    is_active, is_digital, requires_shipping, meta_title, meta_description, tags, created_at, updated_at
) VALUES
    (1, 'Apex Runner Sneakers', 'Lightweight everyday sneakers designed for commute-to-gym days.', 'Responsive knit sneaker with all-day comfort.', 'GB-APEX-RUNNER', 89.00, 109.00, 41.00, 0.850, '32x20x12 cm', TRUE, FALSE, TRUE, 'Apex Runner Sneakers', 'Lightweight Gearbox sneakers for daily wear.', 'sneakers,commute,fitness', '2026-01-05T09:00:00Z', '2026-03-28T09:00:00Z'),
    (2, 'Terra Trek Backpack', 'A weather-resistant backpack with modular pockets and laptop sleeve.', 'Carry-on friendly backpack built for daily travel.', 'GB-TERRA-TREK', 79.99, 95.00, 36.00, 1.200, '48x31x16 cm', TRUE, FALSE, TRUE, 'Terra Trek Backpack', 'A versatile backpack for office and travel.', 'backpack,travel,office', '2026-01-07T10:00:00Z', '2026-03-28T10:00:00Z'),
    (3, 'Nimbus Overshirt', 'A midweight overshirt that layers well in cool offices and mild evenings.', 'Structured overshirt with soft brushed finish.', 'GB-NIMBUS-OVERSHIRT', 59.00, 74.00, 24.00, 0.600, '38x28x3 cm', TRUE, FALSE, TRUE, 'Nimbus Overshirt', 'Clean-lined overshirt for everyday layering.', 'overshirt,layering,apparel', '2026-01-10T11:00:00Z', '2026-03-28T11:00:00Z'),
    (4, 'Aurora Desk Lamp', 'A dimmable desk lamp with warm and cool presets for focused work.', 'Compact aluminum task lamp with USB-C power.', 'GB-AURORA-LAMP', 64.00, 79.00, 28.00, 1.700, '41x17x17 cm', TRUE, FALSE, TRUE, 'Aurora Desk Lamp', 'Modern desk lamp for home office setups.', 'desk,lighting,workspace', '2026-01-12T12:00:00Z', '2026-03-28T12:00:00Z'),
    (5, 'Pulse Wireless Headphones', 'Noise-reducing wireless headphones with 30-hour battery life.', 'Portable headphones tuned for commute and focus.', 'GB-PULSE-WIRELESS', 129.00, 159.00, 58.00, 0.420, '22x19x9 cm', TRUE, FALSE, TRUE, 'Pulse Wireless Headphones', 'Wireless over-ear headphones with long battery life.', 'audio,headphones,wireless', '2026-01-15T13:00:00Z', '2026-03-28T13:00:00Z'),
    (6, 'Starter Strength Training Plan', 'A downloadable 6-week training plan with mobility and recovery notes.', 'Digital training plan with printable weekly progress sheets.', 'GB-STARTER-PLAN', 29.00, 39.00, 3.00, 0.000, 'digital', TRUE, TRUE, FALSE, 'Starter Strength Training Plan', 'A digital Gearbox training plan for beginners.', 'digital,training,fitness', '2026-01-18T14:00:00Z', '2026-03-28T14:00:00Z');

INSERT INTO product_variants (
    id, product_id, name, sku, price, compare_price, cost_price, weight, is_active, position, created_at, updated_at
) VALUES
    (1, 1, 'Sand / US 8', 'GB-APEX-RUNNER-8', 89.00, 109.00, 41.00, 0.820, TRUE, 1, '2026-01-05T09:05:00Z', '2026-03-28T09:05:00Z'),
    (2, 1, 'Sand / US 9', 'GB-APEX-RUNNER-9', 89.00, 109.00, 41.00, 0.840, TRUE, 2, '2026-01-05T09:06:00Z', '2026-03-28T09:06:00Z'),
    (3, 1, 'Sand / US 10', 'GB-APEX-RUNNER-10', 89.00, 109.00, 41.00, 0.860, TRUE, 3, '2026-01-05T09:07:00Z', '2026-03-28T09:07:00Z'),
    (4, 3, 'Olive / Medium', 'GB-NIMBUS-OVERSHIRT-M', 59.00, 74.00, 24.00, 0.580, TRUE, 1, '2026-01-10T11:05:00Z', '2026-03-28T11:05:00Z'),
    (5, 3, 'Olive / Large', 'GB-NIMBUS-OVERSHIRT-L', 59.00, 74.00, 24.00, 0.610, TRUE, 2, '2026-01-10T11:06:00Z', '2026-03-28T11:06:00Z'),
    (6, 5, 'Midnight Black', 'GB-PULSE-WIRELESS-BLK', 129.00, 159.00, 58.00, 0.420, TRUE, 1, '2026-01-15T13:05:00Z', '2026-03-28T13:05:00Z'),
    (7, 5, 'Stone White', 'GB-PULSE-WIRELESS-WHT', 129.00, 159.00, 58.00, 0.420, TRUE, 2, '2026-01-15T13:06:00Z', '2026-03-28T13:06:00Z');

INSERT INTO product_attributes (id, name, value, type, created_at, updated_at) VALUES
    (1, 'Color', 'sand', 'select', '2026-01-04T08:00:00Z', '2026-03-01T08:00:00Z'),
    (2, 'Size', 'us', 'select', '2026-01-04T08:05:00Z', '2026-03-01T08:05:00Z'),
    (3, 'Material', 'fabric', 'text', '2026-01-04T08:10:00Z', '2026-03-01T08:10:00Z'),
    (4, 'Capacity', 'liters', 'number', '2026-01-04T08:15:00Z', '2026-03-01T08:15:00Z'),
    (5, 'Connectivity', 'bluetooth', 'text', '2026-01-04T08:20:00Z', '2026-03-01T08:20:00Z');

INSERT INTO product_attribute_values (id, product_id, attribute_id, value, created_at, updated_at) VALUES
    (1, 1, 1, 'Sand', '2026-01-05T09:10:00Z', '2026-03-28T09:10:00Z'),
    (2, 1, 2, '8-10', '2026-01-05T09:11:00Z', '2026-03-28T09:11:00Z'),
    (3, 2, 4, '24', '2026-01-07T10:10:00Z', '2026-03-28T10:10:00Z'),
    (4, 3, 3, 'Brushed cotton blend', '2026-01-10T11:10:00Z', '2026-03-28T11:10:00Z'),
    (5, 5, 5, 'Bluetooth 5.3', '2026-01-15T13:10:00Z', '2026-03-28T13:10:00Z');

INSERT INTO product_images (
    id, product_id, url, alt, public_id, position, is_primary, created_at, updated_at
) VALUES
    (1, 1, 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=900&h=900&fit=crop', 'Apex Runner Sneakers', 'demo/products/apex-runner-1', 1, TRUE, '2026-01-05T09:20:00Z', '2026-03-28T09:20:00Z'),
    (2, 1, 'https://images.unsplash.com/photo-1549298916-b41d501d3772?w=900&h=900&fit=crop', 'Apex Runner side view', 'demo/products/apex-runner-2', 2, FALSE, '2026-01-05T09:21:00Z', '2026-03-28T09:21:00Z'),
    (3, 2, 'https://images.unsplash.com/photo-1548036328-c9fa89d128fa?w=900&h=900&fit=crop', 'Terra Trek Backpack', 'demo/products/terra-trek-1', 1, TRUE, '2026-01-07T10:20:00Z', '2026-03-28T10:20:00Z'),
    (4, 3, 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=900&h=900&fit=crop', 'Nimbus Overshirt', 'demo/products/nimbus-1', 1, TRUE, '2026-01-10T11:20:00Z', '2026-03-28T11:20:00Z'),
    (5, 4, 'https://images.unsplash.com/photo-1497366754035-f200968a6e72?w=900&h=900&fit=crop', 'Aurora Desk Lamp', 'demo/products/aurora-lamp-1', 1, TRUE, '2026-01-12T12:20:00Z', '2026-03-28T12:20:00Z'),
    (6, 5, 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=900&h=900&fit=crop', 'Pulse Wireless Headphones', 'demo/products/pulse-wireless-1', 1, TRUE, '2026-01-15T13:20:00Z', '2026-03-28T13:20:00Z'),
    (7, 6, 'https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=900&h=900&fit=crop', 'Starter Strength Training Plan', 'demo/products/starter-plan-1', 1, TRUE, '2026-01-18T14:20:00Z', '2026-03-28T14:20:00Z');

INSERT INTO product_categories (id, product_id, category_id, is_primary, created_at) VALUES
    (1, 1, 2, TRUE, '2026-01-05T09:30:00Z'),
    (2, 1, 1, FALSE, '2026-01-05T09:31:00Z'),
    (3, 2, 7, TRUE, '2026-01-07T10:30:00Z'),
    (4, 2, 3, FALSE, '2026-01-07T10:31:00Z'),
    (5, 3, 6, TRUE, '2026-01-10T11:30:00Z'),
    (6, 3, 1, FALSE, '2026-01-10T11:31:00Z'),
    (7, 4, 4, TRUE, '2026-01-12T12:30:00Z'),
    (8, 5, 8, TRUE, '2026-01-15T13:30:00Z'),
    (9, 5, 3, FALSE, '2026-01-15T13:31:00Z'),
    (10, 6, 5, TRUE, '2026-01-18T14:30:00Z');

INSERT INTO carts (id, user_id, session_id, currency, created_at, updated_at, expires_at) VALUES
    (1, 4, 'sess_mason_001', 'USD', '2026-03-30T10:00:00Z', '2026-03-31T07:30:00Z', '2026-04-07T10:00:00Z'),
    (2, NULL, 'guest_demo_cart_001', 'USD', '2026-03-31T06:45:00Z', '2026-03-31T07:50:00Z', '2026-04-02T06:45:00Z');

INSERT INTO cart_items (id, cart_id, product_id, product_variant_id, quantity, unit_price, total_price, created_at, updated_at) VALUES
    (1, 1, 1, 2, 1, 89.00, 89.00, '2026-03-30T10:05:00Z', '2026-03-31T07:30:00Z'),
    (2, 1, 4, NULL, 1, 64.00, 64.00, '2026-03-30T10:10:00Z', '2026-03-31T07:30:00Z'),
    (3, 2, 2, NULL, 1, 79.99, 79.99, '2026-03-31T06:50:00Z', '2026-03-31T07:50:00Z'),
    (4, 2, 6, NULL, 1, 29.00, 29.00, '2026-03-31T06:52:00Z', '2026-03-31T07:50:00Z');

INSERT INTO cart_coupons (id, cart_id, coupon_code, discount_amount, created_at) VALUES
    (1, 1, 'SAVE10', 10.00, '2026-03-31T07:00:00Z');

INSERT INTO cart_shipping (id, cart_id, shipping_method_id, shipping_method, estimated_days, created_at) VALUES
    (1, 1, 101, 'Standard Shipping', 4, '2026-03-31T07:05:00Z'),
    (2, 2, 201, 'Express Shipping', 2, '2026-03-31T07:10:00Z');

INSERT INTO wishlists (id, user_id, name, is_public, created_at, updated_at) VALUES
    (1, 4, 'Spring Upgrade Picks', TRUE, '2026-03-05T09:00:00Z', '2026-03-29T20:00:00Z'),
    (2, 5, 'Travel Setup', FALSE, '2026-03-09T12:00:00Z', '2026-03-30T18:00:00Z');

INSERT INTO wishlist_items (id, wishlist_id, product_id, product_variant_id, notes, created_at) VALUES
    (1, 1, 5, 6, 'Want these for flight days.', '2026-03-05T09:10:00Z'),
    (2, 1, 2, NULL, 'Need a bigger backpack soon.', '2026-03-05T09:12:00Z'),
    (3, 2, 1, 3, 'Preferred size when it is back in stock.', '2026-03-09T12:10:00Z');

INSERT INTO inventory (
    id, product_id, product_variant_id, quantity, reserved_quantity, available_quantity,
    min_stock_level, max_stock_level, reorder_point, last_restocked, created_at, updated_at
) VALUES
    (1, 1, 1, 14, 1, 13, 4, 25, 6, '2026-03-20T08:00:00Z', '2026-01-05T09:40:00Z', '2026-03-31T07:30:00Z'),
    (2, 1, 2, 9, 2, 7, 4, 25, 6, '2026-03-20T08:00:00Z', '2026-01-05T09:41:00Z', '2026-03-31T07:30:00Z'),
    (3, 1, 3, 3, 0, 3, 4, 25, 6, '2026-03-20T08:00:00Z', '2026-01-05T09:42:00Z', '2026-03-31T07:30:00Z'),
    (4, 2, NULL, 18, 1, 17, 5, 40, 8, '2026-03-18T08:00:00Z', '2026-01-07T10:40:00Z', '2026-03-31T07:45:00Z'),
    (5, 3, 4, 11, 1, 10, 3, 20, 5, '2026-03-19T08:00:00Z', '2026-01-10T11:40:00Z', '2026-03-31T07:10:00Z'),
    (6, 3, 5, 6, 0, 6, 3, 20, 5, '2026-03-19T08:00:00Z', '2026-01-10T11:41:00Z', '2026-03-31T07:10:00Z'),
    (7, 4, NULL, 4, 0, 4, 3, 12, 4, '2026-03-15T08:00:00Z', '2026-01-12T12:40:00Z', '2026-03-31T06:20:00Z'),
    (8, 5, 6, 20, 0, 20, 5, 35, 8, '2026-03-17T08:00:00Z', '2026-01-15T13:40:00Z', '2026-03-31T06:30:00Z'),
    (9, 5, 7, 8, 1, 7, 5, 35, 8, '2026-03-17T08:00:00Z', '2026-01-15T13:41:00Z', '2026-03-31T06:30:00Z'),
    (10, 6, NULL, 999, 0, 999, 0, 9999, 0, '2026-03-01T08:00:00Z', '2026-01-18T14:40:00Z', '2026-03-31T06:00:00Z');

INSERT INTO inventory_movements (
    id, product_id, product_variant_id, movement_type, quantity, previous_quantity, new_quantity,
    reference, reference_type, reason, notes, created_by, created_at
) VALUES
    (1, 1, 2, 'out', 1, 10, 9, '202603280001', 'order', 'Customer order placed', 'Order shipped successfully.', 1, '2026-03-28T10:00:00Z'),
    (2, 2, NULL, 'in', 6, 12, 18, 'PO-2403-18', 'purchase', 'Warehouse restock', 'Backpacks restocked after supplier delivery.', 1, '2026-03-18T09:00:00Z'),
    (3, 4, NULL, 'adjustment', 2, 6, 4, 'AUDIT-2026-03', 'adjustment', 'Damaged unit write-off', 'Two lamps marked unsellable after inspection.', 1, '2026-03-24T15:30:00Z'),
    (4, 5, 7, 'out', 1, 9, 8, '202603290004', 'order', 'Reserved for recent order', 'Stone White unit reserved for refunded order.', 1, '2026-03-29T11:30:00Z');

INSERT INTO inventory_alerts (
    id, product_id, product_variant_id, alert_type, current_quantity, threshold_quantity,
    is_resolved, resolved_at, created_at
) VALUES
    (1, 1, 3, 'low_stock', 3, 4, FALSE, NULL, '2026-03-31T07:35:00Z'),
    (2, 4, NULL, 'reorder_point', 4, 4, FALSE, NULL, '2026-03-31T06:25:00Z');

INSERT INTO coupons (
    id, code, name, description, type, value, minimum_amount, maximum_discount, usage_limit,
    used_count, is_active, starts_at, expires_at, created_at, updated_at
) VALUES
    (1, 'SAVE10', 'Save 10%', 'Ten percent off orders above $75.', 'percentage', 10.00, 75.00, 20.00, 500, 1, TRUE, '2026-03-01T00:00:00Z', '2026-05-01T00:00:00Z', '2026-03-01T00:00:00Z', '2026-03-31T07:00:00Z'),
    (2, 'WELCOME25', 'Welcome $25', 'Fixed amount off for first-time buyers over $100.', 'fixed_amount', 25.00, 100.00, NULL, 250, 1, TRUE, '2026-03-10T00:00:00Z', '2026-04-30T00:00:00Z', '2026-03-10T00:00:00Z', '2026-03-30T12:00:00Z'),
    (3, 'FREESHIP', 'Free Shipping', 'Removes shipping cost for standard delivery.', 'free_shipping', 0.00, 50.00, NULL, NULL, 0, TRUE, '2026-03-01T00:00:00Z', NULL, '2026-03-01T00:00:00Z', '2026-03-31T09:00:00Z'),
    (4, 'SPRING15', 'Spring Promo', 'Seasonal promotion that is currently paused.', 'percentage', 15.00, 120.00, 40.00, 100, 0, FALSE, '2026-02-15T00:00:00Z', '2026-03-15T00:00:00Z', '2026-02-15T00:00:00Z', '2026-03-16T00:00:00Z');

INSERT INTO orders (
    id, order_number, user_id, status, payment_status, fulfillment_status, subtotal,
    discount_amount, total_amount, currency, notes, internal_notes, created_at, updated_at,
    confirmed_at, shipped_at, delivered_at, cancelled_at
) VALUES
    (1, 'GBX-20260328-0001', 4, 'delivered', 'paid', 'fulfilled', 153.00, 10.00, 143.00, 'USD', 'Leave package at front desk.', 'Priority customer order.', '2026-03-28T10:00:00Z', '2026-03-30T18:00:00Z', '2026-03-28T10:05:00Z', '2026-03-28T16:00:00Z', '2026-03-30T15:30:00Z', NULL),
    (2, 'GBX-20260329-0002', 3, 'processing', 'paid', 'partial', 138.99, 0.00, 138.99, 'USD', 'Please notify before dispatch.', 'Manual packing check requested.', '2026-03-29T11:15:00Z', '2026-03-31T07:20:00Z', '2026-03-29T11:20:00Z', NULL, NULL, NULL),
    (3, 'GBX-20260330-0003', 2, 'cancelled', 'failed', 'unfulfilled', 129.00, 0.00, 129.00, 'USD', 'Customer payment attempt failed twice.', 'Cancelled after retry window closed.', '2026-03-30T13:40:00Z', '2026-03-30T15:10:00Z', NULL, NULL, NULL, '2026-03-30T15:10:00Z'),
    (4, 'GBX-20260331-0004', 5, 'refunded', 'refunded', 'fulfilled', 108.99, 25.00, 83.99, 'USD', 'Digital plan delivered immediately.', 'Refunded after duplicate order report.', '2026-03-31T08:10:00Z', '2026-03-31T11:00:00Z', '2026-03-31T08:15:00Z', '2026-03-31T09:10:00Z', '2026-03-31T09:10:00Z', NULL);

INSERT INTO order_items (
    id, order_id, product_id, product_variant_id, product_name, product_sku, quantity,
    unit_price, total_price, discount_amount, is_digital, requires_shipping
) VALUES
    (1, 1, 1, 2, 'Apex Runner Sneakers - Sand / US 9', 'GB-APEX-RUNNER-9', 1, 89.00, 89.00, 6.00, FALSE, TRUE),
    (2, 1, 4, NULL, 'Aurora Desk Lamp', 'GB-AURORA-LAMP', 1, 64.00, 64.00, 4.00, FALSE, TRUE),
    (3, 2, 2, NULL, 'Terra Trek Backpack', 'GB-TERRA-TREK', 1, 79.99, 79.99, 0.00, FALSE, TRUE),
    (4, 2, 3, 4, 'Nimbus Overshirt - Olive / Medium', 'GB-NIMBUS-OVERSHIRT-M', 1, 59.00, 59.00, 0.00, FALSE, TRUE),
    (5, 3, 5, 6, 'Pulse Wireless Headphones - Midnight Black', 'GB-PULSE-WIRELESS-BLK', 1, 129.00, 129.00, 0.00, FALSE, TRUE),
    (6, 4, 2, NULL, 'Terra Trek Backpack', 'GB-TERRA-TREK', 1, 79.99, 79.99, 18.37, FALSE, TRUE),
    (7, 4, 6, NULL, 'Starter Strength Training Plan', 'GB-STARTER-PLAN', 1, 29.00, 29.00, 6.63, TRUE, FALSE);

INSERT INTO order_addresses (
    id, order_id, type, first_name, last_name, company, address1, address2, city, state, country, postal_code, phone, email
) VALUES
    (1, 1, 'shipping', 'Mason', 'Lee', NULL, '88 Cedar Avenue', NULL, 'Seattle', 'Washington', 'United States', '98101', '+1-202-555-0104', 'mason.shop@gearbox.local'),
    (2, 1, 'billing', 'Mason', 'Lee', 'Northwind Labs', '200 Industrial Road', 'Building C', 'Seattle', 'Washington', 'United States', '98109', '+1-202-555-0174', 'billing@northwind.example'),
    (3, 2, 'shipping', 'Nina', 'Patel', NULL, '42 River Walk', 'Apt 9B', 'Austin', 'Texas', 'United States', '78701', '+1-202-555-0103', 'nina.ops@gearbox.local'),
    (4, 2, 'billing', 'Nina', 'Patel', NULL, '42 River Walk', 'Apt 9B', 'Austin', 'Texas', 'United States', '78701', '+1-202-555-0103', 'nina.ops@gearbox.local'),
    (5, 3, 'shipping', 'Eli', 'Editor', NULL, '55 Beacon Street', NULL, 'Boston', 'Massachusetts', 'United States', '02108', '+1-202-555-0102', 'eli.editor@gearbox.local'),
    (6, 3, 'billing', 'Eli', 'Editor', NULL, '55 Beacon Street', NULL, 'Boston', 'Massachusetts', 'United States', '02108', '+1-202-555-0102', 'eli.editor@gearbox.local'),
    (7, 4, 'shipping', 'Zoe', 'Kim', NULL, '17 Lake Shore Drive', NULL, 'Chicago', 'Illinois', 'United States', '60601', '+1-202-555-0105', 'zoe.oauth@gearbox.local'),
    (8, 4, 'billing', 'Zoe', 'Kim', NULL, '17 Lake Shore Drive', NULL, 'Chicago', 'Illinois', 'United States', '60601', '+1-202-555-0105', 'zoe.oauth@gearbox.local');

INSERT INTO order_status_history (
    id, order_id, status, previous_status, notes, created_by, created_at
) VALUES
    (1, 1, 'confirmed', 'pending', 'Payment captured successfully.', 1, '2026-03-28T10:05:00Z'),
    (2, 1, 'shipped', 'confirmed', 'Handed off to carrier.', 1, '2026-03-28T16:00:00Z'),
    (3, 1, 'delivered', 'shipped', 'Delivery confirmed by carrier scan.', 1, '2026-03-30T15:30:00Z'),
    (4, 2, 'confirmed', 'pending', 'Order queued for picking.', 1, '2026-03-29T11:20:00Z'),
    (5, 2, 'processing', 'confirmed', 'One item packed, waiting on second quality check.', 1, '2026-03-31T07:20:00Z'),
    (6, 3, 'cancelled', 'pending', 'Payment authorization failed.', 1, '2026-03-30T15:10:00Z'),
    (7, 4, 'confirmed', 'pending', 'Order auto-confirmed after checkout.', 1, '2026-03-31T08:15:00Z'),
    (8, 4, 'refunded', 'confirmed', 'Customer reported duplicate purchase.', 1, '2026-03-31T11:00:00Z');

INSERT INTO order_fulfillment (
    id, order_id, tracking_number, carrier, service, status, shipped_at, delivered_at,
    estimated_delivery, notes, created_at, updated_at
) VALUES
    (1, 1, '1Z999AA10123456784', 'UPS', 'Ground', 'delivered', '2026-03-28T16:00:00Z', '2026-03-30T15:30:00Z', '2026-03-31T18:00:00Z', 'Delivered to front desk.', '2026-03-28T16:00:00Z', '2026-03-30T15:30:00Z'),
    (2, 2, '9400111899223857123456', 'USPS', 'Priority Mail', 'pending', NULL, NULL, '2026-04-03T18:00:00Z', 'Label created; awaiting final handoff.', '2026-03-31T07:20:00Z', '2026-03-31T07:20:00Z'),
    (3, 4, 'DIGITAL-DELIVERY-0004', 'Digital', 'Instant Download', 'delivered', '2026-03-31T09:10:00Z', '2026-03-31T09:10:00Z', '2026-03-31T09:10:00Z', 'Digital asset link sent automatically.', '2026-03-31T09:10:00Z', '2026-03-31T09:10:00Z');

INSERT INTO order_refunds (
    id, order_id, amount, reason, status, processed_at, created_by, created_at
) VALUES
    (1, 4, 83.99, 'Duplicate order reported by customer.', 'processed', '2026-03-31T11:00:00Z', 1, '2026-03-31T10:45:00Z');

INSERT INTO stock_reservations (
    id, product_id, product_variant_id, order_id, quantity, expires_at, created_at
) VALUES
    (1, 1, 2, 1, 1, '2026-04-02T10:00:00Z', '2026-03-28T10:00:00Z'),
    (2, 2, NULL, 2, 1, '2026-04-03T11:15:00Z', '2026-03-29T11:15:00Z'),
    (3, 3, 4, 2, 1, '2026-04-03T11:15:00Z', '2026-03-29T11:15:00Z'),
    (4, 5, 7, 4, 1, '2026-04-01T08:10:00Z', '2026-03-31T08:10:00Z');

INSERT INTO coupon_usage (
    id, coupon_id, order_id, user_id, cart_id, discount_amount, created_at
) VALUES
    (1, 1, 1, 4, 1, 10.00, '2026-03-28T10:00:00Z'),
    (2, 2, 4, 5, NULL, 25.00, '2026-03-31T08:10:00Z');

SELECT setval('categories_id_seq', COALESCE((SELECT MAX(id) FROM categories), 1), TRUE);
SELECT setval('products_id_seq', COALESCE((SELECT MAX(id) FROM products), 1), TRUE);
SELECT setval('product_variants_id_seq', COALESCE((SELECT MAX(id) FROM product_variants), 1), TRUE);
SELECT setval('product_attributes_id_seq', COALESCE((SELECT MAX(id) FROM product_attributes), 1), TRUE);
SELECT setval('product_attribute_values_id_seq', COALESCE((SELECT MAX(id) FROM product_attribute_values), 1), TRUE);
SELECT setval('product_images_id_seq', COALESCE((SELECT MAX(id) FROM product_images), 1), TRUE);
SELECT setval('product_categories_id_seq', COALESCE((SELECT MAX(id) FROM product_categories), 1), TRUE);
SELECT setval('carts_id_seq', COALESCE((SELECT MAX(id) FROM carts), 1), TRUE);
SELECT setval('cart_items_id_seq', COALESCE((SELECT MAX(id) FROM cart_items), 1), TRUE);
SELECT setval('cart_coupons_id_seq', COALESCE((SELECT MAX(id) FROM cart_coupons), 1), TRUE);
SELECT setval('cart_shipping_id_seq', COALESCE((SELECT MAX(id) FROM cart_shipping), 1), TRUE);
SELECT setval('wishlists_id_seq', COALESCE((SELECT MAX(id) FROM wishlists), 1), TRUE);
SELECT setval('wishlist_items_id_seq', COALESCE((SELECT MAX(id) FROM wishlist_items), 1), TRUE);
SELECT setval('inventory_id_seq', COALESCE((SELECT MAX(id) FROM inventory), 1), TRUE);
SELECT setval('inventory_movements_id_seq', COALESCE((SELECT MAX(id) FROM inventory_movements), 1), TRUE);
SELECT setval('inventory_alerts_id_seq', COALESCE((SELECT MAX(id) FROM inventory_alerts), 1), TRUE);
SELECT setval('stock_reservations_id_seq', COALESCE((SELECT MAX(id) FROM stock_reservations), 1), TRUE);
SELECT setval('coupons_id_seq', COALESCE((SELECT MAX(id) FROM coupons), 1), TRUE);
SELECT setval('coupon_usage_id_seq', COALESCE((SELECT MAX(id) FROM coupon_usage), 1), TRUE);
SELECT setval('orders_id_seq', COALESCE((SELECT MAX(id) FROM orders), 1), TRUE);
SELECT setval('order_items_id_seq', COALESCE((SELECT MAX(id) FROM order_items), 1), TRUE);
SELECT setval('order_addresses_id_seq', COALESCE((SELECT MAX(id) FROM order_addresses), 1), TRUE);
SELECT setval('order_status_history_id_seq', COALESCE((SELECT MAX(id) FROM order_status_history), 1), TRUE);
SELECT setval('order_fulfillment_id_seq', COALESCE((SELECT MAX(id) FROM order_fulfillment), 1), TRUE);
SELECT setval('order_refunds_id_seq', COALESCE((SELECT MAX(id) FROM order_refunds), 1), TRUE);

COMMIT;
