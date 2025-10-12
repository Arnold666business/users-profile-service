package create

import (
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type CreateUserProcessor struct {
	logger          *zap.SugaredLogger
	uhRepository    *repository.UserHistory
	userRepository  *repository.User
	newUserProducer *NewUser.Producer
	redis           *redis.RedisProvider
}

func Build(
	logger *zap.SugaredLogger,
	uhRepository *repository.UserHistory,
	userRepository *repository.User,
	newUserProducer *NewUser.Producer,
	redis *redis.RedisProvider,
) (*CreateUserProcessor, error) {
	return &CreateUserProcessor{
		logger:          logger,
		uhRepository:    uhRepository,
		userRepository:  userRepository,
		newUserProducer: newUserProducer,
		redis:           redis,
	}, nil
}
