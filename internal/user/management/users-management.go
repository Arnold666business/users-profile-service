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

func (um *UserManager) EditEmail(ctx context.Context, userId int64, email string) error {
	u, err := um.userRepository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if res := user.ValidateEmail(ctx, email, um.userRepository); res != "" {
		return fmt.Errorf(res)
	}

	err = um.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		u.Email = email
		u.AccessEmailStatus = false
		err = um.userRepository.Update(ctx, u)
		if err != nil {
			um.logger.Errorf("error updating users: %v", err)
			return err
		}

		_, err = um.userHistoryRepository.Save(ctx, u, models.UPDATE)
		if err != nil {
			um.logger.Errorw("error adding user_history", "userId", userId, "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	res, errm := json.Marshal(u)
	if errm != nil {
		return err
	}
	if errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(res)); errR != nil {
		um.logger.Error("error set users with id %d to cache after change email", userId)
	}
	return nil
}

func (um *UserManager) ConfirmEmail(ctx context.Context, userId int64) error {
	u, err := um.userRepository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if u.Email == "" {
		return errors.New("no users email found")
	}

	err = um.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		u.AccessEmailStatus = true
		err = um.userRepository.UpdateAccessEmailStatus(ctx, userId, "true")
		if err != nil {
			um.logger.Errorf("error updating users: %v", err)
			return err
		}

		_, err = um.userHistoryRepository.Save(ctx, u, models.UPDATE)
		if err != nil {
			um.logger.Errorw("error adding user_history", "userId", userId, "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	res, errm := json.Marshal(u)
	if errm != nil {
		return err
	}
	if errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(res)); errR != nil {
		um.logger.Error("error set users with id %d to cache after change access email status", userId)
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

	err = um.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		u.Login = login
		err = um.userRepository.UpdateLogin(ctx, userId, login)
		if err != nil {
			um.logger.Errorf("error updating users: %v", err)
			return err
		}

		_, err = um.userHistoryRepository.Save(ctx, u, models.UPDATE)
		if err != nil {
			um.logger.Errorw("error adding user_history", "userId", userId, "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	res, errm := json.Marshal(u)
	if errm != nil {
		return err
	}
	if errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(res)); errR != nil {
		um.logger.Error("error set users with id %d to cache after change login", userId)
	}

	return nil
}

func (um *UserManager) GetUser(ctx context.Context, userId int64) (*models.UserProfile, error) {
	val, err := um.redis.GetUsersCache(ctx, strconv.FormatInt(userId, 10))
	if err != nil {
		um.logger.Errorf("error getting users with id %d from cache", userId)
	}

	if val != "" {
		var u models.UserProfile
		err = json.Unmarshal([]byte(val), &u)
		if err != nil {
			um.logger.Errorf("error unmarshalling users with id %d from cache", userId)
		} else {
			return &u, nil
		}
	}

	userAggregate, errA := um.userRepository.GetUserAggregate(ctx, userId)
	if errA != nil {
		return nil, errA
	}

	jsonAggregate, _ := json.Marshal(userAggregate)
	errR := um.redis.SetUserCache(ctx, strconv.FormatInt(userId, 10), string(jsonAggregate))
	if errR != nil {
		um.logger.Error("error set users with id %d to cache after change users", userId)
	}

	return userAggregate, nil
}
