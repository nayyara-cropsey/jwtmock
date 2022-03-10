package log

import (
	"log"
	"strings"
)

// Level is used to set the logging level
type Level int

const (
	Debug = iota + 1 // Debug logging level
	Info             // Info logging level
	Warn             // Warn logging level
	Error            // Error logging level
)

func (l Level) String() string {
	var str string

	switch l {
	case Debug:
		str = "DEBUG"
	case Info:
		str = "INFO"
	case Warn:
		str = "WARN"
	case Error:
		str = "ERROR"
	}

	return str
}

// LoggerOption is used to set logger options
type LoggerOption func(*Logger)

// Logger wraps a logging implementation
type Logger struct {
	logger *log.Logger
	level  Level
}

// WithLevel sets the logging level of the logger
func WithLevel(level Level) LoggerOption {
	return func(logger *Logger) {
		logger.level = level
	}
}

// WithLevelStr sets the logging level of the logger as a string
func WithLevelStr(level string) LoggerOption {
	var logLevel Level
	switch strings.ToLower(level) {
	case "info":
		logLevel = Info
	case "warn":
		logLevel = Warn
	case "error":
		logLevel = Error
	default:
		logLevel = Debug
	}

	return func(logger *Logger) {
		logger.level = logLevel
	}
}

// NewLogger is the preferred way to create a new logger instance
func NewLogger(options ...LoggerOption) *Logger {
	logger := &Logger{
		logger: log.Default(),
	}

	for _, option := range options {
		option(logger)
	}

	if logger.level < Debug {
		logger.level = Debug
	}

	if logger.level > Error {
		logger.level = Error
	}

	return logger
}

func (l *Logger) Debug(s string) {
	l.logWithLevel(Debug, s)
}

func (l *Logger) Debugf(s string, args ...interface{}) {
	l.logWithLevel(Debug, s, args...)
}

func (l *Logger) Info(s string) {
	l.logWithLevel(Info, s)
}

func (l *Logger) Infof(s string, args ...interface{}) {
	l.logWithLevel(Info, s, args...)
}

func (l *Logger) Error(s string) {
	l.logWithLevel(Error, s)
}

func (l *Logger) Errorf(s string, args ...interface{}) {
	l.logWithLevel(Error, s, args...)
}

func (l *Logger) Warn(s string) {
	l.logWithLevel(Warn, s)
}

func (l *Logger) Warnf(s string, args ...interface{}) {
	l.logWithLevel(Warn, s, args...)
}

func (l *Logger) logWithLevel(level Level, msg string, args ...interface{}) {
	if l.level <= level {
		l.logger.Printf(level.String()+"  "+msg, args...)
	}
}
