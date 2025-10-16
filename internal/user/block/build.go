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
	transactor     *repository.Transactor
}

func Build(logger *zap.SugaredLogger,
	userRepository *repository.User,
	ubsRepository *repository.UserBlockStatus,
	btdRepository *repository.BlockTypeDictionary,
	uhRepository *repository.UserHistory,
	transactor *repository.Transactor) (*BlockUserProcessor, error) {
	return &BlockUserProcessor{
		logger:         logger,
		userRepository: userRepository,
		ubsRepository:  ubsRepository,
		btdRepository:  btdRepository,
		uhRepository:   uhRepository,
		transactor:     transactor,
	}, nil
}
