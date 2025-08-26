package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

var (
	cfg  *Config
	once sync.Once // Ensures the config is loaded only once.
)

func LoadConfig() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = "8081"
		}

		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("Error: DATABASE_URL environment variable is not set.")
		}

		cfg = &Config{
			Port:        port,
			DatabaseURL: databaseURL,
		}
	})
	return cfg
}
