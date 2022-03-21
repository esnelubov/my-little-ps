package logger

import (
	"go.uber.org/zap"
	"log"
)

type Log struct {
	logger *zap.Logger
}

func New() *Log {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return &Log{
		logger: logger,
	}
}

func (l *Log) Sync() error {
	return l.logger.Sync()
}
