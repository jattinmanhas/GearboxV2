package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/handlers"
)

func NewRouter(db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Auth Service is running"))
	})

	// Auth routes
	router.Post("/api/v1/auth/register", handlers.Register(db))
	
	return router
}
