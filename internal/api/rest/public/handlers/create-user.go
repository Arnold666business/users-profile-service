package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"users-profile-service/internal/user/create"
)

type Request struct {
	Login          string `json:"login"`
	Email          string `json:"email"`
	Role           int    `json:"role"`
	IdempotencyKey string `json:"idempotency_key"`
}

type Response struct {
	UserId int `json:"user_id"`
}

func CreateUserHandler(processor *create.CreateUserProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
		defer cancel()

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			GenerateError(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := processor.Process(ctx, create.CreateRequest{
			Login:          req.Login,
			Email:          req.Email,
			Role:           req.Role,
			IdempotencyKey: req.IdempotencyKey,
		})
		if err != nil {

		}

	}
}
