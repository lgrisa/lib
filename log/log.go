package log

import (
	"fmt"
	"github.com/disgoorg/log"
	"github.com/lgrisa/lib/config"
)

func InitLog() {
	log.SetLevel(log.Level(config.StartConfig.Log.LogLevel))
}

func LogTracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func LogDebugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func LogInfof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func LogWarnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func LogErrorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func LogPrintf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}