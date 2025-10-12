package redis

import "users-profile-service/internal/config"

type RedisConfig struct {
	RedisHost           string `env:"REDIS_HOST"`
	RedisPort           string `env:"REDIS_PORT"`
	RedisPassword       string `env:"REDIS_PASSWORD"`
	RedisUserTTL        int    `env:"REDIS_USER_TTL_HOURS"`
	RedisIdempotencyTTL int    `env:"REDIS_USER_TTL_HOURS"`
}

func NewDBConfig() (*RedisConfig, error) {
	cfg := &RedisConfig{}
	err := config.Load(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
