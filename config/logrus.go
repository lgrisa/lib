package config

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lgrisa/library/utils/lfshook"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func InitLogLevel(logSavePath string, logSaveDay uint, loggersLevel string) {
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
