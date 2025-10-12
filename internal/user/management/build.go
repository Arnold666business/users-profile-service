package management

import (
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type UserManager struct {
	logger         *zap.SugaredLogger
	userRepository *repository.User
	uhRepository   *repository.UserHistory
	redis          *redis.RedisProvider
}

func Build(
	logger *zap.SugaredLogger,
	userRepository *repository.User,
	uhRepository *repository.UserHistory,
	redis *redis.RedisProvider,
) (*UserManager, error) {
	return &UserManager{
		logger:         logger,
		userRepository: userRepository,
		uhRepository:   uhRepository,
		redis:          redis,
	}, nil
}
