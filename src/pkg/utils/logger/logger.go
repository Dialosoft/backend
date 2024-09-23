package logger

import (
	"io"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger = logrus.New()

func InitLogger() {
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}))

	logger.SetLevel(logrus.InfoLevel)
}

// Info logs an informational message.
func Info(msg string, fields map[string]interface{}) {
	logger.WithFields(logrus.Fields(fields)).Info(msg)
}

// Warn logs a warning message.
func Warn(msg string, fields map[string]interface{}) {
	logger.WithFields(logrus.Fields(fields)).Warn(msg)
}

// Error logs an error message, including context such as file and line number.
func Error(msg string, fields map[string]interface{}) {
	// Capture the caller function, file, and line number
	if pc, file, line, ok := runtime.Caller(1); ok {
		fn := runtime.FuncForPC(pc)
		fields["file"] = file
		fields["line"] = line
		fields["function"] = fn.Name()
	}
	logger.WithFields(logrus.Fields(fields)).Error(msg)
}

// Debug logs a debug message, typically used for low-level system information.
func Debug(msg string, fields map[string]interface{}) {
	logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

// Fatal logs a fatal error message and exits the application.
func Fatal(msg string, fields map[string]interface{}) {
	// Capture caller information before exiting
	if pc, file, line, ok := runtime.Caller(1); ok {
		fn := runtime.FuncForPC(pc)
		fields["file"] = file
		fields["line"] = line
		fields["function"] = fn.Name()
	}
	logger.WithFields(logrus.Fields(fields)).Fatal(msg)
}

// CaptureError logs an error message with an additional error object.
func CaptureError(err error, msg string, fields map[string]interface{}) {
	if err != nil {
		fields["error"] = err.Error()
		Error(msg, fields)
	}
}
