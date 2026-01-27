package service

import (
	"sync"
	"users-profile-service/internal/config"
)

type Config struct {
	AccessToken string
	Target      string
}

var (
	once        sync.Once
	cfgInstance *Config
)

func NewRpcServiceConfig() error {
	var err error
	once.Do(func() {
		cfgInstance = &Config{}
		err = config.Load(cfgInstance)
	})
	return err
}
