package router

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	sharedMiddleware "github.com/jattinmanhas/GearboxV2/services/shared/middleware"
)

func NewRouter(authHandler handlers.IAuthHandler, authService services.IAuthService, roleHandler handlers.IRoleHandler, addressHandler handlers.IAddressHandler, oauthHandler handlers.IOAuthHandler, jwtService *jwt.JWTService) *chi.Mux {
	router := chi.NewRouter()

	// Global CORS middleware
	router.Use(sharedMiddleware.CORSMiddleware([]string{"*"}))
	router.Use(middleware.Logger)

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

		// OAuth routes (public)
		r.Get("/oauth/{provider}", oauthHandler.InitiateOAuth)
		r.Get("/oauth/{provider}/callback", oauthHandler.HandleOAuthCallback)

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			sharedAuthService := sharedMiddleware.NewSharedAuthService(jwtService)
			r.Use(sharedMiddleware.AuthMiddleware(sharedAuthService))

			// Authentication routes
			r.Post("/logout", authHandler.Logout)

			// User management routes
			r.Get("/user/{id}", authHandler.GetUserByID)
			r.Put("/user/{id}", authHandler.UpdateUser)
			r.Delete("/user/{id}", authHandler.DeleteUser)
			r.Post("/user/{id}/change-password", authHandler.ChangePassword)
			r.Post("/logout-all", authHandler.LogoutAll)
			r.Get("/users", authHandler.GetAllUsers)                // Temporarily moved outside admin group for testing
			r.Get("/users/analytics", authHandler.GetUserAnalytics) // User analytics

			// Profile routes
			r.Get("/profile", authHandler.GetProfile)
			r.Put("/profile", authHandler.UpdateProfile)

			// OAuth management routes (protected)
			r.Post("/oauth/link/{provider}", oauthHandler.LinkOAuthProvider)
			r.Delete("/oauth/unlink/{provider}", oauthHandler.UnlinkOAuthProvider)
			r.Get("/oauth/providers", oauthHandler.GetLinkedProviders)

			// Admin-only cleanup
			r.Group(func(r chi.Router) {
				r.Use(sharedMiddleware.RequireAdmin())
				r.Post("/cleanup-expired-tokens", authHandler.CleanupExpiredTokens)
			})

			// Role management routes
			r.Route("/roles", func(r chi.Router) {
				r.Get("/", roleHandler.GetAllRoles)      // Everyone can view roles
				r.Get("/my-role", roleHandler.GetMyRole) // Authenticated user gets their role
				r.Get("/user", roleHandler.GetUserRole)  // Authenticated user gets another user's role

				// Editor+ can assign roles
				r.Group(func(r chi.Router) {
					r.Use(sharedMiddleware.RequireEditor())
					r.Post("/assign", roleHandler.AssignRoleToUser)
					r.Delete("/remove", roleHandler.RemoveUserRole)
				})

				// Permission check route (just checks against required role passed in request)
				r.Get("/check-permission", roleHandler.CheckPermission)
			})

			// Address management routes
			r.Route("/addresses", func(r chi.Router) {
				// Address operations
				r.Post("/", addressHandler.CreateAddress)
				r.Get("/", addressHandler.GetAddresses)
				r.Get("/default", addressHandler.GetDefaultAddress)
				r.Get("/{id}", addressHandler.GetAddressByID)
				r.Put("/{id}", addressHandler.UpdateAddress)
				r.Delete("/{id}", addressHandler.DeleteAddress)
				r.Post("/set-default", addressHandler.SetDefaultAddress)

				// Phone number operations
				r.Route("/phones", func(r chi.Router) {
					r.Post("/", addressHandler.CreatePhoneNumber)
					r.Get("/", addressHandler.GetPhoneNumbers)
					r.Get("/primary", addressHandler.GetPrimaryPhone)
					r.Get("/{id}", addressHandler.GetPhoneNumberByID)
					r.Put("/{id}", addressHandler.UpdatePhoneNumber)
					r.Delete("/{id}", addressHandler.DeletePhoneNumber)
					r.Post("/set-primary", addressHandler.SetPrimaryPhone)
				})
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
