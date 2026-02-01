package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/rueidis"
	"go.uber.org/zap"
)

type RedisProvider struct {
	client rueidis.Client
	cfg    *RedisConfig
	logger *zap.SugaredLogger
}

var (
	userCachePrefix        = "get_users_cache"
	idempotencyCachePrefix = "idempotency_storage"
)

func Build(logger *zap.SugaredLogger) (*RedisProvider, error) {
	cfg, err := NewDBConfig()
	if err != nil {
		return nil, err
	}
	address := cfg.RedisHost + ":" + cfg.RedisPort
	logger.Named("redis.client.connection").With(
		zap.String("redis.ulr", address))

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:      []string{address},
		Password:         cfg.RedisPassword,
		DisableRetry:     true,
		DisableCache:     true,
		ConnWriteTimeout: 10 * time.Second,
	})
	if err != nil {
		logger.Errorw("Error creating redis client", "error", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		logger.Errorw("Redis connection test failed", "error", err)
		return nil, fmt.Errorf("redis connection test failed: %w", err)
	}

	return &RedisProvider{client, cfg, logger}, nil
}

func (rd *RedisProvider) Close() {
	if rd.client != nil {
		rd.client.Close()
		rd.logger.Debugf("Redis connection closed")
	}
}

func (rd *RedisProvider) GetUsersCache(ctx context.Context, key string) (string, error) {
	value, err := rd.get(ctx, userCachePrefix, key)
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", nil
		} else {
			return "", err
		}
	}
	return value, nil
}

func (rd *RedisProvider) GetIdempotencyStorage(ctx context.Context, key string) (string, error) {
	value, err := rd.get(ctx, idempotencyCachePrefix, key)
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", nil
		} else {
			return "", err
		}
	}
	return value, nil
}

func (rd *RedisProvider) SetUserCache(ctx context.Context, key string, value string) error {
	expiration := time.Duration(rd.cfg.RedisUserTTL) * time.Hour
	return rd.set(ctx, userCachePrefix, key, value, expiration)
}

func (rd *RedisProvider) DeleteUserCache(ctx context.Context, key string) error {
	return rd.delete(ctx, userCachePrefix, key)
}

func (rd *RedisProvider) SetIdempotencyStorage(ctx context.Context, key string, value string) error {
	expiration := time.Duration(rd.cfg.RedisIdempotencyTTL) * time.Hour
	return rd.set(ctx, idempotencyCachePrefix, key, value, expiration)
}

func (rd *RedisProvider) set(ctx context.Context, prefix string, key string, value string, expiration time.Duration) error {
	cmd := rd.client.B().Set().Key(prefix + "::" + key).Value(value).Ex(expiration).Build()
	return rd.client.Do(ctx, cmd).Error()
}

func (rd *RedisProvider) get(ctx context.Context, prefix string, key string) (string, error) {
	cmd := rd.client.B().Get().Key(prefix + "::" + key).Build()
	return rd.client.Do(ctx, cmd).ToString()
}

func (rd *RedisProvider) delete(ctx context.Context, prefix string, key string) error {
	cmd := rd.client.B().Del().Key(prefix + "::" + key).Build()
	return rd.client.Do(ctx, cmd).Error()
}
