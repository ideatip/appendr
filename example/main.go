package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.ideatip.dev.appendr/appenders"
	"go.ideatip.dev.appendr/logging"
	"go.ideatip.dev.appendr/models"
)

func main() {
	nc, _ := nats.Connect("nats://192.168.0.247:4222")

	fileAppender, err := appenders.NewFileAppender("application.log", 10*1024*1024)
	if err != nil {
		fmt.Printf("Error setting up file appender: %v\n", err)
		return
	}
	defer fileAppender.Close()

	// Get logger factory and add appenders
	loggerFactory := logging.GetLoggerFactory()
	loggerFactory.AddAppender(&appenders.ConsoleAppender{})
	loggerFactory.AddAppender(appenders.NewNATSAppender(nc, "logflux"))
	loggerFactory.AddAppender(fileAppender)

	// Create a logger
	logger := loggerFactory.CreateLogger()

	// Log messages
	logger.Info("Application started")
	logger.Debug("Debugging mode enabled", models.Field{Key: "mode", Value: "debug"})

	contextLogger := logger.WithFields(models.Field{Key: "user", Value: "john"}, models.Field{Key: "request_id", Value: "abc123"})
	contextLogger.Warn("Low disk space", models.Field{Key: "available", Value: "100MB"})
	contextLogger.Error("An error occurred", models.Field{Key: "error", Value: "database connection failed"})
	contextLogger.Fatal("Application crashed")
}
