package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/middleware"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
)

func NewRouter(authHandler handlers.IAuthHandler, authService services.IAuthService, roleHandler handlers.IRoleHandler) *chi.Mux {
	router := chi.NewRouter()

	// Global CORS middleware
	router.Use(middleware.CORSMiddleware([]string{"*"}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"auth-service"}`))
	})

	// Auth routes
	router.Route("/api/v1/auth", func(r chi.Router) {
		// Service health endpoint
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok","service":"auth-service","version":"1.0"}`))
		})

		// Public authentication routes (no auth required)
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/refresh", authHandler.RefreshToken)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))

			// Authentication routes
			r.Post("/logout", authHandler.Logout)

			// User management routes
			r.Get("/user/{id}", authHandler.GetUserByID)
			r.Put("/user/{id}", authHandler.UpdateUser)
			r.Delete("/user/{id}", authHandler.DeleteUser)
			r.Post("/user/{id}/change-password", authHandler.ChangePassword)
			r.Post("/logout-all", authHandler.LogoutAll)

			// Admin-only user listing + cleanup
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin())
				r.Get("/users", authHandler.GetAllUsers)
				r.Post("/cleanup-expired-tokens", authHandler.CleanupExpiredTokens)
			})

			// Role management routes
			r.Route("/roles", func(r chi.Router) {
				r.Get("/", roleHandler.GetAllRoles)      // Everyone can view roles
				r.Get("/my-role", roleHandler.GetMyRole) // Authenticated user gets their role
				r.Get("/user", roleHandler.GetUserRole)  // Authenticated user gets another user's role

				// Editor+ can assign roles
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireEditor())
					r.Post("/assign", roleHandler.AssignRoleToUser)
					r.Delete("/remove", roleHandler.RemoveUserRole)
				})

				// Permission check route (just checks against required role passed in request)
				r.Get("/check-permission", roleHandler.CheckPermission)
			})
		})
	})

	// 404 handler for unmatched routes
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found","message":"The requested resource was not found"}`))
	})

	// Method not allowed handler
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"method not allowed","message":"The requested method is not allowed for this resource"}`))
	})

	return router
}
