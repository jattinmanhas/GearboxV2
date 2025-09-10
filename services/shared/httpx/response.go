package httpx

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Timestamp time.Time   `json:"timestamp"`
	Status    int         `json:"status"`
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, success bool, message string, data interface{}, errPayload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(APIResponse{
		Timestamp: time.Now().UTC(),
		Status:    status,
		Success:   success,
		Message:   message,
		Data:      data,
		Error:     errPayload,
	})
}

func OK(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusOK, true, message, data, nil)
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusCreated, true, message, data, nil)
}

func Error(w http.ResponseWriter, status int, message string, err error) {
	payload := map[string]string{"message": message}
	if err != nil {
		payload["detail"] = err.Error()
	}
	WriteJSON(w, status, false, message, nil, payload)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusBadRequest, false, message, data, nil)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusInternalServerError, false, message, data, nil)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusNotFound, false, message, data, nil)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusUnauthorized, false, message, data, nil)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(w http.ResponseWriter, message string, data interface{}) {
	WriteJSON(w, http.StatusForbidden, false, message, data, nil)
}
