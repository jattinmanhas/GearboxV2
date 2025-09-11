package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the payment service
type Config struct {
	// Server configuration
	Server ServerConfig `json:"server"`

	// Database configuration
	Database DatabaseConfig `json:"database"`

	// JWT configuration
	JWT JWTConfig `json:"jwt"`

	// Payment gateway configuration
	PaymentGateways PaymentGatewayConfig `json:"payment_gateways"`

	// Service URLs
	ServiceURLs ServiceURLsConfig `json:"service_urls"`

	// Logging configuration
	Logging LoggingConfig `json:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret        string        `json:"secret"`
	RefreshSecret string        `json:"refresh_secret"`
}

// PaymentGatewayConfig holds payment gateway configuration
type PaymentGatewayConfig struct {
	Stripe   StripeConfig   `json:"stripe"`
	PayPal   PayPalConfig   `json:"paypal"`
	Razorpay RazorpayConfig `json:"razorpay"`
}

// StripeConfig holds Stripe-specific configuration
type StripeConfig struct {
	SecretKey      string `json:"secret_key"`
	PublishableKey string `json:"publishable_key"`
	WebhookSecret  string `json:"webhook_secret"`
}

// PayPalConfig holds PayPal-specific configuration
type PayPalConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Mode         string `json:"mode"` // sandbox or live
}

// RazorpayConfig holds Razorpay-specific configuration
type RazorpayConfig struct {
	KeyID         string `json:"key_id"`
	KeySecret     string `json:"key_secret"`
	WebhookSecret string `json:"webhook_secret"`
}

// ServiceURLsConfig holds URLs for other services
type ServiceURLsConfig struct {
	AuthService    string `json:"auth_service"`
	ProductService string `json:"product_service"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't return an error
		fmt.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	// Server configuration
	config.Server.Host = getEnv("HOST", "localhost")
	config.Server.Port = getEnv("PORT", "8083")

	// Database configuration
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnvAsInt("DB_PORT", 5432)
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "password")
	config.Database.Name = getEnv("DB_NAME", "gearbox_payment")
	config.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

	// JWT configuration
	config.JWT.Secret = getEnv("JWT_SECRET", "your-jwt-secret-key-here")
	config.JWT.RefreshSecret = getEnv("JWT_REFRESH_SECRET", "your-jwt-refresh-secret-key-here")

	// Payment gateway configuration
	config.PaymentGateways.Stripe.SecretKey = getEnv("STRIPE_SECRET_KEY", "")
	config.PaymentGateways.Stripe.PublishableKey = getEnv("STRIPE_PUBLISHABLE_KEY", "")
	config.PaymentGateways.Stripe.WebhookSecret = getEnv("STRIPE_WEBHOOK_SECRET", "")

	config.PaymentGateways.PayPal.ClientID = getEnv("PAYPAL_CLIENT_ID", "")
	config.PaymentGateways.PayPal.ClientSecret = getEnv("PAYPAL_CLIENT_SECRET", "")
	config.PaymentGateways.PayPal.Mode = getEnv("PAYPAL_MODE", "sandbox")

	config.PaymentGateways.Razorpay.KeyID = getEnv("RAZORPAY_KEY_ID", "")
	config.PaymentGateways.Razorpay.KeySecret = getEnv("RAZORPAY_KEY_SECRET", "")
	config.PaymentGateways.Razorpay.WebhookSecret = getEnv("RAZORPAY_WEBHOOK_SECRET", "")

	// Service URLs
	config.ServiceURLs.AuthService = getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	config.ServiceURLs.ProductService = getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082")

	// Logging configuration
	config.Logging.Level = getEnv("LOG_LEVEL", "info")
	config.Logging.Format = getEnv("LOG_FORMAT", "json")

	return config, nil
}

// GetDatabaseDSN returns the database connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}