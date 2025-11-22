package unblock

import (
	"context"
	"os"
	"strconv"
	"time"
	"users-profile-service/internal/external/kafka/producer"
	"users-profile-service/internal/external/kafka/producer/UnBlockUser"
	"users-profile-service/internal/models"
	"users-profile-service/internal/repository"
	"users-profile-service/internal/repository/redis"

	"go.uber.org/zap"
)

type UnblockUserJob struct {
	logger                    *zap.SugaredLogger
	userBlockStatusRepository *repository.UserBlockStatus
	userHistoryRepository     *repository.UserHistory
	userRepository            *repository.User
	transactor                *repository.Transactor
	redis                     *redis.RedisProvider
	unblockUserProducer       *UnBlockUser.Producer
}

func Build(logger *zap.SugaredLogger, repositories *repository.Repositories, producers *producer.Producers) UnblockUserJob {
	return UnblockUserJob{
		logger:                    logger,
		userBlockStatusRepository: repositories.UserBlockStatusRepository,
		userHistoryRepository:     repositories.UserHistoryRepository,
		transactor:                repositories.Transactor,
		userRepository:            repositories.UserRepository,
		unblockUserProducer:       producers.UnBlockUser,
	}
}

func (job *UnblockUserJob) Run(ctx context.Context) error {
	interval, err := strconv.Atoi(os.Getenv("UNBLOCK_JOB_INTERVAL_SECOND"))
	if err != nil {
		job.logger.Fatal("unblock job interval parse failed", zap.Error(err))
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	limit := 100
	workerCount := 5
	for {
		select {
		case <-ticker.C:
			unblockStatusEntities, err := job.userBlockStatusRepository.FindForUnblockProcessing(ctx, limit)
			if err != nil {
				job.logger.Errorw("unblock job fetch failed", zap.Error(err))
				continue
			}

			for _, unblockStatus := range unblockStatusEntities {
				errProcessing := job.unblockProcessing(ctx, unblockStatus)
				if errProcessing != nil {
					job.logger.Errorw("unblock job processing failed", zap.Error(err))
				}
			}
		case <-ctx.Done():
			return nil
		}

	}

}

func (job *UnblockUserJob) unblockProcessing(ctx context.Context, unblockUserStatus *models.UserBlockStatus) error {
	if unblockUserStatus.IsActive {
		unblockUserStatus.IsActive = false
		err := job.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			err := job.userBlockStatusRepository.Update(ctx, unblockUserStatus)
			if err != nil {
				return err
			}

			user, err := job.userRepository.GetById(ctx, unblockUserStatus.UserId)
			if err != nil {
				return err
			}
			_, err = job.userHistoryRepository.Save(ctx, user, models.UNBLOCKED)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		go func() {
			err := job.redis.DeleteUserCache(ctx, strconv.FormatInt(unblockUserStatus.UserId, 10))
			if err != nil {
				job.logger.Errorw("unblock job redis delete failed", zap.Error(err))
			}
		}()
	}

	err := job.unblockUserProducer.Produce(UnBlockUser.UnblockUserTopicData{UserId: unblockUserStatus.UserId, Timestamp: time.Now()})
	if err != nil {
		return err
	}

	unblockUserStatus.UnBlockEventSent = true
	err = job.userBlockStatusRepository.Update(ctx, unblockUserStatus)
	if err != nil {
		return err
	}
	return nil
}
