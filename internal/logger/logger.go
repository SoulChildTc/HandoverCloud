package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"soul/global"
	"soul/utils/logutil"
)

var (
	log    *logrus.Entry
	logger *logrus.Logger
)

func getMultiOutput() io.Writer {
	var writer []io.Writer

	logPath := global.Config.Log.Path

	if global.Config.Log.Console {
		writer = append(writer, os.Stdout)
	}

	if !global.Config.Log.CloseFileLog {
		if global.Config.Log.Rotate.Enable {
			writer = append(writer, &lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    global.Config.Log.Rotate.MaxSize,    // 单个文件最大大小, 单位M
				MaxBackups: global.Config.Log.Rotate.MaxBackups, // 最多保留多少个文件
				MaxAge:     global.Config.Log.Rotate.MaxAge,     // 每个最多保留多少天
				Compress:   global.Config.Log.Rotate.Compress,   // 启用压缩
				LocalTime:  global.Config.Log.Rotate.Localtime,  // 默认使用UTC时间, 改为使用本地时间
			})
		} else {
			f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				fmt.Println("创建日志文件失败")
				panic(err.Error())
			}
			writer = append(writer, f)
		}
	}

	return io.MultiWriter(writer...)

}

func InitLogger() {
	logger = logrus.New()

	// 获取日志输出目标
	outputs := getMultiOutput()
	// 设置日志输出目标
	logger.SetOutput(outputs)

	// 输出代码位置
	//logger.SetReportCaller(true)  // 因为封装的原因,自带的输出位置有问题(暂时无法使用),使用log.InfoC可以输出代码位置

	// 设置日志格式
	if global.Config.Env == "dev" {
		logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", PadLevelText: true})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	}

	level, err := logrus.ParseLevel(global.Config.Log.Level)
	if err != nil {
		panic("未知的日志级别,可选项为[TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC]")
	}
	logger.SetLevel(level)

	log = logger.WithFields(logrus.Fields{
		"service": global.Config.AppName,
	})
}

func GetLogger() *logrus.Logger {
	return logger
}

func GetEntry() *logrus.Entry {
	return log
}

func Trace(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Debug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panic(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func TraceC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Tracef(format, args...)
}

func DebugC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Debugf(format, args...)
}

func InfoC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Infof(format, args...)
}

func WarnC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Warnf(format, args...)
}

func ErrorC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Errorf(format, args...)
}

func FatalC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Fatalf(format, args...)
}

func PanicC(format string, args ...interface{}) {
	log.WithField("file", logutil.CallerInfo(2)).Panicf(format, args...)
}
