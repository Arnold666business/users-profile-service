package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"users-profile-service/internal/api/rest"
	"users-profile-service/internal/api/rest/http"
	"users-profile-service/internal/api/rpc"
	"users-profile-service/internal/api/rpc/users"
	"users-profile-service/internal/external/kafka/producer"
	"users-profile-service/internal/logger"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/postgres"
	"users-profile-service/internal/repository/redis"
	"users-profile-service/internal/user/block"
	"users-profile-service/internal/user/create"
	"users-profile-service/internal/user/delete"
	"users-profile-service/internal/user/management"
)

type Stopper interface {
	stop()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	l, err := logger.Build()
	if err != nil {
		panic("logger init failed" + err.Error())
	}

	db, err := postgres.Build(l)
	if err != nil {
		panic("logger init postgres" + err.Error())
	}
	defer db.Close()

	repos := repository.Build(db)

	redisProvider, err := redis.Build(l)
	if err != nil {
		l.Error("redis init failed" + err.Error())
	}

	producers, err := producer.Build(l)
	if err != nil {
		panic("producer build failed" + err.Error())
	}

	грпс

	deleteUserService := delete.Build(l, repos, грпс, producers, redisProvider)
	createUserService := create.Build(l, repos, producers, redisProvider)
	blockUserService := block.Build(l, repos, redisProvider)

	лисенер

	userManagement := management.Build(l, repos, redisProvider)

	router := rest.BuildRouter(l, createUserService, deleteUserService, userManagement)
	baseServer := http.Build(l, router)
	go baseServer.Start()

	usersServiceImpl := users.Build(l, userManagement)
	grpcServer := rpc.Build(l, usersServiceImpl)
	go grpcServer.Start()

	<-ctx.Done()
	l.Info("Shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	baseServer.Stop(shutdownCtx)

}
