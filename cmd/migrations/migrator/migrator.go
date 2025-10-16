package migrator

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Migrator struct {
	db     *pgxpool.Pool
	logger *zap.SugaredLogger
}
