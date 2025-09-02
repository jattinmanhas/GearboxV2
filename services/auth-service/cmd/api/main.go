package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	router "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/router"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/db"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Println("Configuration loaded successfully.")

	// Connect to DB
	database, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	}
	defer database.Close()
	log.Println("✅ Database connection established successfully.")

	// Run migrations
	migrationsPath := "migrations"
	if err := db.RunMigrations(database, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	refreshTokenRepo := repository.NewRefreshTokenRepository(database)
	roleRepo := repository.NewRoleRepository(database)

	// Initialize services
	jwtService := services.NewJWTService(cfg.JWTSecret, cfg.JWTRefreshSecret)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, roleRepo, jwtService)
	userService := services.NewUserService(userRepo, authService)
	roleService := services.NewRoleService(roleRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, authService, jwtService)
	roleHandler := handlers.NewRoleHandler(roleService)

	// Initialize router
	appRouter := router.NewRouter(authHandler, authService, roleHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      appRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Auth Service Running on port: %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
