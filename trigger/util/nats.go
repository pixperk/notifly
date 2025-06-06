package util

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
)

func PublishNotif(nc *nats.Conn, event common.NotificationEvent) (string, error) {

	notifID := uuid.New().String()

	data, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	subject := fmt.Sprintf("notifications.%s", event.Type)

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
	}

	err = nc.PublishMsg(msg)
	return notifID, err
}
