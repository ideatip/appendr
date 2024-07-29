package appenders

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"go.ideatip.dev/appendr/models"
	"go.ideatip.dev/appendr/utils"
	"log"
	"time"
)

type SentryAppendr struct{}

func NewSentryAppendr() *SentryAppendr {
	return &SentryAppendr{}
}

func (s *SentryAppendr) Append(level models.LogLevel, message string, fields []models.Field) {
	if sentry.CurrentHub().Client() == nil {
		log.Fatal("sentry not connected")
	}

	if level == models.INFO {
		sentry.CaptureEvent(&sentry.Event{
			Message:   message,
			Level:     appendrLevelToSentryLevel(level),
			Extra:     utils.FieldsToMap(fields),
			Timestamp: time.Now().UTC(),
		})
	} else {
		sentry.CaptureException(errors.New(message))
	}
}

func appendrLevelToSentryLevel(level models.LogLevel) sentry.Level {
	switch level {
	case models.DEBUG:
		return sentry.LevelDebug
	case models.INFO:
		return sentry.LevelInfo
	case models.WARN:
		return sentry.LevelWarning
	case models.ERROR:
		return sentry.LevelError
	case models.FATAL:
		return sentry.LevelFatal
	}

	return sentry.LevelError
}
