package producer

import (
	"fmt"
	"users-profile-service/internal/external/kafka/producer/Audit"
	"users-profile-service/internal/external/kafka/producer/DeletedUser"
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/external/kafka/producer/UnBlockUser"

	"go.uber.org/zap"
)

type Producers struct {
	Audit       *Audit.Producer
	DeleteUser  *DeletedUser.Producer
	NewUser     *NewUser.Producer
	UnBlockUser *UnBlockUser.Producer
}

func Build(logger *zap.SugaredLogger) (*Producers, error) {
	audit, err := Audit.Build(logger)
	if err != nil {
		return nil, fmt.Errorf("error building audit producer: %w", err)
	}

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
	return &Producers{audit, deleteUser, newUser, unBlockUser}, nil
}
