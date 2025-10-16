package block

import (
	"context"
	"users-profile-service/internal/models"
)

type BlockRequest struct {
	UserId      int64
	BlockTypeId int
	ForeverFlag bool
}

func (processor *BlockUserProcessor) Process(ctx context.Context, data BlockRequest) error {
	l := processor.logger.Named("block.user.processing")

	userId := data.UserId
	user, err := processor.userRepository.GetById(ctx, userId)
	if err != nil {
		l.Errorw("error getting user", "userId", userId, "err", err)
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

		_, err = processor.uhRepository.Save(ctx, user, models.BLOCKED)
		if err != nil {
			l.Errorw("error adding user_history", "userId", userId, "err", err)
			return err
		}
		return nil
	})
}
