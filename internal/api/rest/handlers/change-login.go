package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"users-profile-service/internal/user/management"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ChangeLoginRequest struct {
	Login string `json:"login"`
}

func ChangeLoginHandler(logger *zap.SugaredLogger, manager *management.UserManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 7*time.Second)
		defer cancel()

		var req ChangeLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			GenerateError(w, err.Error(), http.StatusBadRequest)
			return
		}

		userIDStr := chi.URLParam(r, "user_id")
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			GenerateError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = manager.EditLogin(ctx, userID, req.Login)
		if err != nil {
			HandleError(w, err, logger)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	}
}
