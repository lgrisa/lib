package logutil

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

func InitLog(logLevel int) {
	logger.Level(zerolog.Level(logLevel))
}

func LogTraceF(format string, args ...interface{}) {
	log.Trace().Msgf(format, args...)
}

func LogDebugF(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

func LogInfoF(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func LogWarnF(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func LogErrorF(format string, args ...interface{}) {
	log.Error().Msgf(format, args...)
}

func LogPrintf(format string, args ...interface{}) {
	log.Print(fmt.Sprintf(format, args...))
}
