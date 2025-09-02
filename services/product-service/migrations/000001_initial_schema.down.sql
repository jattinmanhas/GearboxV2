-- Drop triggers first
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP TRIGGER IF EXISTS update_product_variants_updated_at ON product_variants;
DROP TRIGGER IF EXISTS update_product_attributes_updated_at ON product_attributes;
DROP TRIGGER IF EXISTS update_product_attribute_values_updated_at ON product_attribute_values;
DROP TRIGGER IF EXISTS update_product_images_updated_at ON product_images;
DROP TRIGGER IF EXISTS update_inventory_updated_at ON inventory;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (due to foreign key constraints)
DROP TABLE IF EXISTS inventory_alerts;
DROP TABLE IF EXISTS stock_reservations;
DROP TABLE IF EXISTS inventory_movements;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS product_attribute_values;
DROP TABLE IF EXISTS product_attributes;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
