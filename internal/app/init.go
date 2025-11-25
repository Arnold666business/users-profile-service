package app

import (
	"context"
	"errors"
	"users-profile-service/internal/api/rest"
	"users-profile-service/internal/api/rest/http"
	"users-profile-service/internal/api/rpc"
	"users-profile-service/internal/api/rpc/users"
	"users-profile-service/internal/external/kafka/producer/build"
	"users-profile-service/internal/job"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/postgres"
	"users-profile-service/internal/repository/redis"
	"users-profile-service/internal/user/block"
	"users-profile-service/internal/user/create"
	delete_user "users-profile-service/internal/user/delete"
	"users-profile-service/internal/user/management"
	"users-profile-service/internal/user/unblock"
	"users-profile-service/pkg/close"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	logger                *zap.SugaredLogger
	db                    *pgxpool.Pool
	repositories          *repository.Repositories
	producers             *build.Producers
	redis                 *redis.RedisProvider
	deleteUserProcessor   *delete_user.DeleteUserProcessor
	creteUsersProcessor   *create.CreateUserProcessor
	blockUsersProcessor   *block.BlockUserProcessor
	unblockUsersProcessor *unblock.UnblockUserProcessor
	unblockUserJob        *job.CheckUsersForUnblockJob
	baseServer            *http.Server
	grpcServer            *rpc.GrpcServer
}

func Build(closer *close.Closer, l *zap.SugaredLogger) (*App, error) {

	db, err := postgres.Build(l)
	if err != nil {
		return nil, errors.New("postgres init failed" + err.Error())
	}
	closer.Add(db)

	repos := repository.Build(db)

	redisProvider, err := redis.Build(l)

	if err != nil {
		return nil, errors.New("redis init failed" + err.Error())
	}
	closer.Add(redisProvider.Close)

	producers, err := build.Build(l)
	if err != nil {
		return nil, errors.New("producer build failed" + err.Error())
	}

	//KitchenService

	deleteUserProcessor := delete_user.Build(l, repos, KitchenService, producers, redisProvider)
	createUserProcessor := create.Build(l, repos, producers, redisProvider)
	blockUserProcessor := block.Build(l, repos, redisProvider)
	unblockUserProcessor := unblock.Build(l, repos, redisProvider, producers)

	unblockUserJob := job.Build(l, unblockUserProcessor)

	//консумеры  blockUserProcessor use

	userManagement := management.Build(l, repos, redisProvider)

	router := rest.BuildRouter(l, createUserProcessor, deleteUserProcessor, userManagement)
	baseServer := http.Build(l, router)
	closer.Add(baseServer.Closer)

	usersServiceImpl := users.Build(l, userManagement)
	grpcServer := rpc.Build(l, usersServiceImpl)
	closer.Add(grpcServer.Close)

	return &App{
		db:                    db,
		repositories:          repos,
		producers:             producers,
		redis:                 redisProvider,
		deleteUserProcessor:   deleteUserProcessor,
		creteUsersProcessor:   createUserProcessor,
		blockUsersProcessor:   blockUserProcessor,
		unblockUsersProcessor: unblockUserProcessor,
		unblockUserJob:        unblockUserJob,
		baseServer:            baseServer,
		grpcServer:            grpcServer,
	}, nil
}

func (app *App) Run(ctx context.Context) {
	go app.baseServer.Start()
	go app.grpcServer.Start()

	//consumer

	go app.unblockUserJob.Run(ctx)
}
