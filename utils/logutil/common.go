package logutil

import (
	"fmt"
	"github.com/rs/zerolog"
)

func GetLogger() *zerolog.Logger {
	return GetZerolog()
}

func LogTraceF(format string, args ...interface{}) {
	GetLogger().Trace().Msgf(format, args...)
}

func LogDebugF(format string, args ...interface{}) {
	GetLogger().Debug().Msgf(format, args...)
}

func LogInfoF(format string, args ...interface{}) {
	GetLogger().Info().Msgf(format, args...)
}

func LogWarnF(format string, args ...interface{}) {
	GetLogger().Warn().Msgf(format, args...)
}

func LogErrorF(format string, args ...interface{}) {
	GetLogger().Error().Msgf(format, args...)
}

func LogPrintf(format string, args ...interface{}) {
	GetLogger().Print(fmt.Sprintf(format, args...))
}
