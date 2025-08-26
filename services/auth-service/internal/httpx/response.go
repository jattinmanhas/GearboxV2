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
