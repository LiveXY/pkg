package logx

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func NewDBLogger(cfg LogConfig, name string) logger.Interface {
	return New(cfg)
}

type dblogger struct {
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
}

func New(cfg LogConfig) *dblogger {
	return &dblogger{SkipErrRecordNotFound: true, SlowThreshold: time.Duration(cfg.Slow) * time.Second}
}

func (l *dblogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *dblogger) Info(ctx context.Context, s string, args ...any) {
	//DB.Infof(s, args...)
}

func (l *dblogger) Warn(ctx context.Context, s string, args ...any) {
	//DB.Warnf(s, args...)
}

func (l *dblogger) Error(ctx context.Context, s string, args ...any) {
	//DB.Errorf(s, args...)
}

func (l *dblogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	line := utils.FileWithLineNum()
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		DB.Error(err.Error(), zap.String("line", line), zap.Duration("elapsed", elapsed), zap.String("sql", sql))
		return
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		DB.Warn("慢查询", zap.String("line", line), zap.Duration("elapsed", elapsed), zap.String("sql", sql))
		return
	}
	DB.Debug("", zap.String("line", line), zap.Duration("elapsed", elapsed), zap.String("sql", sql))
}
