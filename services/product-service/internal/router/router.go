package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/handlers"
)

func NewRouter(categoryHandler handlers.ICategoryHandler, productHandler handlers.IProductHandler, cartHandler handlers.ICartHandler) *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

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

		// Product routes
		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.CreateProduct)
			r.Get("/", productHandler.ListProducts)
			r.Get("/search", productHandler.SearchProducts)
			r.Get("/tags", productHandler.GetProductsByTags)
			r.Get("/sku/{sku}", productHandler.GetProductBySKU)
			r.Get("/{id}", productHandler.GetProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
			r.Patch("/{id}/quantity", productHandler.UpdateProductQuantity)

			// Product variants
			r.Post("/{id}/variants", productHandler.CreateProductVariant)
			r.Get("/{id}/variants", productHandler.GetProductVariants)
			r.Put("/variants/{id}", productHandler.UpdateProductVariant)
			r.Delete("/variants/{id}", productHandler.DeleteProductVariant)
			r.Get("/variants/{id}", productHandler.GetProductVariant)

			// Product categories
			r.Post("/{id}/categories", productHandler.AddProductToCategory)
			r.Put("/{id}/categories", productHandler.UpdateProductCategories)
			r.Get("/{id}/categories", productHandler.GetProductCategories)
			r.Delete("/{id}/categories/{category_id}", productHandler.RemoveProductFromCategory)
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
			r.Get("/{id}/products", productHandler.GetProductsByCategory)
		})

		// Cart routes
		r.Route("/carts", func(r chi.Router) {
			r.Post("/", cartHandler.CreateCart)
			r.Get("/get-or-create", cartHandler.GetOrCreateCart)
			r.Get("/analytics", cartHandler.GetCartAnalytics)
			r.Get("/{id}", cartHandler.GetCart)
			r.Put("/{id}", cartHandler.UpdateCart)
			r.Delete("/{id}", cartHandler.DeleteCart)

			// Cart items
			r.Post("/{id}/items", cartHandler.AddItemToCart)
			r.Get("/{id}/items", cartHandler.GetCartItems)
			r.Delete("/{id}/items", cartHandler.ClearCartItems)
			r.Get("/items/{id}", cartHandler.GetCartItem)
			r.Put("/items/{id}", cartHandler.UpdateCartItem)
			r.Delete("/items/{id}", cartHandler.DeleteCartItem)

			// Cart summary & calculations
			r.Get("/{id}/summary", cartHandler.GetCartSummary)
			r.Get("/{id}/total", cartHandler.GetCartTotal)
			r.Get("/{id}/count", cartHandler.GetCartItemCount)

			// Cart coupons
			r.Post("/{id}/coupons", cartHandler.ApplyCouponToCart)
			r.Get("/{id}/coupons", cartHandler.GetCartCoupons)
			r.Delete("/{id}/coupons", cartHandler.RemoveCouponFromCart)

			// Cart shipping
			r.Post("/{id}/shipping", cartHandler.SetCartShipping)
			r.Get("/{id}/shipping", cartHandler.GetCartShipping)
			r.Put("/{id}/shipping", cartHandler.UpdateCartShipping)
			r.Delete("/{id}/shipping", cartHandler.DeleteCartShipping)

			// Cart operations
			r.Post("/{id}/merge", cartHandler.MergeCarts)
			r.Delete("/{id}/clear", cartHandler.ClearCart)
		})

		// Wishlist routes
		r.Route("/wishlists", func(r chi.Router) {
			r.Post("/", cartHandler.CreateWishlist)
			r.Get("/", cartHandler.GetWishlists)
			r.Get("/{id}", cartHandler.GetWishlist)
			r.Put("/{id}", cartHandler.UpdateWishlist)
			r.Delete("/{id}", cartHandler.DeleteWishlist)

			// Wishlist items
			r.Post("/{id}/items", cartHandler.AddItemToWishlist)
			r.Get("/{id}/items", cartHandler.GetWishlistItems)
			r.Get("/items/{id}", cartHandler.GetWishlistItem)
			r.Put("/items/{id}", cartHandler.UpdateWishlistItem)
			r.Delete("/items/{id}", cartHandler.DeleteWishlistItem)
			r.Post("/items/{id}/move-to-cart", cartHandler.MoveItemToCart)
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
