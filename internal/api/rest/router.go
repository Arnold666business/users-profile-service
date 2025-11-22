package rest

import (
	"users-profile-service/internal/api/rest/handlers"
	"users-profile-service/internal/user/create"
	"users-profile-service/internal/user/delete"
	"users-profile-service/internal/user/management"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func BuildRouter(logger *zap.SugaredLogger,
	createProcessor *create.CreateUserProcessor,
	deleteProcessor *delete.DeleteUserProcessor,
	userManager *management.UserManager,
) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/{version}", func(r chi.Router) {
		r.Use(xTokenMiddleware)

		r.Route("/users", func(r chi.Router) {
			r.Post("", handlers.CreateUserHandler(logger, createProcessor))
			r.Route("/{userID}", func(r chi.Router) {
				r.Put("/verify-email", handlers.VerifyEmailHandler(logger, userManager))
				r.Post("/email", handlers.ChangeEmailHandler(logger, userManager))
				r.Post("/login", handlers.ChangeLoginHandler(logger, userManager))
				r.Delete("/", handlers.DeleteUserHandler(logger, deleteProcessor))
			})
		})
	})
	return r
}
