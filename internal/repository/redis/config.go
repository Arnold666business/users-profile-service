package redis

import (
	"sync"
	"users-profile-service/internal/config"
)

type RedisConfig struct {
	RedisHost           string `env:"REDIS_HOST"`
	RedisPort           string `env:"REDIS_PORT"`
	RedisPassword       string `env:"REDIS_PASSWORD"`
	RedisUserTTL        int    `env:"REDIS_USER_TTL_HOURS"`
	RedisIdempotencyTTL int    `env:"REDIS_IDEMPOTENCY_TTL_HOURS"`
}

var (
	once sync.Once
	cfg  *RedisConfig
)

func NewDBConfig() (*RedisConfig, error) {
	var err error
	once.Do(func() {
		cfg = &RedisConfig{}
		err = config.Load(cfg)
	})
	return cfg, err
}
