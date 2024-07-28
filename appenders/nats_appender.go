package appenders

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"go.ideatip.dev.appendr/models"
	"go.ideatip.dev.appendr/utils"
	"time"
)

type NATSAppender struct {
	conn    *nats.Conn
	subject string
}

func NewNATSAppender(conn *nats.Conn, subject string) *NATSAppender {
	return &NATSAppender{conn: conn, subject: subject}
}

func (n *NATSAppender) Append(level models.LogLevel, message string, fields []models.Field) {
	logEntry := map[string]interface{}{
		"level":     level.String(),
		"message":   message,
		"timestamp": time.Now().UTC(),
		"fields":    utils.FieldsToMap(fields),
	}
	jsonData, _ := json.Marshal(logEntry)
	err := n.conn.Publish(n.subject, jsonData)
	if err != nil {
		return
	}
}
