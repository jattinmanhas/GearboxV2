-- Add phone_number column to users table
ALTER TABLE users ADD COLUMN phone_number VARCHAR(20);

-- Add comment for the new column
COMMENT ON COLUMN users.phone_number IS 'User phone number for contact purposes';
