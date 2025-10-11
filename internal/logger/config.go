package logger

import (
	"users-profile-service/internal/config"
)

type Logger struct {
	Level      string `env:"LOGGING_LEVEL"`
	FileName   string `env:"LOGGING_FILE_NAME"`
	MaxSize    int    `env:"LOGGING_MAX_SIZE"`
	MaxBackups int    `env:"LOGGING_MAX_BACKUPS"`
	MaxAge     int    `env:"LOGGING_MAX_AGE"`
}

func NewLoggerConfig() (*Logger, error) {
	cfg := &Logger{}
	err := config.Load(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
