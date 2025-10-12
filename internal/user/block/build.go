package block

import (
	"users-profile-service/internal/repository"

	"go.uber.org/zap"
)

type BlockUserProcessor struct {
	logger         *zap.SugaredLogger
	userRepository *repository.User
	ubsRepository  *repository.UserBlockStatus
	btdRepository  *repository.BlockTypeDictionary
	uhRepository   *repository.UserHistory
}

func Build(logger *zap.SugaredLogger,
	userRepository *repository.User,
	ubsRepository *repository.UserBlockStatus,
	btdRepository *repository.BlockTypeDictionary,
	uhRepository *repository.UserHistory) (*BlockUserProcessor, error) {
	return &BlockUserProcessor{
		logger:         logger,
		userRepository: userRepository,
		ubsRepository:  ubsRepository,
		btdRepository:  btdRepository,
		uhRepository:   uhRepository,
	}, nil
}
