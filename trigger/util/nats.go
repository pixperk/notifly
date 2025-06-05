package util

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
)

func PublishNotif( /* ctx context.Context */ nc *nats.Conn, event common.NotificationEvent) error {
	//use context

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("notifications.%s", event.Type)

	msg := &nats.Msg{
		Subject: subject,
		Data:    data,
	}

	return nc.PublishMsg(msg)
}
