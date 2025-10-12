package management

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"users-profile-service/internal/models"
	"users-profile-service/internal/user"
)

//todo: транзакции

func (um *UserManager) EditEmail(ctx context.Context, userId int64, email string) error {
	u, err := um.userRepository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if res := user.ValidateEmail(ctx, email, um.userRepository); res != "" {
		return fmt.Errorf(res)
	}

	u.Email = email
	u.AccessEmailStatus = false
	err = um.userRepository.Update(ctx, u)
	if err != nil {
		return err
	}

	res, errm := json.Marshal(u)
	if errm != nil {
		return err
	}
	if errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(res)); errR != nil {
		um.logger.Error("error set user with id %d to cache after change email", userId)
	}

	_, err = um.uhRepository.Save(ctx, u, models.UPDATE)
	if err != nil {
		um.logger.Errorw("error adding user_history", "userId", userId, "err", err)
		return err
	}
	return nil
}

func (um *UserManager) ConfirmEmail(ctx context.Context, userId int64) error {
	u, err := um.userRepository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if u.Email == "" {
		return errors.New("no user email found")
	}

	u.AccessEmailStatus = true
	err = um.userRepository.UpdateAccessEmailStatus(ctx, userId, "true")
	if err != nil {
		return err
	}

	res, errm := json.Marshal(u)
	if errm != nil {
		return err
	}
	if errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(res)); errR != nil {
		um.logger.Error("error set user with id %d to cache after change access email status", userId)
	}

	_, err = um.uhRepository.Save(ctx, u, models.UPDATE)
	if err != nil {
		um.logger.Errorw("error adding user_history", "userId", userId, "err", err)
		return err
	}
	return nil
}

func (um *UserManager) EditLogin(ctx context.Context, userId int64, login string) error {
	u, err := um.userRepository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if res := user.ValidateLogin(ctx, login, um.userRepository); res != "" {
		return fmt.Errorf(res)
	}

	u.Login = login
	err = um.userRepository.UpdateLogin(ctx, userId, login)
	if err != nil {
		return err
	}

	res, errm := json.Marshal(u)
	if errm != nil {
		return err
	}
	if errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(res)); errR != nil {
		um.logger.Error("error set user with id %d to cache after change login", userId)
	}

	_, err = um.uhRepository.Save(ctx, u, models.UPDATE)
	if err != nil {
		um.logger.Errorw("error adding user_history", "userId", userId, "err", err)
		return err
	}
	return nil
}

func (um *UserManager) GetUser(ctx context.Context, userId int64) (*models.User, error) {
	val, err := um.redis.GetUsersCache(ctx, strconv.FormatInt(userId, 10))
	if err != nil {
		um.logger.Errorf("error getting user with id %d from cache", userId)
	}

	if val != "" {
		var u models.User
		err = json.Unmarshal([]byte(val), &u)
		if err != nil {
			um.logger.Errorf("error unmarshalling user with id %d from cache", userId)
		} else {
			return &u, nil
		}
	}

	//todo: когда приду: 1) вынести повторяющуюся мочу от сюда. 2) дописать GetUser 3) с трнанзакциями разобраться
}
