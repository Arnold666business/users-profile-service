package build

import (
	"fmt"
	"users-profile-service/internal/external/kafka/producer/DeletedUser"
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/external/kafka/producer/UnBlockUser"
	"users-profile-service/pkg/close"

	"go.uber.org/zap"
)

type Producers struct {
	DeleteUser  *DeletedUser.Producer
	NewUser     *NewUser.Producer
	UnBlockUser *UnBlockUser.Producer
}

func Build(logger *zap.SugaredLogger, closer *close.Closer) (*Producers, error) {
	deleteUser, err := DeletedUser.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building deleted user producer: %w", err)
	}
	closer.Add(deleteUser.Close)

	newUser, err := NewUser.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building new user producer: %w", err)
	}
	closer.Add(newUser.Close)

	unBlockUser, err := UnBlockUser.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building unblock user producer: %w", err)
	}
	closer.Add(unBlockUser.Close)

	return &Producers{deleteUser, newUser, unBlockUser}, nil
}
