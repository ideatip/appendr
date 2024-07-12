package logv3

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

type Logger interface {
	Debug(message string, fields ...Field)
	Info(message string, fields ...Field)
	Warn(message string, fields ...Field)
	Error(message string, fields ...Field)
	Fatal(message string, fields ...Field)
	WithFields(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type LoggerFactory struct {
	appenders []Appender
	mu        sync.Mutex
}

var factoryInstance *LoggerFactory
var once sync.Once

func GetLoggerFactory() *LoggerFactory {
	once.Do(func() {
		factoryInstance = &LoggerFactory{}
	})
	return factoryInstance
}

func (f *LoggerFactory) AddAppender(appender Appender) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.appenders = append(f.appenders, appender)
}

func (f *LoggerFactory) CreateLogger() Logger {
	return &DefaultLogger{appenders: f.appenders}
}

type Appender interface {
	Append(level LogLevel, message string, fields []Field)
}

type ConsoleAppender struct{}

func (c *ConsoleAppender) Append(level LogLevel, message string, fields []Field) {
	fmt.Printf("[%s] %s %s\n", level, message, fieldsToString(fields))
}

type NATSAppender struct {
	conn    *nats.Conn
	subject string
}

func NewNATSAppender(conn *nats.Conn, subject string) *NATSAppender {
	return &NATSAppender{conn: conn, subject: subject}
}

func (n *NATSAppender) Append(level LogLevel, message string, fields []Field) {
	logEntry := map[string]interface{}{
		"level":     level.String(),
		"message":   message,
		"timestamp": time.Now().UTC(),
		"fields":    fieldsToMap(fields),
	}
	jsonData, _ := json.Marshal(logEntry)
	n.conn.Publish(n.subject, jsonData)
}

type DefaultLogger struct {
	appenders []Appender
	fields    []Field
}

func (l *DefaultLogger) log(level LogLevel, message string, fields ...Field) {
	allFields := append(l.fields, fields...)
	for _, appender := range l.appenders {
		appender.Append(level, message, allFields)
	}
}

func (l *DefaultLogger) Debug(message string, fields ...Field) {
	l.log(DEBUG, message, fields...)
}

func (l *DefaultLogger) Info(message string, fields ...Field) {
	l.log(INFO, message, fields...)
}

func (l *DefaultLogger) Warn(message string, fields ...Field) {
	l.log(WARN, message, fields...)
}

func (l *DefaultLogger) Error(message string, fields ...Field) {
	l.log(ERROR, message, fields...)
}

func (l *DefaultLogger) Fatal(message string, fields ...Field) {
	l.log(FATAL, message, fields...)
}

func (l *DefaultLogger) WithFields(fields ...Field) Logger {
	return &DefaultLogger{
		appenders: l.appenders,
		fields:    append(l.fields, fields...),
	}
}

func fieldsToString(fields []Field) string {
	result := ""
	for _, field := range fields {
		result += fmt.Sprintf("%s=%v ", field.Key, field.Value)
	}
	return result
}

func fieldsToMap(fields []Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}

func ExampleUsage() {
	// Setup NATS connection
	nc, _ := nats.Connect("nats://192.168.0.247:4222")

	fileAppender, err := NewFileAppender("application.log", 10*1024*1024)
	if err != nil {
		fmt.Printf("Error setting up file appender: %v\n", err)
		return
	}
	defer fileAppender.Close()

	// Get logger factory and add appenders
	loggerFactory := GetLoggerFactory()
	loggerFactory.AddAppender(&ConsoleAppender{})
	loggerFactory.AddAppender(NewNATSAppender(nc, "logflux"))
	loggerFactory.AddAppender(fileAppender)

	// Create a logger
	logger := loggerFactory.CreateLogger()

	// Log messages
	logger.Info("Application started")
	logger.Debug("Debugging mode enabled", Field{"mode", "debug"})

	contextLogger := logger.WithFields(Field{"user", "john"}, Field{"request_id", "abc123"})
	contextLogger.Warn("Low disk space", Field{"available", "100MB"})
	contextLogger.Error("An error occurred", Field{"error", "database connection failed"})
	contextLogger.Fatal("Application crashed")
}
