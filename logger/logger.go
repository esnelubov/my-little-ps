package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"my-little-ps/config"
)

type Log struct {
	logger *zap.Logger
}

func New(config config.IConfig) *Log {
	var (
		logger *zap.Logger
		err    error
		env    string
	)

	env = config.GetString("env")

	if env == "prod" {
		logger, err = newProd()
	} else {
		logger, err = newDev()
	}

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return &Log{
		logger: logger,
	}
}

func newProd() (*zap.Logger, error) {
	var cfg = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}

func newDev() (*zap.Logger, error) {
	var cfg = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}

func (l *Log) Sync() error {
	return l.logger.Sync()
}

func (l *Log) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *Log) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Log) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Log) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Log) Debugf(msg string, param ...interface{}) {
	l.logger.Debug(fmt.Sprintf(msg, param...))
}

func (l *Log) Infof(msg string, param ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, param...))
}

func (l *Log) Warnf(msg string, param ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, param...))
}

func (l *Log) Errorf(msg string, param ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, param...))
}
