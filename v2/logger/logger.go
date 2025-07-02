package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger

type LoggerLogLevel string

const (
	LOG_LEVEL_DEBUG LoggerLogLevel = "debug"
	LOG_LEVEL_INFO                 = "info"
	LOG_LEVEL_WARN                 = "warn"
	LOG_LEVEL_ERROR                = "error"
)

type LogWriter struct {
	logger *logrus.Logger
}

func NewLogWriter() *LogWriter {
	return &LogWriter{logger: Logger}
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	lw.logger.Info(string(p))
	return len(p), nil
}

func Initialize(logLevel LoggerLogLevel) error {
	Logger = logrus.New()

	// Set log format (can be JSON or text)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // Show full timestamp
	})

	// Set log level (you can change to logrus.DebugLevel or others)
	Logger.SetLevel(logrusLogLevel(logLevel))

	return nil
}

func logrusLogLevel(logLevel LoggerLogLevel) logrus.Level {
	var lvl logrus.Level

	switch logLevel {
	case LOG_LEVEL_DEBUG:
		lvl = logrus.DebugLevel
		break
	case LOG_LEVEL_INFO:
		lvl = logrus.InfoLevel
		break
	case LOG_LEVEL_WARN:
		lvl = logrus.WarnLevel
		break
	case LOG_LEVEL_ERROR:
		lvl = logrus.ErrorLevel
		break
	default:
		lvl = logrus.InfoLevel
	}
	return lvl
}



