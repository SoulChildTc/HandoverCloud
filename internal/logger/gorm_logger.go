package logger

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"soul/global"
	"soul/utils/logutil"
	"strings"
	"time"
)

type logrusAdapter struct {
	glogger.Config
	reportCaller bool
}

func (l logrusAdapter) LogMode(level glogger.LogLevel) glogger.Interface {
	l.LogLevel = level
	return l
}

func (l logrusAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.reportCaller {
		log.WithContext(ctx).WithField("file", logutil.CallerInfo(3)).Infof(msg, data...)
	} else {
		log.WithContext(ctx).Infof(msg, data...)
	}
}

func (l logrusAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.reportCaller {
		log.WithContext(ctx).WithField("file", logutil.CallerInfo(3)).Warnf(msg, data...)
	} else {
		log.WithContext(ctx).Warnf(msg, data...)
	}
}

func (l logrusAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.reportCaller {
		log.WithContext(ctx).WithField("file", logutil.CallerInfo(3)).Errorf(msg, data...)
	} else {
		log.WithContext(ctx).Errorf(msg, data...)
	}
}

func (l logrusAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	traceStr := "[%.3fms] [rows:%v] %s"
	traceWarnStr := "%s\n[%.3fms] [rows:%v] %s"
	traceErrStr := "%s\n[%.3fms] [rows:%v] %s"

	if l.LogLevel <= glogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Config.LogLevel >= glogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.WithContext(ctx).WithField("file", logutil.CallerInfo(4)).Errorf(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.WithContext(ctx).WithField("file", logutil.CallerInfo(4)).Errorf(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= glogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			log.WithContext(ctx).WithField("file", logutil.CallerInfo(4)).Warnf(traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.WithContext(ctx).WithField("file", logutil.CallerInfo(4)).Warnf(traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == glogger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.WithContext(ctx).WithField("file", logutil.CallerInfo(4)).Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.WithContext(ctx).WithField("file", logutil.CallerInfo(4)).Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func NewGormLogger() glogger.Interface {
	var logLevel glogger.LogLevel
	switch strings.ToLower(global.Config.Database.LogLevel) {
	case "info":
		logLevel = glogger.Info
	case "warn":
		logLevel = glogger.Warn
	case "error":
		logLevel = glogger.Error
	case "silent":
		logLevel = glogger.Silent
	}
	gConfig := glogger.Config{
		SlowThreshold:             time.Second, // 慢 SQL 阈值
		LogLevel:                  logLevel,    // 日志级别
		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
	}

	return logrusAdapter{
		Config:       gConfig,
		reportCaller: global.Config.Database.ReportCaller,
	}
}
