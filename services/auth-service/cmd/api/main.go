package main

import (
	"log"
	"net/http"

	router "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/router"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/db"
    "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
    "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
    "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
)

func main() {
    cfg := config.LoadConfig()
    log.Println("Configuration loaded successfully.")

    // Connect to DB
    database, err := db.NewConnection(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("❌ Failed to connect to the database: %v", err)
    }
    defer database.Close()
    log.Println("✅ Database connection established successfully.")

    // Initialize repository
    userRepo := repository.NewUserRepository(database)

    // Initialize service
    userService := services.NewUserService(userRepo)

    // Initialize handler
    authHandler := handlers.NewAuthHandler(userService)

    // Initialize router
    r := router.NewRouter(authHandler)

    log.Println("🚀 Auth Service Running on port:", cfg.Port)
    log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}