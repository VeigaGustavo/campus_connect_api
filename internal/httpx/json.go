package httpx

import (
	"encoding/json"
	"net/http"

	"campus_connect_api/internal/models"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, models.APIError{Code: code, Message: message})
}
