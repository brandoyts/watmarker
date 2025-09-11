package logger

import (
	"log"
	"os"

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
	if err := l.SugaredLogger.Sync(); err != nil && err != os.ErrInvalid {
		log.Printf("Failed to sync logger: %v", err)
	}
}
