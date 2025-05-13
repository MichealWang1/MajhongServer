package mysql

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"kxmj.common/log"
	"time"
)

type Logger struct {
	logger *zap.Logger
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	return NewLogger()
}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
	l.logger.Sugar().Infof(s, i)
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.logger.Sugar().Warnf(s, i)
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
	l.logger.Sugar().Errorf(s, i)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	l.logger.Debug("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))

	elapsedMillisecond := elapsed / time.Millisecond
	// 慢sql增加一条warn记录
	if elapsedMillisecond >= 100*2 {
		l.logger.Warn("warn", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}

	// 错误sql增加一条错误日志
	if err != nil {
		l.logger.Error("error", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.Default(),
	}
}
