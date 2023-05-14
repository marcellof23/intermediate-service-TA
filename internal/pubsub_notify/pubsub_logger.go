package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type ZapLoggerAdapter struct {
	*zap.Logger
}

func (l *ZapLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	l.Logger.Error(msg, zap.Error(err), zap.Any("fields", fields))
}

func (l *ZapLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	l.Logger.Info(msg, zap.Any("fields", fields))
}

func (l *ZapLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	l.Logger.Debug(msg, zap.Any("fields", fields))
}

func (l *ZapLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	l.Logger.Debug(msg, zap.Any("fields", fields))
}

func (l *ZapLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &ZapLoggerAdapter{
		Logger: l.Logger.With(zap.Any("fields", fields)),
	}
}
