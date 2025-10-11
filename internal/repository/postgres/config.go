package postgres

import (
	"fmt"
	"time"
	"users-profile-service/internal/config"
)

type DB struct {
	Username          string        `env:"DB_USER"`
	Password          string        `env:"DB_PASSWORD"`
	Host              string        `env:"DB_HOST"`
	Port              string        `env:"DB_PORT"`
	Name              string        `env:"DB_NAME"`
	ConnMaxLifetime   time.Duration `env:"DB_CONN_MAX_LIFETIME"`
	MaxOpenConns      int32         `env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns      int32         `env:"DB_MAX_IDLE_CONNS"`
	RetryInterval     int           `env:"DB_RETRY_INTERVAL"`
	MaxRetries        int           `env:"DB_MAX_RETRIES"`
	BackoffMultiplier float64       `env:"DB_BACKOFF_MULTIPLIER"`
}

func NewDBConfig() (*DB, error) {
	cfg := &DB{}
	err := config.Load(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (db *DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.Username, db.Password, db.Host, db.Port, db.Name)
}
