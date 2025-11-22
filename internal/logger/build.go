package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Build() (*zap.SugaredLogger, error) {
	cfg, err := NewLoggerConfig()
	if err != nil {
		return nil, err
	}
	logLevel := cfg.Level
	var level zapcore.Level
	err = level.Set(logLevel)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,    //будет создан новый файл
		MaxBackups: cfg.MaxBackups, //сколько старых файлов хранить
		MaxAge:     cfg.MaxAge,
		Compress:   true, //gzip compression
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(fileWriter),
		level,
	)

	stdoutCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		level,
	)

	core := zapcore.NewTee(fileCore, stdoutCore)
	logger := zap.New(core, zap.AddCaller()).With(
		zap.String("users", "users-profile-users"))
	return logger.Sugar(), nil
}
