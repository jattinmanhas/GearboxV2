package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/db"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/repository"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/router"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	database, err := db.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	migrationsPath := "migrations"
	if err := database.RunMigrations(migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	cartRepo := repository.NewCartRepository(database.DB)
	inventoryRepo := repository.NewInventoryRepository(database.DB)
	couponRepo := repository.NewCouponRepository(database.DB)
	orderRepo := repository.NewOrderRepository(database.DB)

	// Initialize services
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	productService := services.NewProductService(productRepo)
	inventoryService := services.NewInventoryService(inventoryRepo, productRepo)
	couponService := services.NewCouponService(couponRepo, cartRepo)
	cartService := services.NewCartService(cartRepo, productRepo, inventoryService, couponService)

	// Initialize auth client for user address management
	authClient := services.NewAuthServiceClient(cfg.AuthServiceURL)

	// Initialize payment client for payment processing
	paymentClient := services.NewPaymentServiceClient(cfg.PaymentServiceURL)

	orderService := services.NewOrderService(orderRepo, productRepo, inventoryService, cartService, couponService, authClient, paymentClient)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService, cfg.SecureCookies)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	couponHandler := handlers.NewCouponHandler(couponService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Initialize router
	appRouter := router.NewRouter(categoryHandler, productHandler, cartHandler, inventoryHandler, couponHandler, orderHandler, cfg.JWTSecret, cfg.JWTRefreshSecret)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      appRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", server.Addr)
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
