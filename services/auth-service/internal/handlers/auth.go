package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/models"
	"github.com/jattinmanhas/GearboxV2/services/auth-service/internal/repositories"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	MiddleName  string    `json:"middle_name"`
	LastName    string    `json:"last_name"`
	Avatar      string    `json:"avatar"`
	Gender      string    `json:"gender"`
	DateOfBirth *time.Time `json:"date_of_birth"`
}

type apiError struct {
	Error string `json:"error"`
}

type apiOk struct {
	ID uint `json:"id"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Register(db *sqlx.DB) http.HandlerFunc {
	userRepo := repositories.NewUserRepository(db)
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "invalid JSON"})
			return
		}

		// basic validation
		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)
		if req.Username == "" || req.Email == "" || req.Password == "" {
			writeJSON(w, http.StatusBadRequest, apiError{Error: "username, email and password are required"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, apiError{Error: "failed to hash password"})
			return
		}

		user := models.User{
			Username:   req.Username,
			Password:   string(hash),
			Email:      req.Email,
			FirstName:  req.FirstName,
			MiddleName: req.MiddleName,
			LastName:   req.LastName,
			Avatar:     req.Avatar,
			Gender:     req.Gender,
		}
		if req.DateOfBirth != nil {
			user.DateOfBirth = *req.DateOfBirth
		}

		if err := userRepo.Create(r.Context(), &user); err != nil {
			// likely unique constraint violation on username/email
			writeJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusCreated, apiOk{ID: user.ID})
	}
}
