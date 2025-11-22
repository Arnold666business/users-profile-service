package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	common_error "users-profile-service/internal/common-error"

	"go.uber.org/zap"
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

func HandleError(w http.ResponseWriter, err error, logger *zap.SugaredLogger) {
	var appErr *common_error.ErrorDefinition
	if !errors.As(err, &appErr) {
		logger.Error("Unknown error", zap.Error(err))
		GenerateError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch appErr.Type {
	case common_error.TypeInternal:
		GenerateError(w, err.Error(), http.StatusInternalServerError)
		return
	case common_error.TypeNotFound:
		GenerateError(w, err.Error(), http.StatusNotFound)
		return
	case common_error.TypeUnauthorized:
		GenerateError(w, err.Error(), http.StatusUnauthorized)
		return
	case common_error.TypeValidation:
		GenerateError(w, err.Error(), http.StatusBadRequest)
		return
	default:
		logger.Error("Unknown error type", zap.Error(err))
		GenerateError(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
