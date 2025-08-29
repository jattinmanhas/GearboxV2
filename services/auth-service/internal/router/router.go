package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/middleware"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
)

func NewRouter(authHandler handlers.IAuthHandler, authService services.IAuthService, roleHandler *handlers.RoleHandler) *chi.Mux {
	router := chi.NewRouter()

	// CORS middleware
	router.Use(middleware.CORSMiddleware([]string{"*"}))

	// Auth routes
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Auth Service is running"))
		})

		// Public authentication routes (no auth required)
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/refresh", authHandler.RefreshToken)
		r.Post("/logout", authHandler.Logout)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))

			// User management routes
			r.Get("/user/{id}", authHandler.GetUserByID)
			r.Put("/user/{id}", authHandler.UpdateUser)
			r.Delete("/user/{id}", authHandler.DeleteUser)
			r.Post("/user/{id}/change-password", authHandler.ChangePassword)
			r.Get("/users", authHandler.GetAllUsers)
			r.Post("/logout-all", authHandler.LogoutAll)

			// Role management routes
			r.Route("/roles", func(r chi.Router) {
				r.Get("/", roleHandler.GetAllRoles)                     // Get all available roles
				r.Get("/my-role", roleHandler.GetMyRole)                // Get current user's role
				r.Get("/user", roleHandler.GetUserRole)                 // Get specific user's role
				r.Post("/assign", roleHandler.AssignRoleToUser)         // Assign role to user
				r.Put("/update", roleHandler.UpdateUserRole)            // Update user's role
				r.Delete("/remove", roleHandler.RemoveUserRole)         // Remove user's role
				r.Get("/check-permission", roleHandler.CheckPermission) // Check if user has permission
			})
		})
	})

	return router
}
