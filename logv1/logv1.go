package logv1

import (
	"fmt"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger is the main interface for logging
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	WithFields(keyvals ...interface{}) Logger
}

// Appender interface for different log outputs
type Appender interface {
	Append(level LogLevel, msg string, fields map[string]interface{})
}

// ConsoleAppender implements Appender interface for console output
type ConsoleAppender struct{}

func (ca *ConsoleAppender) Append(level LogLevel, msg string, fields map[string]interface{}) {
	fmt.Printf("[%s] %s - %s %v\n", time.Now().Format(time.RFC3339), level, msg, fields)
}

// DefaultLogger implements the Logger interface
type DefaultLogger struct {
	appenders []Appender
	fields    map[string]interface{}
}

func NewLogger(appenders ...Appender) *DefaultLogger {
	return &DefaultLogger{
		appenders: appenders,
		fields:    make(map[string]interface{}),
	}
}

func (l *DefaultLogger) log(level LogLevel, msg string, keyvals ...interface{}) {
	fields := make(map[string]interface{})
	for k, v := range l.fields {
		fields[k] = v
	}
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		}
	}

	for _, appender := range l.appenders {
		appender.Append(level, msg, fields)
	}
}

func (l *DefaultLogger) Debug(msg string, keyvals ...interface{}) {
	l.log(DEBUG, msg, keyvals...)
}

func (l *DefaultLogger) Info(msg string, keyvals ...interface{}) {
	l.log(INFO, msg, keyvals...)
}

func (l *DefaultLogger) Warn(msg string, keyvals ...interface{}) {
	l.log(WARN, msg, keyvals...)
}

func (l *DefaultLogger) Error(msg string, keyvals ...interface{}) {
	l.log(ERROR, msg, keyvals...)
}

func (l *DefaultLogger) WithFields(keyvals ...interface{}) Logger {
	newLogger := &DefaultLogger{
		appenders: l.appenders,
		fields:    make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			newLogger.fields[fmt.Sprint(keyvals[i])] = keyvals[i+1]
		}
	}
	return newLogger
}

// Example usage
func ExampleUsage() {
	consoleAppender := &ConsoleAppender{}
	logger := NewLogger(consoleAppender)

	logger.Info("Application started")
	logger.Debug("Debug message", "key1", "value1", "key2", 42)

	contextLogger := logger.WithFields("user", "john", "request_id", "abc123")
	contextLogger.Warn("User action", "action", "login")
}
