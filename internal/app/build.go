package app

import (
	"fmt"
	"users-profile-service/internal/api/rest/public/http"
	"users-profile-service/internal/external/kafka/producer/DeletedUser"
	"users-profile-service/internal/external/kafka/producer/NewUser"
	"users-profile-service/internal/logger"
	"users-profile-service/internal/repository/postgres"
)

func Build() error {
	l, err := logger.Build()
	if err != nil {
		return err
	}

	db, err := postgres.Build(l)
	if err != nil {
		l.Fatalf("init postgres failed: %s", err)
	}

	newUserProducer, err := NewUser.Build(l)
	if err != nil {
		return fmt.Errorf("init newUserProducer failed: %s", err)
	}
	deletedUserProducer, err := DeletedUser.Build(l)
	if err != nil {
		return fmt.Errorf("init deletedUserProducer failed: %s", err)
	}

	server := http.Build()
	go server.Start()

}
