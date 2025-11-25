package logger

import (
	"sync"
	"users-profile-service/internal/config"
)

type Logger struct {
	Level      string `env:"LOGGING_LEVEL"`
	FileName   string `env:"LOGGING_FILE_NAME"`
	MaxSize    int    `env:"LOGGING_MAX_SIZE"`
	MaxBackups int    `env:"LOGGING_MAX_BACKUPS"`
	MaxAge     int    `env:"LOGGING_MAX_AGE"`
}

var (
	once        sync.Once
	cfgInstance *Logger
)

func NewLoggerConfig() error {
	var err error
	once.Do(func() {
		cfgInstance = &Logger{}
		err = config.Load(cfgInstance)
	})
	return err
}
