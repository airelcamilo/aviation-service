package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.SugaredLogger
)

const (
	INFO_DEBUG_LEVEL string = "info"
)

func init() {
	log = NewLogger(INFO_DEBUG_LEVEL)
}

func GetLogger() *zap.SugaredLogger {
	return log
}

// NewLogger creates a new Zap logger with the specified log level
func NewLogger(logLevel string) *zap.SugaredLogger {
	level := zap.NewAtomicLevel()

	// Parse log level from string
	switch logLevel {
	case "debug":
		level.SetLevel(zap.DebugLevel)
	case INFO_DEBUG_LEVEL:
		level.SetLevel(zap.InfoLevel)
	case "warn":
		level.SetLevel(zap.WarnLevel)
	case "error":
		level.SetLevel(zap.ErrorLevel)
	default:
		level.SetLevel(zap.InfoLevel)
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// For development, use console encoder and write to stderr
	// For production, use JSON encoder
	var encoder zapcore.Encoder
	if os.Getenv("APP_ENV") != "development" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger.Sugar()
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	log.Errorf(msg, args...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}

// With returns a logger with the specified fields
func With(fields ...interface{}) *zap.SugaredLogger {
	return log.With(fields...)
}

// Infow logs an info message with fields
func Infow(msg string, fields ...interface{}) {
	log.Infow(msg, fields...)
}

// Errorw logs an error message with fields
func Errorw(msg string, fields ...interface{}) {
	log.Errorw(msg, fields...)
}

// Fatalw logs a fatal message with fields and exits
func Fatalw(msg string, fields ...interface{}) {
	log.Fatalw(msg, fields...)
}