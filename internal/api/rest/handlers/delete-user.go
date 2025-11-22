package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"users-profile-service/internal/user/delete"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func DeleteUserHandler(logger *zap.SugaredLogger, processor *delete.DeleteUserProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
		defer cancel()

		userIDStr := chi.URLParam(r, "user_id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			GenerateError(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = processor.Process(ctx, userID)
		if err != nil {
			HandleError(w, err, logger)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
