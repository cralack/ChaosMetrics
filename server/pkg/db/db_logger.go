package db

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type ZapLogger struct {
	logger *zap.Logger
}

func (l *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *ZapLogger) Info(context.Context, string, ...interface{}) {
	l.logger.Info("GORM log")
}

func (l *ZapLogger) Warn(context.Context, string, ...interface{}) {
	l.logger.Warn("GORM log")
}

func (l *ZapLogger) Error(context.Context, string, ...interface{}) {
	l.logger.Error("GORM log")
}

func (l *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	l.logger.Debug("GORM SQL",
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("duration", time.Since(begin)),
	)
}
