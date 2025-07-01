package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger() *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if os.Getenv("GIN_MODE") == "release" {
		// Production format (JSON)
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
		})
	} else {
		// Development format (Text with colors)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set output
	logger.SetOutput(os.Stdout)

	// Set global logger
	Logger = logger
	return logger
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		return InitLogger()
	}
	return Logger
}

// LogError logs an error with context
func LogError(err error, message string, fields logrus.Fields) {
	logger := GetLogger()
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["error"] = err.Error()
	logger.WithFields(fields).Error(message)
}

// LogInfo logs an info message with context
func LogInfo(message string, fields logrus.Fields) {
	logger := GetLogger()
	if fields == nil {
		fields = logrus.Fields{}
	}
	logger.WithFields(fields).Info(message)
}

// LogWarn logs a warning message with context
func LogWarn(message string, fields logrus.Fields) {
	logger := GetLogger()
	if fields == nil {
		fields = logrus.Fields{}
	}
	logger.WithFields(fields).Warn(message)
}

// LogDebug logs a debug message with context
func LogDebug(message string, fields logrus.Fields) {
	logger := GetLogger()
	if fields == nil {
		fields = logrus.Fields{}
	}
	logger.WithFields(fields).Debug(message)
} 