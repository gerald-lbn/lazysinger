package log

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const (
	PanicLevel = zerolog.PanicLevel
	FatalLevel = zerolog.FatalLevel
	ErrorLevel = zerolog.ErrorLevel
	WarnLevel  = zerolog.WarnLevel
	InfoLevel  = zerolog.InfoLevel
	DebugLevel = zerolog.DebugLevel
	TraceLevel = zerolog.TraceLevel
)

var (
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
)

func GetLevel() zerolog.Level {
	return zerolog.GlobalLevel()
}

func SetLevel(l zerolog.Level) {
	zerolog.SetGlobalLevel(l)
}

func SetLevelString(l string) {
	level := levelFromString(l)
	zerolog.SetGlobalLevel(level)
}

func levelFromString(l string) zerolog.Level {
	envLevel := strings.ToLower(l)
	switch envLevel {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "trace":
		return TraceLevel
	default:
		return InfoLevel
	}
}

func Fatal() *zerolog.Event {
	return logger.Fatal()
}

func Error() *zerolog.Event {
	return logger.Error()
}

func Warn() *zerolog.Event {
	return logger.Warn()
}

func Log() *zerolog.Event {
	return logger.Log()
}

func Info() *zerolog.Event {
	return logger.Info()
}

func Debug() *zerolog.Event {
	return logger.Debug()
}

func Trace() *zerolog.Event {
	return logger.Trace()
}
