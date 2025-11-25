package middleware

import (
	"net/http"
	"os"
	"users-profile-service/internal/api/rest/handlers"

	"go.uber.org/zap"
)

func ValidateTokenMiddleware(l *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const authHeader = "x-Token"
			reqToken := r.Header.Get(authHeader)
			if reqToken == "" {
				l.Debugw("Token not found in header", "header", authHeader, "path", r.URL.Path, "method", r.Method)
				handlers.GenerateError(w, "not found x-Token", http.StatusUnauthorized)
				return
			}

			validToken := os.Getenv("SPERMA")
			if reqToken != validToken {
				l.Debugw("not valid token", "header", authHeader, "path", r.URL.Path, "method", r.Method)
				handlers.GenerateError(w, "not found x-Token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
