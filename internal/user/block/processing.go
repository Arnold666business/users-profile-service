package block

import (
	"context"
	"errors"
	"users-profile-service/internal/models"
	"users-profile-service/internal/repository"
)

type BlockRequest struct {
	UserId      int64
	BlockTypeId int
	ForeverFlag bool
}

// todo: добавить транзакции ну или ненадо
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
	}
	ubs.SetUnblockDate(btd.Hour)
	currentUbs, err := processor.ubsRepository.GetByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.NotFoundUserBlockStatusError) {
			_, err := processor.ubsRepository.Save(ctx, ubs)
			if err != nil {
				l.Errorw("error adding block", "userId", userId, "err", err)
				return err
			}
		} else {
			l.Warnw("error getting block", "userId", userId, "err", err)
			return err
		}
	} else {
		ubs.Id = currentUbs.Id
		err = processor.ubsRepository.Update(ctx, ubs)
		if err != nil {
			l.Errorw("error updating block", "userId", userId, "err", err)
			return err
		}
	}

	_, err = processor.uhRepository.Save(ctx, user, models.BLOCKED)
	if err != nil {
		l.Errorw("error adding user_history", "userId", userId, "err", err)
		return err
	}
	return nil
}
