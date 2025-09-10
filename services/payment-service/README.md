# Payment Service

A comprehensive payment processing service for the Gearbox e-commerce platform. This service handles payment methods, payment gateways, payment processing, refunds, and webhook management.

## Features

### Payment Methods
- **Credit Card**: Visa, Mastercard, American Express, etc.
- **Debit Card**: Direct debit payments
- **PayPal**: PayPal integration
- **Bank Transfer**: Direct bank transfers
- **Digital Wallet**: Apple Pay, Google Pay, etc.

### Payment Gateways
- **Stripe**: Full Stripe integration with webhooks
- **PayPal**: PayPal payment processing
- **Razorpay**: Razorpay integration for Indian markets

### Core Functionality
- Payment creation and processing
- Payment status management
- Refund processing
- Webhook handling
- Payment analytics and reporting
- Multi-currency support
- Transaction tracking

### Technical Features
- **SQLx Integration**: Enhanced database operations with better type safety
- **Named Queries**: Cleaner SQL with named parameters
- **Struct Scanning**: Automatic mapping of database rows to Go structs
- **Transaction Support**: Built-in transaction management

## API Endpoints

### Public Endpoints
- `GET /api/v1/public/payment-methods` - List available payment methods
- `GET /api/v1/public/payment-gateways` - List available payment gateways
- `POST /api/v1/public/webhooks/{gatewayId}` - Handle webhook events

### Protected Endpoints (Authentication Required)
- `POST /api/v1/protected/payments` - Create a new payment
- `POST /api/v1/protected/payments/process` - Process a payment
- `GET /api/v1/protected/payments` - List payments with filtering
- `GET /api/v1/protected/payments/{id}` - Get payment details
- `GET /api/v1/protected/payments/transaction/{transactionId}` - Get payment by transaction ID
- `PUT /api/v1/protected/payments/{id}/status` - Update payment status
- `POST /api/v1/protected/payments/refund` - Process a refund
- `GET /api/v1/protected/payments/summary` - Get payment summary statistics

### Admin Endpoints (Admin Role Required)
- `POST /api/v1/protected/payment-methods` - Create payment method
- `GET /api/v1/protected/payment-methods/{id}` - Get payment method
- `PUT /api/v1/protected/payment-methods/{id}` - Update payment method
- `DELETE /api/v1/protected/payment-methods/{id}` - Delete payment method
- `POST /api/v1/protected/payment-gateways` - Create payment gateway
- `GET /api/v1/protected/payment-gateways/{id}` - Get payment gateway
- `PUT /api/v1/protected/payment-gateways/{id}` - Update payment gateway
- `DELETE /api/v1/protected/payment-gateways/{id}` - Delete payment gateway

## Configuration

Copy `env.example` to `.env` and configure the following:

### Database
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gearbox_payment
DB_SSL_MODE=disable
```

### Server
```env
PORT=8083
HOST=localhost
```

### JWT
```env
JWT_SECRET=your-jwt-secret-key-here
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

### Payment Gateways

#### Stripe
```env
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here
```

#### PayPal
```env
PAYPAL_CLIENT_ID=your_paypal_client_id_here
PAYPAL_CLIENT_SECRET=your_paypal_client_secret_here
PAYPAL_MODE=sandbox
```

#### Razorpay
```env
RAZORPAY_KEY_ID=your_razorpay_key_id_here
RAZORPAY_KEY_SECRET=your_razorpay_key_secret_here
RAZORPAY_WEBHOOK_SECRET=your_razorpay_webhook_secret_here
```

### Service URLs
```env
AUTH_SERVICE_URL=http://localhost:8081
PRODUCT_SERVICE_URL=http://localhost:8082
```

## Database Setup

1. Create the database:
```sql
CREATE DATABASE gearbox_payment;
```

2. Run migrations:
```bash
cd migrations
chmod +x run-migrations.sh
./run-migrations.sh up
```

## Running the Service

1. Install dependencies:
```bash
go mod tidy
```

2. Set up environment variables:
```bash
cp env.example .env
# Edit .env with your configuration
```

3. Run the service:
```bash
go run cmd/api/main.go
```

The service will start on `http://localhost:8083` by default.

## API Usage Examples

### Create a Payment
```bash
curl -X POST http://localhost:8083/api/v1/protected/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "order_id": 123,
    "payment_method_id": 1,
    "amount": 99.99,
    "currency": "USD",
    "gateway_id": "stripe",
    "metadata": {
      "customer_id": "cust_123",
      "description": "Order #123"
    }
  }'
```

### Process a Payment
```bash
curl -X POST http://localhost:8083/api/v1/protected/payments/process \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "payment_id": 1,
    "payment_data": {
      "payment_method_id": "pm_1234567890",
      "customer_id": "cust_123"
    }
  }'
```

### List Payments
```bash
curl -X GET "http://localhost:8083/api/v1/protected/payments?status=completed&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Process a Refund
```bash
curl -X POST http://localhost:8083/api/v1/protected/payments/refund \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "payment_id": 1,
    "amount": 50.00,
    "reason": "Customer requested partial refund"
  }'
```

## Webhook Handling

The service supports webhooks from various payment gateways:

- **Stripe**: `POST /api/v1/public/webhooks/stripe`
- **PayPal**: `POST /api/v1/public/webhooks/paypal`
- **Razorpay**: `POST /api/v1/public/webhooks/razorpay`

Webhook events are automatically processed and payment statuses are updated accordingly.

## Payment Flow

1. **Create Payment**: Create a payment record with order details
2. **Process Payment**: Send payment to gateway for processing
3. **Webhook Processing**: Gateway sends webhook with payment status
4. **Status Update**: Payment status is updated based on webhook
5. **Refund Processing**: Handle refunds if needed

## Security

- JWT-based authentication
- Role-based access control (RBAC)
- Webhook signature validation
- Input validation and sanitization
- SQL injection prevention
- CORS configuration

## Monitoring

- Health check endpoint: `GET /health`
- Payment analytics and reporting
- Error logging and monitoring
- Database connection pooling

## Development

### Project Structure
```
payment-service/
├── cmd/api/                 # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── db/                 # Database connection (SQLx)
│   ├── domain/             # Domain models
│   ├── dto/                # Data transfer objects
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── repository/         # Data access layer (SQLx)
│   ├── router/             # HTTP routing
│   ├── services/           # Business logic
│   └── validation/         # Input validation
├── migrations/             # Database migrations
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── env.example            # Environment variables example
└── README.md              # This file
```

### SQLx Benefits

The repository layer uses **SQLx** for enhanced database operations:

#### Before (Standard database/sql):
```go
// Manual scanning and error handling
rows, err := db.QueryContext(ctx, query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()

var payments []*Payment
for rows.Next() {
    payment := &Payment{}
    err := rows.Scan(
        &payment.ID, &payment.OrderID, &payment.Amount,
        &payment.Currency, &payment.Status, &payment.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    payments = append(payments, payment)
}
```

#### After (SQLx):
```go
// Automatic struct scanning
var payments []*Payment
err := db.SelectContext(ctx, &payments, query, args...)
if err != nil {
    return nil, err
}
```

#### Named Queries:
```go
// Clean, readable SQL with named parameters
query := `
    INSERT INTO payments (order_id, payment_method_id, transaction_id, amount, currency, status)
    VALUES (:order_id, :payment_method_id, :transaction_id, :amount, :currency, :status)
    RETURNING id`

result, err := db.NamedQueryContext(ctx, query, payment)
```

#### Benefits:
- **Type Safety**: Compile-time checking of struct fields
- **Cleaner Code**: Less boilerplate for scanning results
- **Named Parameters**: More readable SQL queries
- **Better Performance**: Optimized scanning operations
- **Error Reduction**: Fewer manual error handling cases

### Adding New Payment Gateways

1. Implement the `GatewayProvider` interface
2. Add gateway configuration to config
3. Register the gateway in `PaymentGatewayService`
4. Add webhook handling logic

### Testing

Run tests with:
```bash
go test ./...
```

## License

This project is part of the Gearbox e-commerce platform.
