package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
)

func NewRouter(authHandler handlers.IAuthHandler) *chi.Mux {
	router := chi.NewRouter()


	// Auth routes
	router.Route("/api/v1/auth", func(r chi.Router){
		r.Get("/health", func(w http.ResponseWriter, r *http.Request){
			w.Write([]byte("Auth Service is running"))
		})

		r.Get("/user/{id}", authHandler.GetUserByID)
		r.Post("/register", authHandler.RegisterUser)
		r.Get("/users", authHandler.GetAllUsers)
	})
	
	return router
}
