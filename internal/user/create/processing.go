package create

import (
	"context"
	"strconv"
	common_error "users-profile-service/internal/common-error"
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/models"
	"users-profile-service/internal/user"
)

type CreateRequest struct {
	Login          string `json:"login"`
	Email          string `json:"email"`
	Role           int    `json:"role"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (processor *CreateUserProcessor) Process(ctx context.Context, req CreateRequest) (int64, error) {
	l := processor.logger.Named("create.users.processing")

	val, err := processor.redis.GetIdempotencyStorage(ctx, req.IdempotencyKey)
	if val != "" {
		l.Debugf("get idempotency key: %s, value: %s", req.IdempotencyKey, val)
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, common_error.NewError(err.Error(), common_error.TypeInternal)
		}
		return num, nil
	}
	if err != nil {
		l.Errorf("error with getting idempotency key %s: %s", req.IdempotencyKey, err)
	}

	if validateErr := user.ValidateEmail(ctx, req.Email, processor.userRepository); validateErr != "" {
		l.Debugf("User email is invalid: %s", validateErr)
		return 0, common_error.NewError(validateErr, common_error.TypeInternal)
	}

	if validateErr := user.ValidateLogin(ctx, req.Login, processor.userRepository); validateErr != "" {
		l.Debugf("User login is invalid: %s", validateErr)
		return 0, common_error.NewError(validateErr, common_error.TypeInternal)
	}
	newUser := &models.User{
		Login:             req.Login,
		Email:             req.Email,
		AccessEmailStatus: false,
		Role:              req.Role,
		IsDeleted:         false,
	}
	errTx := processor.transactor.WithinTransaction(ctx, func(ctx context.Context) error {

		userId, err := processor.userRepository.Save(ctx, newUser)
		if err != nil {
			l.Debugf("Error creating users: %s", err)
			return err
		}

		_, err = processor.userHistoryRepository.Save(ctx, newUser, models.CREATE)
		if err != nil {
			l.Errorw("error adding user_history", "userId", userId, "err", err)
			return err
		}
		return nil
	})
	if errTx != nil {
		return 0, common_error.NewError(errTx.Error(), common_error.TypeInternal)
	}

	errI := processor.redis.SetIdempotencyStorage(ctx, req.IdempotencyKey, strconv.FormatInt(newUser.Id, 10))
	if errI != nil {
		l.Errorf("error with save to idempotency storage %s: %s", req.IdempotencyKey, errI)
	}

	//todo: залупа
	go func() {
		errProducer := processor.newUserProducer.Produce(NewUser.NewUserTopicData{Id: newUser.Id, Role: req.Role})
		if errProducer != nil {
			l.Debugf("Error send user new user event: %s", errProducer)
		}
	}()

	return newUser.Id, nil
}
