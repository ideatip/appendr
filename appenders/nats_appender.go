package appenders

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"go.ideatip.dev/appendr/models"
	"go.ideatip.dev/appendr/utils"
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
	err := n.conn.Publish(fmt.Sprintf("%s.%s", n.subject, level.String()), jsonData)
	if err != nil {
		fmt.Printf("error trying to pulish to nats %s: %s \n", n.subject, err)
		return
	}
}
