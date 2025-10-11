package delete

import (
	"context"
	"time"
	"users-profile-service/internal/models"
)

type Booking struct {
	??
}

type KitchenService interface {
	GetKitchenActiveBookingsByOwnerId(ctx context.Context, id int64) []Booking
	UnPublishKitchenByOwerId(ctx context.Context, id int64)
}

func (processor *DeleteUserProcessor) process(id int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l := processor.logger.Named("deleted.processing")

	user, err := processor.userRepository.GetById(ctx, id)
	if err != nil {
		l.Errorw("error getting user", "userId", id, "err", err)
		return 0, err
	}

	if len(processor.kitchenService.GetKitchenActiveBookingsByOwnerId(ctx, user.Id)) != 0 {
		l.Warnw("kitchen has active bookings", "user", user)
		return 0, nil
	}

	go processor.kitchenService.UnPublishKitchenByOwerId(ctx, user.Id)

	user.IsDeleted = true
	user.DeletedAt = time.Now()

	processor.auditProducer.Produce() ??

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
