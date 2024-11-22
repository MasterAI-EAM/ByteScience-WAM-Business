package logger

import (
	"ByteScience-WAM-Business/conf"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"time"
)

// Logger 实例
var Logger *logrus.Logger
var GormLogger *LogrusGormLogger

// LogrusGormLogger 适配 GORM 的 Logrus 日志记录器
type LogrusGormLogger struct {
	Logger *logrus.Logger
}

// NewLogger 创建 GORM 日志实例
func NewLogger() error {
	config := conf.GlobalConf.Mysql
	// 创建 logrus 实例
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", config.Level)
	}
	logger.SetLevel(level)

	// 设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Logger = logger
	GormLogger = &LogrusGormLogger{Logger: logger}

	return nil
}

// LogMode 日志钩子实现
func (l *LogrusGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &LogrusGormLogger{Logger: l.Logger}
}

// Trace GORM 的 Trace 日志钩子
func (l *LogrusGormLogger) Trace(ctx context.Context, start time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(start)

	// 如果有错误，则记录错误日志
	if err != nil {
		l.Logger.Errorf("SQL Error: %v | Query: %s | Duration: %vms | Rows: %d", err, sql, elapsed.Milliseconds(), rows)
	} else {
		// 慢查询的日志
		slowThreshold := conf.GlobalConf.Mysql.SlowThreshold
		if elapsed.Seconds() > float64(slowThreshold) {
			l.Logger.Warnf("Slow Query: %s | Duration: %vms | Rows: %d", sql, elapsed.Milliseconds(), rows)
		} else {
			l.Logger.Infof("SQL Query: %s | Duration: %vms | Rows: %d", sql, elapsed.Milliseconds(), rows)
		}
	}
}

// Info 日志打印
func (l *LogrusGormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.Infof(msg, args...)
}

// Warn 日志打印
func (l *LogrusGormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.Warnf(msg, args...)
}

// Error 日志打印
func (l *LogrusGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.Errorf(msg, args...)
}
