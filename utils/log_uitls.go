package utils

import (
	"fmt"
	"github.com/disgoorg/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lgrisa/lib/utils/lfshook"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func InitLogrus(logSavePath string, logSaveDay uint, loggersLevel string) {
	if path := logSavePath; len(path) > 0 {
		writer, _ := rotatelogs.New(
			path+".%Y%m%d",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithRotationCount(logSaveDay),
		)

		writerMap := lfshook.WriterMap{}
		for _, lv := range logrus.AllLevels {
			writerMap[lv] = writer
		}

		logrus.AddHook(lfshook.NewHook(writerMap))
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	switch strings.ToLower(loggersLevel) {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	default:
		logrus.Infof("未知的log_level: %s, 使用info级别", loggersLevel)
	}

	logrus.Infof("设置日志等级: %s", logrus.GetLevel())
}

func InitLog(logLevel int) {
	log.SetLevel(log.Level(logLevel))
}

func LogTraceF(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func LogDebugF(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func LogInfoF(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func LogWarnF(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func LogErrorF(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func LogPrintf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
