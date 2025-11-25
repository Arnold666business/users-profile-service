package create

import (
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/external/kafka/producer/build"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type CreateUserProcessor struct {
	logger                *zap.SugaredLogger
	userHistoryRepository *repository.UserHistory
	userRepository        *repository.User
	newUserProducer       *NewUser.Producer
	redis                 *redis.RedisProvider
	transactor            *repository.Transactor
}

func Build(
	logger *zap.SugaredLogger,
	repositories *repository.Repositories,
	producers *build.Producers,
	redis *redis.RedisProvider,
) *CreateUserProcessor {
	return &CreateUserProcessor{
		logger:                logger,
		userHistoryRepository: repositories.UserHistoryRepository,
		userRepository:        repositories.UserRepository,
		newUserProducer:       producers.NewUser,
		redis:                 redis,
		transactor:            repositories.Transactor,
	}
}
