package logutil

import (
	"github.com/lgrisa/lib/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
	"time"
)

var (
	zeroLogger  zerolog.Logger
	zeroLogOnce sync.Once
)

func GetZerolog() *zerolog.Logger {
	return &zeroLogger
}

func InitZeroLog(logLevel int, isDebugMode bool) {
	initZeroLog(zerolog.Level(logLevel), isDebugMode)
}

func InitZeroLogConfig() {
	initZeroLog(zerolog.Level(config.StartConfig.Log.LogLevel), config.StartConfig.SwitchController.IsDebugMode)
}

func initZeroLog(logLevel zerolog.Level, isDebugMode bool) {
	zeroLogOnce.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339}

		//注意不要在生产环境中使用 ， ConsoleWriter 因为它会大大减慢日志记录的速度。它只是为了帮助在开发应用程序时使日志更易于阅读.
		if !isDebugMode {
			fileLogger := &lumberjack.Logger{
				Filename:   "demo.log",
				MaxSize:    5,
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		zeroLogger = zerolog.New(output).
			Level(logLevel).
			With().
			Timestamp().
			Logger()
	})
}
