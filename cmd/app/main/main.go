package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"users-profile-service/internal/app"
	"users-profile-service/internal/logger"
	"users-profile-service/pkg/close"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	l, err := logger.Build()
	if err != nil {
		panic("logger init failed" + err.Error())
	}

	closer := &close.Closer{}

	readyApp, err := app.Build(closer, l)
	if err != nil {
		panic("app build failed" + err.Error())
	}

	readyApp.Run(ctx)

	<-ctx.Done()
	l.Info("Shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	closer.CloseAll(shutdownCtx)
}
