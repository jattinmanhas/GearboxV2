-- Remove tax_amount and shipping_amount from orders table
ALTER TABLE orders DROP COLUMN tax_amount;
ALTER TABLE orders DROP COLUMN shipping_amount;

-- Remove tax_amount from order_items table
ALTER TABLE order_items DROP COLUMN tax_amount;

-- Remove shipping_amount from cart_shipping table
ALTER TABLE cart_shipping DROP COLUMN shipping_amount;
