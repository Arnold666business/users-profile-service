package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Build(logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	err := NewDBConfig()
	if err != nil {
		return nil, err
	}
	logger.Named("db.connection").With(
		zap.String("db.url", cfgInstance.DSN()))

	conn, err := pgxpool.New(context.Background(), cfgInstance.DSN())
	if err != nil {
		return nil, err
	}
	conn.Config().MinIdleConns = cfgInstance.MaxIdleConns
	conn.Config().MaxConns = cfgInstance.MaxOpenConns
	conn.Config().MaxConnLifetime = cfgInstance.ConnMaxLifetime

	go pingWithAttempts(conn, &retryPingPolicy{
		maxAttempts: 5,
		delaySec:    10,
		delayKef:    1.5,
	}, logger)

	return conn, nil
}

type retryPingPolicy struct {
	maxAttempts int
	delaySec    float64
	delayKef    float64
}

func (p *retryPingPolicy) changeDelay() {
	p.delaySec = p.delaySec * p.delayKef
}

func pingWithAttempts(conn *pgxpool.Pool, policy *retryPingPolicy, logger *zap.SugaredLogger) {
	attempts := 0
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := conn.Ping(ctx)
		cancel()

		attempts++
		if err == nil {
			logger.Debugw("ping postgres successful", "attempts", attempts)
			return
		}

		logger.Errorw("ping db failed", "error", err, "attempts", attempts)
		if attempts > policy.maxAttempts {
			panic("All attempts are exceeded. Unable to connect to PostgreSQL")
		}
		policy.changeDelay()
		time.Sleep(time.Duration(policy.delaySec) * time.Second)
	}
}
