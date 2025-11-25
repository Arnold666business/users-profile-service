package producer

import (
	"context"
	"sync"
	"time"
	"users-profile-service/internal/config"

	"go.uber.org/zap"
)

type ProducerConfig struct {
	Brokers          []string `env:"KAFKA_BROKERS" env-separator:","`
	NewUserTopic     string   `env:"NEW_USER_TOPIC"`
	DeletedUserTopic string   `env:"DELETED_USER_TOPIC"`
	RetryAttempts    int      `env:"KAFKA_RETRY_ATTEMPTS"`
	RetryInterval    int      `env:"KAFKA_RETRY_INTERVAL_SECONDS"`
	Timeout          int      `env:"KAFKA_SEND_TIMEOUT_SECONDS"`
}

var (
	once        sync.Once
	cfgInstance *ProducerConfig
)

func NewProducerConfig() error {
	var err error
	once.Do(func() {
		cfgInstance = &ProducerConfig{}
		err = config.Load(cfgInstance)
	})
	return err
}

func ProduceWithRetry(l *zap.SugaredLogger, cfg *ProducerConfig, produce func(ctx context.Context) error) error {
	var err error
	for i := 0; i < cfg.RetryAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.RetryInterval)*time.Second)
		err = produce(ctx)
		cancel()
		if err == nil {
			return nil
		}
		l.Warnf("%d sand message attempt failed with error: %s", i, err)
		if i+1 == cfg.RetryAttempts {
			return err
		}
		time.Sleep(time.Duration(cfg.RetryInterval) * time.Second)
	}
	return err
}
