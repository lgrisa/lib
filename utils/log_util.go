package utils

import "github.com/disgoorg/log"

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
