package delete

import (
	"context"
	"time"
	"users-profile-service/internal/models"
)

type Booking struct {
	//???asdfghj
}

type KitchenService interface {
	GetKitchenActiveBookingsByOwnerId(ctx context.Context, id int64) []Booking
	UnPublishKitchenByOwerId(ctx context.Context, id int64)
}

func (processor *DeleteUserProcessor) process(ctx context.Context, id int64) (int64, error) {
	l := processor.logger.Named("deleted.user.processing")

	user, err := processor.userRepository.GetById(ctx, id)
	if err != nil {
		l.Errorw("error getting user", "userId", id, "err", err)
		return 0, err
	}

	//todo: вот эти хуйни все сделать нормально
	if len(processor.kitchenService.GetKitchenActiveBookingsByOwnerId(ctx, user.Id)) != 0 {
		l.Warnw("kitchen has active bookings", "user", user)
		return 0, nil
	}

	//todo: вот эти хуйни все сделать нормально
	go processor.kitchenService.UnPublishKitchenByOwerId(ctx, user.Id)

	err = processor.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		user.IsDeleted = true
		user.DeletedAt = time.Now()
		err = processor.userRepository.Update(ctx, user)
		if err != nil {
			l.Errorw("error updating user", "user", user)
			return err
		}

		_, err = processor.uhRepository.Save(ctx, user, models.DELETE)
		if err != nil {
			l.Errorw("error adding user_history", "userId", id, "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	//todo: вот эти хуйни все сделать нормально
	//processor.auditProducer.Produce()

	//todo: вот эти хуйни все сделать нормально
	err = processor.deleteUserProducer.Produce(ctx, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
