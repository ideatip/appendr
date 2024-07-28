package appendr

import (
	"github.com/getsentry/sentry-go"
	"ideatip.dev.appendr"
	"log"
	"time"
)

type SentryAppendr struct{}

func NewSentryAppendr() *SentryAppendr {
	return &SentryAppendr{}
}

func (s *SentryAppendr) Append(level appendr.LogLevel, message string, fields []appendr.Field) {
	if sentry.CurrentHub().Client() == nil {
		log.Fatal("sentry not connected")
	}

	sentry.CaptureEvent(&sentry.Event{
		Message:   message,
		Level:     appendrLevelToSentryLevel(level),
		Extra:     appendr.FieldsToMap(fields),
		Timestamp: time.Now().UTC(),
	})
}

func appendrLevelToSentryLevel(level appendr.LogLevel) sentry.Level {
	switch level {
	case appendr.DEBUG:
		return sentry.LevelDebug
	case appendr.INFO:
		return sentry.LevelInfo
	case appendr.WARN:
		return sentry.LevelWarning
	case appendr.ERROR:
		return sentry.LevelError
	case appendr.FATAL:
		return sentry.LevelFatal
	}

	return sentry.LevelError
}
