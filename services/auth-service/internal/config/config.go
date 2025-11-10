package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	JWTRefreshSecret string
	Environment      string

	// OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURL  string
	FrontendURL        string

	// Email Configuration
	EmailSMTPHost     string
	EmailSMTPPort     string
	EmailSMTPUser     string
	EmailSMTPPassword string
	EmailFromAddress  string
	EmailFromName     string
	EmailProvider     string // "smtp", "sendgrid", "mailgun", "console" (for logging)
}

var (
	cfg  *Config
	once sync.Once  // Ensures the config is loaded only once.
	mu   sync.Mutex // Protects the reset functionality
)

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	once.Do(func() {
		loadConfig()
	})
	return cfg
}

// loadConfig is the internal function that actually loads the config
func loadConfig() {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Error: JWT_SECRET environment variable is not set.")
	}

	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtRefreshSecret == "" {
		log.Fatal("Error: JWT_REFRESH_SECRET environment variable is not set.")
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	// OAuth Configuration
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectURL := os.Getenv("GITHUB_REDIRECT_URL")

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Email Configuration
	emailSMTPHost := os.Getenv("EMAIL_SMTP_HOST")
	emailSMTPPort := os.Getenv("EMAIL_SMTP_PORT")
	if emailSMTPPort == "" {
		emailSMTPPort = "587" // Default to TLS port
	}
	emailSMTPUser := os.Getenv("EMAIL_SMTP_USER")
	emailSMTPPassword := os.Getenv("EMAIL_SMTP_PASSWORD")
	emailFromAddress := os.Getenv("EMAIL_FROM_ADDRESS")
	if emailFromAddress == "" {
		emailFromAddress = "noreply@gearbox.com"
	}
	emailFromName := os.Getenv("EMAIL_FROM_NAME")
	if emailFromName == "" {
		emailFromName = "GearBox"
	}
	emailProvider := os.Getenv("EMAIL_PROVIDER")
	if emailProvider == "" {
		emailProvider = "console" // Default to console logging in development
	}

	cfg = &Config{
		Port:             port,
		DatabaseURL:      databaseURL,
		JWTSecret:        jwtSecret,
		JWTRefreshSecret: jwtRefreshSecret,
		Environment:      environment,

		// OAuth
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectURL:  googleRedirectURL,
		GithubClientID:     githubClientID,
		GithubClientSecret: githubClientSecret,
		GithubRedirectURL:  githubRedirectURL,
		FrontendURL:        frontendURL,

		// Email
		EmailSMTPHost:     emailSMTPHost,
		EmailSMTPPort:     emailSMTPPort,
		EmailSMTPUser:     emailSMTPUser,
		EmailSMTPPassword: emailSMTPPassword,
		EmailFromAddress:  emailFromAddress,
		EmailFromName:     emailFromName,
		EmailProvider:     emailProvider,
	}
}

// ResetConfig resets the singleton for testing purposes
// This should only be used in tests
func ResetConfig() {
	mu.Lock()
	defer mu.Unlock()
	cfg = nil
	once = sync.Once{}
}
