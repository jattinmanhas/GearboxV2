package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/db"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/repository"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/router"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT service
	jwtService := jwt.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.RefreshSecret,
	)

	// Initialize repositories
	paymentRepo := repository.NewPaymentRepository(db.GetDB())

	// Initialize services
	gatewayService := services.NewPaymentGatewayService(paymentRepo, cfg)
	paymentService := services.NewPaymentService(paymentRepo, gatewayService)

	// Initialize router
	r := router.NewRouter(paymentService, jwtService)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Payment service starting on %s", cfg.GetServerAddress())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down payment service...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Payment service stopped")
}
