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

	cfg = &Config{
		Port:             port,
		DatabaseURL:      databaseURL,
		JWTSecret:        jwtSecret,
		JWTRefreshSecret: jwtRefreshSecret,
		Environment:      environment,
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
