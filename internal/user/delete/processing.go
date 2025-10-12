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

// todo: добавить транзакции ну или ненадо
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

	user.IsDeleted = true
	user.DeletedAt = time.Now()

	//todo: вот эти хуйни все сделать нормально
	//processor.auditProducer.Produce()

	//todo: вот эти хуйни все сделать нормально
	err = processor.deleteUserProducer.Produce(ctx, id)
	if err != nil {
		return 0, err
	}

	_, err = processor.uhRepository.Save(ctx, user, models.DELETE)
	if err != nil {
		l.Errorw("error adding user_history", "userId", id, "err", err)
		return 0, err
	}
	return id, nil
}
