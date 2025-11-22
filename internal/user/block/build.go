package block

import (
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type BlockUserProcessor struct {
	logger                *zap.SugaredLogger
	userRepository        *repository.User
	ubsRepository         *repository.UserBlockStatus
	btdRepository         *repository.BlockTypeDictionary
	userHistoryRepository *repository.UserHistory
	transactor            *repository.Transactor
	redis                 *redis.RedisProvider
}

func Build(logger *zap.SugaredLogger,
	repositories *repository.Repositories,
	redis *redis.RedisProvider) *BlockUserProcessor {
	return &BlockUserProcessor{
		logger:                logger,
		userRepository:        repositories.UserRepository,
		ubsRepository:         repositories.UserBlockStatusRepository,
		btdRepository:         repositories.BlockTypeDictionaryRepository,
		userHistoryRepository: repositories.UserHistoryRepository,
		transactor:            repositories.Transactor,
		redis:                 redis,
	}
}
