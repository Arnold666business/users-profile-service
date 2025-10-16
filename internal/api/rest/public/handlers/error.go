package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	Status       int    `json:"status"`
}

func GenerateError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: message, Status: code})
}
