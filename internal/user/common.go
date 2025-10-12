package user

import (
	"context"
	"errors"
	"strings"
	"users-profile-service/internal/repository"
)

var (
	minLoginSymbols   = 5
	maxLoginSymbols   = 30
	wrongLoginSymbols = "%$#*"
)

// "" если валиден
func ValidateEmail(ctx context.Context, email string, userRepository *repository.User) string {
	if strings.TrimSpace(email) == "" {
		return "email cant be empty"
	}

	if _, err := userRepository.GetByEmail(ctx, email); err != nil {
		if !errors.Is(err, repository.NotFoundUserError) {
			return err.Error()
		}
	} else {
		return "user with email " + email + " already exists"
	}
	return ""
}

// "" если валиден
func ValidateLogin(ctx context.Context, login string, userRepository *repository.User) string {
	if strings.TrimSpace(login) == "" {
		return "login cant be empty"
	}
	if len(login) < minLoginSymbols {
		return "login must be at least 5 characters"
	}
	if len(login) > maxLoginSymbols {
		return "login must be at most 30 characters"
	}
	for _, r := range wrongLoginSymbols {
		if strings.ContainsRune(login, r) {
			return "login contains invalid symbol: " + string(r)
		}
	}

	if _, err := userRepository.GetByEmail(ctx, login); err != nil {
		if !errors.Is(err, repository.NotFoundUserError) {
			return err.Error()
		}
	} else {
		return "user with login " + login + " already exists"
	}
	return ""
}
