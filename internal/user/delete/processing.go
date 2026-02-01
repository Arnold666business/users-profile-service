package delete

import (
	"context"
	"strconv"
	"time"
	"users-profile-service/internal/models"
)

type Booking struct {
	BookingId int64
	KitchenId int64
}

type KitchenService interface {
	GetKitchenActiveBookingsByOwnerId(ctx context.Context, id int64) []Booking
	UnPublishKitchenByOwerId(ctx context.Context, id int64) error
}

func (processor *DeleteUserProcessor) Process(ctx context.Context, id int64) (int64, error) {
	l := processor.logger.Named("deleted.users.processing")

	user, err := processor.userRepository.GetById(ctx, id)
	if err != nil {
		l.Errorw("error getting users", "userId", id, "err", err)
		return 0, err
	}
	forHistory := *user
	ownerRole := 1
	if user.Role == ownerRole {
		if len(processor.kitchenService.GetKitchenActiveBookingsByOwnerId(ctx, user.Id)) != 0 {
			l.Warnw("kitchen has active bookings", "users", user)
			return 0, nil
		}

		err = processor.kitchenService.UnPublishKitchenByOwerId(ctx, user.Id)

		if err != nil {
			return 0, err
		}
	}

	err = processor.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		user.IsDeleted = true
		now := time.Now()
		user.DeletedAt = &now
		err = processor.userRepository.Update(ctx, user)
		if err != nil {
			l.Errorw("error updating users", "users", user)
			return err
		}

		_, err = processor.userHistoryRepository.Save(ctx, forHistory, models.DELETE)
		if err != nil {
			l.Errorw("error adding user_history", "userId", id, "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	go func() {
		if errR := processor.redis.DeleteUserCache(ctx, strconv.FormatInt(id, 10)); errR != nil {
			processor.logger.Error("error delete users with id %d to cache after delete user", id)
		}
	}()

	go func() {
		err = processor.deleteUserProducer.Produce(id)
		if err != nil {
			processor.logger.Errorw("error deleting user", "userId", id, "err", err)
		}
	}()

	return id, nil
}
