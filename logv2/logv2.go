package logv2

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"sync"
)

type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
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
	Append(level string, message string)
}

type ConsoleAppender struct{}

func (c *ConsoleAppender) Append(level string, message string) {
	fmt.Printf("[%s] %s\n", level, message)
}

type NATSAppender struct {
	conn    *nats.Conn
	subject string
}

func NewNATSAppender(conn *nats.Conn, subject string) *NATSAppender {
	return &NATSAppender{conn: conn, subject: subject}
}

func (n *NATSAppender) Append(level string, message string) {
	n.conn.Publish(n.subject, []byte(fmt.Sprintf("[%s] %s", level, message)))
}

type DefaultLogger struct {
	appenders []Appender
}

func (l *DefaultLogger) log(level string, message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)
	for _, appender := range l.appenders {
		appender.Append(level, formattedMessage)
	}
}

func (l *DefaultLogger) Debug(message string, args ...interface{}) {
	l.log("DEBUG", message, args...)
}

func (l *DefaultLogger) Info(message string, args ...interface{}) {
	l.log("INFO", message, args...)
}

func (l *DefaultLogger) Warn(message string, args ...interface{}) {
	l.log("WARN", message, args...)
}

func (l *DefaultLogger) Error(message string, args ...interface{}) {
	l.log("ERROR", message, args...)
}

func (l *DefaultLogger) Fatal(message string, args ...interface{}) {
	l.log("FATAL", message, args...)
}

func ExampleUsage() {
	// Setup NATS connection
	nc, _ := nats.Connect(nats.DefaultURL)

	// Get logger factory and add appenders
	loggerFactory := GetLoggerFactory()
	loggerFactory.AddAppender(&ConsoleAppender{})
	loggerFactory.AddAppender(NewNATSAppender(nc, "logs"))

	// Create a logger
	logger := loggerFactory.CreateLogger()

	// Log messages
	logger.Info("Application started")
	logger.Debug("Debugging mode enabled")
	logger.Warn("Low disk space")
	logger.Error("An error occurred")
	logger.Fatal("Application crashed")
}
