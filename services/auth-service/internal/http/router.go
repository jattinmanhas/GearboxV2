package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Auth Service is running"))
	})
	
	return router
}
