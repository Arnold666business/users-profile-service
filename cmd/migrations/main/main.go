package main

import (
	"fmt"
	"users-profile-service/internal/logger"
	"users-profile-service/internal/repository/postgres"
)

func main() {

}

func doMigrate() error {
	l, err := logger.Build()
	if err != nil {
		return fmt.Errorf("error with init logger in migrations: %v", err)
	}
	l.Named("db.migrations")

	db, err := postgres.Build(l)
	if err != nil {
		return fmt.Errorf("error with init postgres connection in migrations: %v", err)
	}

}
