package management

import (
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type UserManager struct {
	logger                *zap.SugaredLogger
	userRepository        *repository.User
	userHistoryRepository *repository.UserHistory
	redis                 *redis.RedisProvider
	transactor            *repository.Transactor
}

func Build(
	logger *zap.SugaredLogger,
	repositories *repository.Repositories,
	redis *redis.RedisProvider,
) *UserManager {
	return &UserManager{
		logger:                logger,
		userRepository:        repositories.UserRepository,
		userHistoryRepository: repositories.UserHistoryRepository,
		redis:                 redis,
		transactor:            repositories.Transactor,
	}
}
