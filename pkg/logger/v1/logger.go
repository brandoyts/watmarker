package logger

import (
	"errors"
	"log"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

type Config struct {
	LogLevel string
}

func NewLogger(config *Config) (*Logger, error) {
	var level zapcore.Level

	err := level.Set(config.LogLevel)
	if err != nil {
		return nil, err
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.Level.SetLevel(level)

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger.Sugar()}, nil
}

func (l *Logger) Sync() {
	err := l.SugaredLogger.Sync()
	if err == nil {
		return
	}

	if errors.Is(err, syscall.EINVAL) {
		return
	}

	log.Printf("Failed to sync logger: %v", err)
}
