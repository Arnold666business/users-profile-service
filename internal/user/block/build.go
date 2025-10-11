package block

import (
	"users-profile-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BlockUserProcessor struct {
	logger         *zap.SugaredLogger
	db             *pgxpool.Pool
	userRepository *repository.User
	ubsRepository  *repository.UserBlockStatus
	btdRepository  *repository.BlockTypeDictionary
	uhRepository   *repository.UserHistory
}

func Build(logger *zap.SugaredLogger,
	db *pgxpool.Pool,
	userRepository *repository.User,
	ubsRepository *repository.UserBlockStatus,
	btdRepository *repository.BlockTypeDictionary,
	uhRepository *repository.UserHistory) (*BlockUserProcessor, error) {
	return &BlockUserProcessor{
		logger:         logger,
		db:             db,
		userRepository: userRepository,
		ubsRepository:  ubsRepository,
		btdRepository:  btdRepository,
		uhRepository:   uhRepository,
	}, nil
}
