package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"ideatip.dev.appendr"
	natsAppendr "ideatip.dev.appendr/nats"
)

func main() {
	nc, _ := nats.Connect("nats://192.168.0.247:4222")

	fileAppender, err := appendr.NewFileAppender("application.log", 10*1024*1024)
	if err != nil {
		fmt.Printf("Error setting up file appender: %v\n", err)
		return
	}
	defer fileAppender.Close()

	// Get logger factory and add appenders
	loggerFactory := appendr.GetLoggerFactory()
	loggerFactory.AddAppender(&appendr.ConsoleAppender{})
	loggerFactory.AddAppender(natsAppendr.NewNATSAppender(nc, "logflux"))
	loggerFactory.AddAppender(fileAppender)

	// Create a logger
	logger := loggerFactory.CreateLogger()

	// Log messages
	logger.Info("Application started")
	logger.Debug("Debugging mode enabled", appendr.Field{"mode", "debug"})

	contextLogger := logger.WithFields(appendr.Field{"user", "john"}, appendr.Field{"request_id", "abc123"})
	contextLogger.Warn("Low disk space", appendr.Field{"available", "100MB"})
	contextLogger.Error("An error occurred", appendr.Field{"error", "database connection failed"})
	contextLogger.Fatal("Application crashed")
}
