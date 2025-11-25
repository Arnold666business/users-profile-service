package http

import (
	"context"
	"errors"
	"net/http"
)

func (s *Server) Start() {
	go func() {
		s.Logger.Debugf("starting http server on %s", s.Instance.Addr)
		if err := s.Instance.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Warnw("http server unexcepted error", "error", err)
		}
	}()
}

func (s *Server) Closer(ctx context.Context) error {
	s.Logger.Info("Shutting down base server...")
	return s.Instance.Shutdown(ctx)
}
