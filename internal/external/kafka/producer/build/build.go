package build

import (
	"fmt"
	"users-profile-service/internal/external/kafka/producer/DeletedUser"
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/external/kafka/producer/UnBlockUser"

	"go.uber.org/zap"
)

type Producers struct {
	DeleteUser  *DeletedUser.Producer
	NewUser     *NewUser.Producer
	UnBlockUser *UnBlockUser.Producer
}

func Build(logger *zap.SugaredLogger) (*Producers, error) {
	deleteUser, err := DeletedUser.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building deleted user producer: %w", err)
	}

	newUser, err := NewUser.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building new user producer: %w", err)
	}

	unBlockUser, err := UnBlockUser.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building unblock user producer: %w", err)
	}
	return &Producers{deleteUser, newUser, unBlockUser}, nil
}
