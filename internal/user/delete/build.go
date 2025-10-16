package delete

import (
	"users-profile-service/internal/external/kafka/producer/Audit"
	"users-profile-service/internal/external/kafka/producer/DeletedUser"
	"users-profile-service/internal/repository"

	"go.uber.org/zap"
)

type DeleteUserProcessor struct {
	logger             *zap.SugaredLogger
	userRepository     *repository.User
	ubsRepository      *repository.UserBlockStatus
	btdRepository      *repository.BlockTypeDictionary
	uhRepository       *repository.UserHistory
	kitchenService     KitchenService
	deleteUserProducer *DeletedUser.Producer
	auditProducer      *Audit.Producer
	transactor         *repository.Transactor
}

func Build(logger *zap.SugaredLogger,
	userRepository *repository.User,
	ubsRepository *repository.UserBlockStatus,
	btdRepository *repository.BlockTypeDictionary,
	uhRepository *repository.UserHistory,
	kitchenService KitchenService,
	deleteUserProducer *DeletedUser.Producer,
	auditProducer *Audit.Producer,
	transactor *repository.Transactor) (*DeleteUserProcessor, error) {
	return &DeleteUserProcessor{
		logger:             logger,
		userRepository:     userRepository,
		ubsRepository:      ubsRepository,
		btdRepository:      btdRepository,
		uhRepository:       uhRepository,
		kitchenService:     kitchenService,
		deleteUserProducer: deleteUserProducer,
		auditProducer:      auditProducer,
		transactor:         transactor,
	}, nil
}
