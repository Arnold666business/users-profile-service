package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Build(logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	cfg, err := NewDBConfig()
	if err != nil {
		return nil, err
	}
	logger.Named("db.connection").With(
		zap.String("db.url", cfg.DSN()))

	conn, err := pgxpool.New(context.Background(), cfg.DSN())
	if err != nil {
		return nil, err
	}
	conn.Config().MinIdleConns = cfg.MaxIdleConns
	conn.Config().MaxConns = cfg.MaxOpenConns
	conn.Config().MaxConnLifetime = cfg.ConnMaxLifetime

	return conn, nil
}
