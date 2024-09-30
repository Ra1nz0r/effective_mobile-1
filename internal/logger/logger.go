package logger

import (
	"time"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapService interface {
	Debug(fields ...interface{})
	Info(fields ...interface{})
	Error(fields ...interface{})
	Fatal(fields ...interface{})
}

type ZapStorage struct {
	*zap.Logger
}

var Zap ZapService = &ZapStorage{zap.NewNop()}

// Initialize инициализирует логгер.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("parse atomic level error: %w", err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = lvl
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	config.DisableStacktrace = true

	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("logger build error: %w", err)
	}

	Zap = &ZapStorage{logger}

	return nil
}

// Debug логирует сообщения уровня DEBUG.
func (z *ZapStorage) Debug(fields ...interface{}) {
	z.Logger.Sugar().Debugln(fields...)
}

// Info логирует сообщения уровня INFO.
func (z *ZapStorage) Info(fields ...interface{}) {
	z.Logger.Sugar().Infoln(fields...)
}

// Error логирует сообщения уровня ERROR.
func (z *ZapStorage) Error(fields ...interface{}) {
	z.Logger.Sugar().Errorln(fields...)
}

// Fatal логирует сообщения уровня FATAL.
func (z *ZapStorage) Fatal(fields ...interface{}) {
	z.Logger.Sugar().Fatalln(fields...)
}
