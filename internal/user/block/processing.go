package block

import (
	"context"
	"strconv"
	"users-profile-service/internal/models"
)

type BlockRequest struct {
	UserId      int64
	BlockTypeId int
	ForeverFlag bool
}

func (processor *BlockUserProcessor) Process(ctx context.Context, data BlockRequest) error {
	l := processor.logger.Named("block.users.processing")

	userId := data.UserId
	user, err := processor.userRepository.GetById(ctx, userId)
	if err != nil {
		l.Errorw("error getting users", "userId", userId, "err", err)
		return err
	}

	btd, err := processor.btdRepository.GetByBlockType(ctx, data.BlockTypeId)
	if err != nil {
		l.Errorw("error getting block_type_dictionary", "blockTypeId", data.BlockTypeId)
		return err
	}
	ubs := &models.UserBlockStatus{
		UserId:      userId,
		BlockTypeId: btd.BlockType,
		ForeverFlag: data.ForeverFlag,
		IsActive:    true,
	}
	ubs.SetUnblockDate(btd.Hour)

	return processor.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		err = processor.ubsRepository.UpsertByUserId(ctx, ubs)
		if err != nil {
			l.Errorf("error upserting block_type_dictionary: %v", err)
			return err
		}

		_, err = processor.userHistoryRepository.Save(ctx, user, models.BLOCKED)
		if err != nil {
			l.Errorw("error adding user_history", "userId", userId, "err", err)
			return err
		}

		go func() {
			if errR := processor.redis.DeleteUserCache(ctx, strconv.FormatInt(userId, 10)); errR != nil {
				processor.logger.Error("error delete users with id %d to cache after block user", userId)
			}
		}()

		return nil
	})
}
