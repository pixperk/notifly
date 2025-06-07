package util

import (
	"encoding/json"
	"fmt"
	"log"

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

func SubscribeToNotifications(nc *nats.Conn, queue chan<- common.NotificationEvent) error {
	subj := "notifications.*"

	_, err := nc.QueueSubscribe(subj, "notif-workers", func(msg *nats.Msg) {
		var event common.NotificationEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Invalid message: %v", err)
			return
		}

		job := event
		queue <- job
	})

	return err
}
