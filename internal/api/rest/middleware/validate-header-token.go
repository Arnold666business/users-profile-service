package middleware

import "net/http"

func ValidateTokenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

	}
}
