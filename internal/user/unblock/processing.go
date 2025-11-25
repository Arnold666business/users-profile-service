package unblock

import (
	"context"
	"strconv"
	"time"
	"users-profile-service/internal/external/kafka/producer/UnBlockUser"
	"users-profile-service/internal/models"

	"go.uber.org/zap"
)

func (processor *UnblockUserProcessor) Process(ctx context.Context, unblockUserStatus *models.UserBlockStatus) error {
	if unblockUserStatus.IsActive {
		unblockUserStatus.IsActive = false
		err := processor.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			err := processor.userBlockStatusRepository.Update(ctx, unblockUserStatus)
			if err != nil {
				return err
			}

			user, err := processor.userRepository.GetById(ctx, unblockUserStatus.UserId)
			if err != nil {
				return err
			}
			_, err = processor.userHistoryRepository.Save(ctx, user, models.UNBLOCKED)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		go func() {
			err := processor.redis.DeleteUserCache(ctx, strconv.FormatInt(unblockUserStatus.UserId, 10))
			if err != nil {
				processor.logger.Errorw("unblock job redis delete failed", zap.Error(err))
			}
		}()
	}

	err := processor.unblockUserProducer.Produce(UnBlockUser.UnblockUserTopicData{UserId: unblockUserStatus.UserId, Timestamp: time.Now()})
	if err != nil {
		return err
	}

	unblockUserStatus.UnBlockEventSent = true
	err = processor.userBlockStatusRepository.Update(ctx, unblockUserStatus)
	if err != nil {
		return err
	}
	return nil
}

func (processor *UnblockUserProcessor) FindForUnblock(ctx context.Context, findLimit int) ([]*models.UserBlockStatus, error) {
	return processor.userBlockStatusRepository.FindForUnblockProcessing(ctx, findLimit)
}
