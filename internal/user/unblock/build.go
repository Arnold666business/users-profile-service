package unblock

import (
	"users-profile-service/internal/external/kafka/producer/UnBlockUser"
	"users-profile-service/internal/external/kafka/producer/build"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type UnblockUserProcessor struct {
	logger                    *zap.SugaredLogger
	userBlockStatusRepository *repository.UserBlockStatus
	userHistoryRepository     *repository.UserHistory
	userRepository            *repository.User
	transactor                *repository.Transactor
	redis                     *redis.RedisProvider
	unblockUserProducer       *UnBlockUser.Producer
}

func Build(logger *zap.SugaredLogger,
	repositories *repository.Repositories,
	redis *redis.RedisProvider,
	producers *build.Producers) *UnblockUserProcessor {
	return &UnblockUserProcessor{
		logger:                    logger,
		userBlockStatusRepository: repositories.UserBlockStatusRepository,
		userHistoryRepository:     repositories.UserHistoryRepository,
		userRepository:            repositories.UserRepository,
		transactor:                repositories.Transactor,
		redis:                     redis,
		unblockUserProducer:       producers.UnBlockUser,
	}
}
