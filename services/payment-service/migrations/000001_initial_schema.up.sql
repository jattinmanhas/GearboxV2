-- Create payment_methods table
CREATE TABLE payment_methods (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('credit_card', 'debit_card', 'paypal', 'bank_transfer', 'digital_wallet')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_default BOOLEAN NOT NULL DEFAULT false,
    sort_order INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    icon VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for payment_methods
CREATE INDEX idx_payment_methods_code ON payment_methods(code);
CREATE INDEX idx_payment_methods_type ON payment_methods(type);
CREATE INDEX idx_payment_methods_is_active ON payment_methods(is_active);
CREATE INDEX idx_payment_methods_sort_order ON payment_methods(sort_order);

-- Create payment_gateways table
CREATE TABLE payment_gateways (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_test_mode BOOLEAN NOT NULL DEFAULT true,
    config TEXT NOT NULL, -- JSON configuration
    webhook_url VARCHAR(500),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for payment_gateways
CREATE INDEX idx_payment_gateways_code ON payment_gateways(code);
CREATE INDEX idx_payment_gateways_is_active ON payment_gateways(is_active);
CREATE INDEX idx_payment_gateways_sort_order ON payment_gateways(sort_order);

-- Create payments table
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL, -- References orders table (from product service)
    payment_method_id BIGINT NOT NULL REFERENCES payment_methods(id),
    transaction_id VARCHAR(255) NOT NULL UNIQUE,
    gateway_id VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled', 'refunded')),
    gateway_status VARCHAR(50),
    gateway_response TEXT,
    failure_reason TEXT,
    processed_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for payments
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_payment_method_id ON payments(payment_method_id);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_gateway_id ON payments(gateway_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);
CREATE INDEX idx_payments_processed_at ON payments(processed_at);

-- Create payment_refunds table
CREATE TABLE payment_refunds (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    refund_id VARCHAR(255) NOT NULL UNIQUE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processed', 'failed')),
    gateway_response TEXT,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_by BIGINT NOT NULL, -- References users table (from auth service)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for payment_refunds
CREATE INDEX idx_payment_refunds_payment_id ON payment_refunds(payment_id);
CREATE INDEX idx_payment_refunds_refund_id ON payment_refunds(refund_id);
CREATE INDEX idx_payment_refunds_status ON payment_refunds(status);
CREATE INDEX idx_payment_refunds_created_at ON payment_refunds(created_at);
CREATE INDEX idx_payment_refunds_created_by ON payment_refunds(created_by);

-- Create payment_webhooks table
CREATE TABLE payment_webhooks (
    id BIGSERIAL PRIMARY KEY,
    gateway_id VARCHAR(50) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_id VARCHAR(255) NOT NULL,
    payload TEXT NOT NULL, -- JSON payload
    signature VARCHAR(500) NOT NULL,
    is_processed BOOLEAN NOT NULL DEFAULT false,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for payment_webhooks
CREATE INDEX idx_payment_webhooks_gateway_id ON payment_webhooks(gateway_id);
CREATE INDEX idx_payment_webhooks_event_type ON payment_webhooks(event_type);
CREATE INDEX idx_payment_webhooks_event_id ON payment_webhooks(event_id);
CREATE INDEX idx_payment_webhooks_is_processed ON payment_webhooks(is_processed);
CREATE INDEX idx_payment_webhooks_created_at ON payment_webhooks(created_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at
CREATE TRIGGER update_payment_methods_updated_at 
    BEFORE UPDATE ON payment_methods 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payment_gateways_updated_at 
    BEFORE UPDATE ON payment_gateways 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at 
    BEFORE UPDATE ON payments 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default payment methods
INSERT INTO payment_methods (name, code, type, is_active, is_default, sort_order, description, icon) VALUES
('Credit Card', 'credit_card', 'credit_card', true, true, 1, 'Pay with your credit card', 'credit-card'),
('Debit Card', 'debit_card', 'debit_card', true, false, 2, 'Pay with your debit card', 'debit-card'),
('PayPal', 'paypal', 'paypal', true, false, 3, 'Pay with PayPal', 'paypal'),
('Bank Transfer', 'bank_transfer', 'bank_transfer', true, false, 4, 'Direct bank transfer', 'bank'),
('Digital Wallet', 'digital_wallet', 'digital_wallet', true, false, 5, 'Pay with digital wallet', 'wallet');

-- Insert default payment gateways
INSERT INTO payment_gateways (name, code, is_active, is_test_mode, config, sort_order) VALUES
('Stripe', 'stripe', true, true, '{"publishable_key": "pk_test_...", "secret_key": "sk_test_...", "webhook_secret": "whsec_..."}', 1),
('PayPal', 'paypal', true, true, '{"client_id": "...", "client_secret": "...", "mode": "sandbox"}', 2),
('Razorpay', 'razorpay', true, true, '{"key_id": "...", "key_secret": "...", "webhook_secret": "..."}', 3);
