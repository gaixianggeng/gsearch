// log/log.go
package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// 默认为 Info 级别
var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func init() {
	// 设置时间格式 2006-01-02 15:04:05.000
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	// show file name and line number
	logger = logger.With().Caller().Logger()
	// set level
	SetLevel("debug")
	// output to stdout as line
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

// SetLevel 设置日志级别
func SetLevel(level string) {
	switch level {
	case "debug":
		logger = logger.Level(zerolog.DebugLevel)
	case "info":
		logger = logger.Level(zerolog.InfoLevel)
	case "warn":
		logger = logger.Level(zerolog.WarnLevel)
	case "error":
		logger = logger.Level(zerolog.ErrorLevel)
	case "fatal":
		logger = logger.Level(zerolog.FatalLevel)
	case "panic":
		logger = logger.Level(zerolog.PanicLevel)
	default:
		logger = logger.Level(zerolog.InfoLevel)
	}
}
func Debug(args ...interface{}) {
	logger.Debug().Msg(joinWithSpace(args...))
}

func Debugf(format string, args ...interface{}) {
	logger.Debug().Msgf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info().Msg(joinWithSpace(args...))
}

func Infof(format string, args ...interface{}) {
	logger.Info().Msgf(format, args...)
}

func Warn(args ...interface{}) {
	logger.Warn().Msg(joinWithSpace(args...))
}

func Warnf(format string, args ...interface{}) {
	logger.Warn().Msgf(format, args...)
}
func Error(args ...interface{}) {
	logger.Error().Msg(joinWithSpace(args...))
}
func Errorf(format string, args ...interface{}) {
	logger.Error().Msgf(format, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal().Msg(joinWithSpace(args...))
}
func Fatalf(format string, args ...interface{}) {
	logger.Fatal().Msgf(format, args...)
}

func Panic(args ...interface{}) {
	logger.Panic().Msg(joinWithSpace(args...))
}

// WithFields 添加字段到日志中
func WithFields(fields map[string]interface{}) zerolog.Context {
	ctx := logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return ctx
}

func joinWithSpace(args ...any) string {
	var builder strings.Builder
	for i, arg := range args {
		if i > 0 {
			builder.WriteByte(' ')
		}
		fmt.Fprint(&builder, arg)
	}
	return builder.String()
}
