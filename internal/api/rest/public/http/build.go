package http

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

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

func (s *Server) Start() {
	go func() {
		s.Logger.Debugf("starting http server on %s", s.Instance.Addr)
		if err := s.Instance.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Warnw("http server unexcepted error", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.stop()
}

func (s *Server) stop() {
	s.Instance.Close()
	if s.DB != nil {
		s.DB.Close()
	}

}
