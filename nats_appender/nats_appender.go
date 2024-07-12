package appendr

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"ideatip.dev.appendr"
	"time"
)

type NATSAppender struct {
	conn    *nats.Conn
	subject string
}

func NewNATSAppender(conn *nats.Conn, subject string) *NATSAppender {
	return &NATSAppender{conn: conn, subject: subject}
}

func (n *NATSAppender) Append(level appendr.LogLevel, message string, fields []appendr.Field) {
	logEntry := map[string]interface{}{
		"level":     level.String(),
		"message":   message,
		"timestamp": time.Now().UTC(),
		"fields":    appendr.FieldsToMap(fields),
	}
	jsonData, _ := json.Marshal(logEntry)
	err := n.conn.Publish(n.subject, jsonData)
	if err != nil {
		return
	}
}
