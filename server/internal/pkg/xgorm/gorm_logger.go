package xgorm

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type ZapLogger struct {
	level  logger.LogLevel
	logger *zap.Logger
}

func (l *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l *ZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.logger.Info(msg, zap.Any("details", data))
	}
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.logger.Warn(msg, zap.Any("details", data))
	}
}

func (l *ZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.logger.Error(msg, zap.Any("details", data))
	}
}

func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("duration", time.Since(begin)),
	}

	if err != nil {
		fields = append(fields, zap.String("error", err.Error()))
	}

	l.logger.Debug("GORM SQL", fields...)
}
