package logging

import (
	"go.ideatip.dev/appendr/models"
	"sync"
)

type LoggerFactory struct {
	appenders []models.Appender
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

func (f *LoggerFactory) AddAppender(appender models.Appender) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.appenders = append(f.appenders, appender)
}

func (f *LoggerFactory) CreateLogger() models.Logger {
	return &DefaultLogger{appenders: f.appenders}
}

type DefaultLogger struct {
	appenders []models.Appender
	fields    []models.Field
}

func (l *DefaultLogger) log(level models.LogLevel, message string, fields ...models.Field) {
	allFields := append(l.fields, fields...)
	for _, appender := range l.appenders {
		appender.Append(level, message, allFields)
	}
}

func (l *DefaultLogger) Debug(message string, fields ...models.Field) {
	l.log(models.DEBUG, message, fields...)
}

func (l *DefaultLogger) Info(message string, fields ...models.Field) {
	l.log(models.INFO, message, fields...)
}

func (l *DefaultLogger) Warn(message string, fields ...models.Field) {
	l.log(models.WARN, message, fields...)
}

func (l *DefaultLogger) Error(message string, fields ...models.Field) {
	l.log(models.ERROR, message, fields...)
}

func (l *DefaultLogger) Fatal(message string, fields ...models.Field) {
	l.log(models.FATAL, message, fields...)
}

func (l *DefaultLogger) WithFields(fields ...models.Field) models.Logger {
	return &DefaultLogger{
		appenders: l.appenders,
		fields:    append(l.fields, fields...),
	}
}
