package delete

import (
	"users-profile-service/internal/external/kafka/producer/DeletedUser"
	"users-profile-service/internal/external/kafka/producer/build"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type DeleteUserProcessor struct {
	logger                *zap.SugaredLogger
	userRepository        *repository.User
	ubsRepository         *repository.UserBlockStatus
	btdRepository         *repository.BlockTypeDictionary
	userHistoryRepository *repository.UserHistory
	kitchenService        *KitchenService
	deleteUserProducer    *DeletedUser.Producer
	transactor            *repository.Transactor
	redis                 *redis.RedisProvider
}

func Build(logger *zap.SugaredLogger,
	repositories *repository.Repositories,
	kitchenService *KitchenService,
	producers *build.Producers,
	redis *redis.RedisProvider) *DeleteUserProcessor {
	return &DeleteUserProcessor{
		logger:                logger,
		userRepository:        repositories.UserRepository,
		ubsRepository:         repositories.UserBlockStatusRepository,
		btdRepository:         repositories.BlockTypeDictionaryRepository,
		userHistoryRepository: repositories.UserHistoryRepository,
		kitchenService:        kitchenService,
		deleteUserProducer:    producers.DeleteUser,
		transactor:            repositories.Transactor,
		redis:                 redis,
	}
}
