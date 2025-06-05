package util

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
)

func publishNotif(nc *nats.Conn, event common.NotificationEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("notifications.%s", event.Type)
	return nc.Publish(subject, data)
}
