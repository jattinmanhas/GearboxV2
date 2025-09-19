package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/handlers"
	"github.com/jattinmanhas/GearboxV2/services/shared/jwt"
	sharedMiddleware "github.com/jattinmanhas/GearboxV2/services/shared/middleware"
)

func NewRouter(categoryHandler handlers.ICategoryHandler, productHandler handlers.IProductHandler, cartHandler handlers.ICartHandler, inventoryHandler handlers.IInventoryHandler, couponHandler handlers.CouponHandler, orderHandler handlers.OrderHandler, jwtSecret, jwtRefreshSecret string) *chi.Mux {
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

			// Product variants
			r.Post("/{id}/variants", productHandler.CreateProductVariant)
			r.Get("/{id}/variants", productHandler.GetProductVariants)
			r.Get("/{id}/variants-with-inventory", productHandler.GetProductVariantsWithInventory)
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

		// Cart routes (supports both guest and authenticated users)
		r.Route("/carts", func(r chi.Router) {
			// Add optional authentication middleware to check for logged-in users
			jwtService := jwt.NewJWTService(jwtSecret, jwtRefreshSecret)
			authService := sharedMiddleware.NewSharedAuthService(jwtService)
			r.Use(sharedMiddleware.OptionalAuthMiddleware(authService))

			r.Get("/session", cartHandler.GetCartBySession)
			r.Get("/get-or-create", cartHandler.GetOrCreateCart)
			r.Get("/analytics", cartHandler.GetCartAnalytics)
			r.Get("/analytics/date-range", cartHandler.GetCartAnalyticsByDateRange)
			r.Get("/analytics/top-products", cartHandler.GetTopProductsInCarts)
			r.Get("/analytics/abandonment-rate", cartHandler.GetCartAbandonmentRate)
			r.Get("/analytics/conversion-funnel", cartHandler.GetCartConversionFunnel)
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

		// Wishlist routes (requires authentication - user-specific)
		r.Route("/wishlists", func(r chi.Router) {
			jwtService := jwt.NewJWTService(jwtSecret, jwtRefreshSecret)
			authService := sharedMiddleware.NewSharedAuthService(jwtService)
			r.Use(sharedMiddleware.AuthMiddleware(authService))
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

		// Inventory routes (requires editor/admin roles)
		r.Route("/inventory", func(r chi.Router) {
			jwtService := jwt.NewJWTService(jwtSecret, jwtRefreshSecret)
			authService := sharedMiddleware.NewSharedAuthService(jwtService)
			r.Use(sharedMiddleware.AuthMiddleware(authService))
			r.Use(sharedMiddleware.RequireEditor())
			r.Post("/", inventoryHandler.CreateInventory)
			r.Get("/summary", inventoryHandler.GetInventorySummary)
			r.Get("/", inventoryHandler.ListInventory)
			r.Get("/product", inventoryHandler.GetInventoryByProduct)
			r.Get("/{id}", inventoryHandler.GetInventoryByID)
			r.Put("/{id}", inventoryHandler.UpdateInventory)
			r.Delete("/{id}", inventoryHandler.DeleteInventory)

			// Stock movements
			r.Post("/movements", inventoryHandler.RecordStockMovement)
			r.Get("/movements", inventoryHandler.GetStockMovements)
			r.Get("/movements/{id}", inventoryHandler.GetStockMovementByID)

			// Stock reservations
			r.Post("/reservations", inventoryHandler.ReserveStock)
			r.Delete("/reservations", inventoryHandler.ReleaseStock)
			r.Get("/reservations", inventoryHandler.GetStockReservations)

			// Inventory alerts
			r.Get("/alerts", inventoryHandler.GetInventoryAlerts)
			r.Put("/alerts/{id}/resolve", inventoryHandler.ResolveInventoryAlert)
			r.Post("/alerts/check", inventoryHandler.CheckLowStockAlerts)

			// Bulk operations
			r.Post("/bulk-update", inventoryHandler.BulkUpdateStock)
		})

		// Coupon routes
		r.Route("/coupons", func(r chi.Router) {
			// Public routes (no authentication required)
			r.Post("/validate", couponHandler.ValidateCoupon)

			// Protected routes (requires editor/admin roles)
			r.Group(func(r chi.Router) {
				jwtService := jwt.NewJWTService(jwtSecret, jwtRefreshSecret)
				authService := sharedMiddleware.NewSharedAuthService(jwtService)
				r.Use(sharedMiddleware.AuthMiddleware(authService))
				r.Use(sharedMiddleware.RequireEditor())

				r.Post("/", couponHandler.CreateCoupon)
				r.Get("/", couponHandler.ListCoupons)
				r.Get("/{id}", couponHandler.GetCoupon)
				r.Put("/{id}", couponHandler.UpdateCoupon)
				r.Delete("/{id}", couponHandler.DeleteCoupon)

				// Coupon usage (admin/editor only)
				r.Get("/usage", couponHandler.GetCouponUsage)
			})
		})

		// Order routes (protected - requires authentication)
		r.Route("/orders", func(r chi.Router) {
			jwtService := jwt.NewJWTService(jwtSecret, jwtRefreshSecret)
			authService := sharedMiddleware.NewSharedAuthService(jwtService)
			r.Use(sharedMiddleware.AuthMiddleware(authService))
			r.Post("/", orderHandler.CreateOrder)
			r.Get("/", orderHandler.ListOrders)
			r.Get("/number/{orderNumber}", orderHandler.GetOrderByNumber)
			r.Get("/{id}", orderHandler.GetOrder)
			r.Put("/{id}", orderHandler.UpdateOrder)
			r.Delete("/{id}", orderHandler.DeleteOrder)

			// Order items
			r.Get("/{id}/items", orderHandler.GetOrderItems)
			r.Put("/items/{id}", orderHandler.UpdateOrderItem)
			r.Delete("/items/{id}", orderHandler.DeleteOrderItem)

			// Order addresses
			r.Get("/{id}/addresses", orderHandler.GetOrderAddresses)
			r.Put("/addresses/{id}", orderHandler.UpdateOrderAddress)

			// Order status management
			r.Put("/{id}/status", orderHandler.UpdateOrderStatus)
			r.Get("/{id}/status-history", orderHandler.GetOrderStatusHistory)

			// Order fulfillment
			r.Post("/{id}/fulfillment", orderHandler.CreateOrderFulfillment)
			r.Get("/{id}/fulfillment", orderHandler.GetOrderFulfillment)
			r.Put("/fulfillment/{id}", orderHandler.UpdateOrderFulfillment)

			// Order refunds
			r.Post("/{id}/refunds", orderHandler.CreateOrderRefund)
			r.Get("/{id}/refunds", orderHandler.GetOrderRefunds)

			// Order analytics
			r.Get("/analytics", orderHandler.GetOrderAnalytics)
			r.Get("/analytics/date-range", orderHandler.GetOrderAnalyticsByDateRange)
			r.Get("/analytics/top-products", orderHandler.GetTopSellingProducts)

			// Cart integration
			r.Post("/from-cart", orderHandler.CreateOrderFromCart)

			// Payment integration
			r.Post("/{id}/payments", orderHandler.CreateOrderPayment)
			r.Post("/{id}/payments/process", orderHandler.ProcessOrderPayment)
			r.Get("/{id}/payments", orderHandler.GetOrderPayment)
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
