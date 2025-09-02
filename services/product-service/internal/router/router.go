package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/handlers"
)

func NewRouter(categoryHandler handlers.ICategoryHandler) *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60))

	// Global CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure this properly for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"product-service"}`))
	})

	// Product service routes
	router.Route("/api/v1", func(r chi.Router) {
		// Service health endpoint
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok","service":"product-service","version":"1.0"}`))
		})

		// Category routes
		r.Route("/categories", func(r chi.Router) {
			r.Post("/", categoryHandler.CreateCategory)
			r.Get("/", categoryHandler.ListCategories)
			r.Get("/hierarchy", categoryHandler.GetCategoryHierarchy)
			r.Get("/slug/{slug}", categoryHandler.GetCategoryBySlug)
			r.Get("/{id}", categoryHandler.GetCategory)
			r.Put("/{id}", categoryHandler.UpdateCategory)
			r.Delete("/{id}", categoryHandler.DeleteCategory)
			r.Get("/{id}/children", categoryHandler.GetCategoryChildren)
		})

		// TODO: Add product routes here
		// r.Route("/products", func(r chi.Router) {
		//     r.Post("/", productHandler.CreateProduct)
		//     r.Get("/", productHandler.ListProducts)
		//     r.Get("/{id}", productHandler.GetProduct)
		//     r.Put("/{id}", productHandler.UpdateProduct)
		//     r.Delete("/{id}", productHandler.DeleteProduct)
		// })
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
