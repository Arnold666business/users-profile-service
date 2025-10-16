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
	transactor      *repository.Transactor
}

func Build(
	logger *zap.SugaredLogger,
	uhRepository *repository.UserHistory,
	userRepository *repository.User,
	newUserProducer *NewUser.Producer,
	redis *redis.RedisProvider,
	transactor *repository.Transactor,
) (*CreateUserProcessor, error) {
	return &CreateUserProcessor{
		logger:          logger,
		uhRepository:    uhRepository,
		userRepository:  userRepository,
		newUserProducer: newUserProducer,
		redis:           redis,
		transactor:      transactor,
	}, nil
}
