package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router   *chi.Mux
	Instance *http.Server
	Logger   *zap.SugaredLogger
	DB       *pgxpool.Pool
}

func Build(logger *zap.SugaredLogger, router *chi.Mux) *Server {
	l := logger.Named("http.server")
	return &Server{
		Router: router,
		Instance: &http.Server{
			Addr:    fmt.Sprintf(":%d", os.Getenv("SERVER_PORT")),
			Handler: router,
		},
		Logger: l,
	}
}
