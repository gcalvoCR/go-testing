package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// Set log level from environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// Set JSON formatter for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Set output to stdout
	Log.SetOutput(os.Stdout)
}

// RequestLogger logs HTTP requests
func RequestLogger(method, path string, statusCode int, duration time.Duration) {
	Log.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration.Milliseconds(),
	}).Info("HTTP Request")
}

// Error logs error messages
func Error(message string, err error) {
	if err != nil {
		Log.WithError(err).Error(message)
	} else {
		Log.Error(message)
	}
}

// Info logs info messages
func Info(message string, fields map[string]interface{}) {
	if fields != nil {
		Log.WithFields(fields).Info(message)
	} else {
		Log.Info(message)
	}
}

// Warn logs warning messages
func Warn(message string, fields map[string]interface{}) {
	if fields != nil {
		Log.WithFields(fields).Warn(message)
	} else {
		Log.Warn(message)
	}
}
