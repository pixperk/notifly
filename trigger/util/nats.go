package util

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
)

func ConnectNats(natsUrl, clientId string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(clientId),
		nats.MaxReconnects(-1),       // Infinite reconnects
		nats.ReconnectWait(2 * 1000), // 2 seconds wait between reconnects
	}

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return nc, nil
}

func PublishNotif(nc *nats.Conn, event common.NotificationEvent) error {

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
