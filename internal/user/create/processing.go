package create

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
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

// todo:транзакции
func (processor *CreateUserProcessor) Process(req CreateRequest) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	l := processor.logger.Named("create.user.processing")

	val, err := processor.redis.GetIdempotencyStorage(ctx, req.IdempotencyKey)
	if val != "" {
		l.Debugf("get idempotency key: %s, value: %s", req.IdempotencyKey, val)
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, err
		}
		return num, nil
	}
	if err != nil {
		l.Errorf("error with getting idempotency key %s: %s", req.IdempotencyKey, err)
	}

	if validateErr := user.ValidateEmail(ctx, req.Email, processor.userRepository); validateErr != "" {
		l.Debugf("User email is invalid: %s", validateErr)
		return 0, fmt.Errorf(validateErr)
	}

	if validateErr := user.ValidateLogin(ctx, req.Login, processor.userRepository); validateErr != "" {
		l.Debugf("User login is invalid: %s", validateErr)
		return 0, fmt.Errorf(validateErr)
	}

	newUser := &models.User{
		Login:             req.Login,
		Email:             req.Email,
		AccessEmailStatus: false,
		Role:              req.Role,
		IsDeleted:         false,
	}
	userId, err := processor.userRepository.Save(ctx, newUser)
	if err != nil {
		l.Debugf("Error creating user: %s", err)
		return 0, err
	}
	newUser.Id = userId // на всякий

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		l.Errorf("failed to marshal user with id %d: %s", userId, err)
		return 0, err
	}
	errI := processor.redis.SetIdempotencyStorage(ctx, req.IdempotencyKey, string(jsonData))
	if errI != nil {
		l.Errorf("error with save to idempotency storage %s: %s", req.IdempotencyKey, errI)
	}

	//todo: вот эти хуйни все сделать нормально
	err = processor.newUserProducer.Produce(NewUser.NewUserTopicData{Id: userId, Role: req.Role})
	if err != nil {
		l.Debugf("Error creating user: %s", err)
	}

	_, err = processor.uhRepository.Save(ctx, newUser, models.BLOCKED)
	if err != nil {
		l.Errorw("error adding user_history", "userId", userId, "err", err)
		return 0, err
	}

	return userId, nil
}
