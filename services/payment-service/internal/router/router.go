package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	sharedMiddleware "github.com/jattinmanhas/GearboxV2/services/shared/middleware"
)

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(paymentService *services.PaymentService, jwtService *jwt.JWTService) *chi.Mux {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * 1000000000)) // 60 seconds

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment service is healthy"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Initialize handlers
		paymentHandler := handlers.NewPaymentHandler(paymentService)
		authService := sharedMiddleware.NewSharedAuthService(jwtService)

		// Public routes (no authentication required)
		r.Route("/public", func(r chi.Router) {

			// Webhook endpoints (no auth required as they come from external services)
			r.Post("/webhooks/{gatewayId}", paymentHandler.HandleWebhook)
		})

		// Protected routes (authentication required)
		r.Route("/protected", func(r chi.Router) {
			r.Use(sharedMiddleware.AuthMiddleware(authService))

			// Payment operations (authenticated users)
			r.Route("/payments", func(r chi.Router) {
				r.Post("/", paymentHandler.CreatePayment)
				r.Post("/process", paymentHandler.ProcessPayment)
				r.Get("/order/{orderId}", paymentHandler.GetPaymentByOrderID)
				r.Get("/transaction/{transactionId}", paymentHandler.GetPaymentByTransactionID)
				r.Get("/{id}", paymentHandler.GetPayment)
			})
		})

		// Admin routes (admin only)
		r.Route("/admin", func(r chi.Router) {
			r.Use(sharedMiddleware.AuthMiddleware(authService))
			r.Use(sharedMiddleware.RequireAdmin())

			// Payment analytics and management
			r.Get("/payments", paymentHandler.ListPayments)
			r.Get("/payments/summary", paymentHandler.GetPaymentSummary)
			r.Get("/payments/analytics", paymentHandler.GetPaymentSummary)
			r.Put("/payments/{id}/status", paymentHandler.UpdatePaymentStatus)
			r.Post("/payments/refund", paymentHandler.RefundPayment)
			r.Get("/payments/export", paymentHandler.ListPayments) // TODO: Add export functionality
		})
	})

	return r
}
