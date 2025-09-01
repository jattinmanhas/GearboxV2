package main

import (
	"log"
	"net/http"

	router "github.com/jattinmanhas/GearboxV2/services/auth-service/internal/router"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/config"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/db"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repository"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
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
	r := router.NewRouter(authHandler, authService, roleHandler)

	log.Println("🚀 Auth Service Running on port:", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
