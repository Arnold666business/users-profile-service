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

// todo: нормально настроить
func Build(logger *zap.SugaredLogger) (*RedisProvider, error) {
	cfg, err := NewDBConfig()
	if err != nil {
		return nil, err
	}

	address := cfg.RedisHost + ":" + cfg.RedisPort
	logger.Named("redis.client.connection").With(
		zap.String("redis.ulr", address))

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{address},
		Password:    cfg.RedisPassword,
		//DisableRetry: true,
		//DisableCache: true,
		//ConnWriteTimeout: 10 * time.Second,
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

func (rb *RedisProvider) GetUsersCache(ctx context.Context, key string) (string, error) {
	value, err := rb.get(ctx, userCachePrefix, key)
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", nil
		} else {
			return "", err
		}
	}
	return value, nil
}

func (rb *RedisProvider) GetIdempotencyStorage(ctx context.Context, key string) (string, error) {
	value, err := rb.get(ctx, idempotencyCachePrefix, key)
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", nil
		} else {
			return "", err
		}
	}
	return value, nil
}

func (rb *RedisProvider) SetUserCache(ctx context.Context, key string, value string) error {
	expiration := time.Duration(rb.cfg.RedisUserTTL) * time.Hour
	return rb.set(ctx, userCachePrefix, key, value, expiration)
}

func (rb *RedisProvider) SetIdempotencyStorage(ctx context.Context, key string, value string) error {
	expiration := time.Duration(rb.cfg.RedisIdempotencyTTL) * time.Hour
	return rb.set(ctx, idempotencyCachePrefix, key, value, expiration)
}

func (rd *RedisProvider) set(ctx context.Context, prefix string, key string, value string, expiration time.Duration) error {
	cmd := rd.client.B().Set().Key(prefix + "::" + key).Value(value).Ex(expiration).Build()
	return rd.client.Do(ctx, cmd).Error()
}

func (rd *RedisProvider) get(ctx context.Context, prefix string, key string) (string, error) {
	cmd := rd.client.B().Get().Key(prefix + "::" + key).Build()
	return rd.client.Do(ctx, cmd).ToString()
}
