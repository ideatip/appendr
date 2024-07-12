package appendr

import (
	"fmt"
	"sync"
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

func FieldsToMap(fields []Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}
