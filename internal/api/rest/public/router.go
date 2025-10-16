package publiс

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.SugaredLogger, db *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Route("/api/{version}", func(r chi.Router) {
		r.Use(xTokenMiddleware)

		r.Route("/user", func(r chi.Router) {
			r.Post("", createUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("", getUserHandler)
				r.Put("/verify-email", verifyEmailHandler)
				r.Post("/email", changeEmailHandler)
				r.Post("/login", changeLoginHandler)
				r.Delete("/", deleteUserHandler)
			})
		})
	})
	return r
}
