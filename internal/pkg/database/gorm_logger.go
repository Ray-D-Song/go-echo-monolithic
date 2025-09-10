package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ray-d-song/go-echo-monolithic/internal/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// GormLogger wraps our custom logger for GORM
type GormLogger struct {
	Logger                    *logger.Logger
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

// NewGormLogger creates a new GORM logger using our custom logger
func NewGormLogger(l *logger.Logger) *GormLogger {
	return &GormLogger{
		Logger:                    l,
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             200 * time.Millisecond,
	}
}

// LogMode implements gorm logger interface
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	return &newLogger
}

// Info implements gorm logger interface
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, data...))
}

// Warn implements gorm logger interface
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, data...))
}

// Error implements gorm logger interface
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, data...))
}

// Trace implements gorm logger interface for SQL logging
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) {
		l.Logger.Error("SQL execution failed",
			zap.String("sql", sql),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.Error(err),
		)
		return
	}

	if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
		l.Logger.Warn("Slow SQL query detected",
			zap.String("sql", sql),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.Duration("threshold", l.SlowThreshold),
		)
		return
	}

	l.Logger.Info("SQL executed",
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	)
}