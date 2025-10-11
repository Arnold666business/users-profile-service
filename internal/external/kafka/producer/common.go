package producer

import (
	"context"
	"time"
	"users-profile-service/internal/config"

	"go.uber.org/zap"
)

func NewProducerConfig() (*ProducerConfig, error) {
	cfg := &ProducerConfig{}
	err := config.Load(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
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
