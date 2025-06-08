// util/nats.go
package util

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	util "github.com/pixperk/notifly/notification/util/notification"
)

func ConnectNats(natsUrl, clientId string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(clientId),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
	}

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return nc, nil
}

func SubscribeToNotifications(nc *nats.Conn, cfg notification.Config) error {
	subj := "notifications.*"
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}

	_, err = js.QueueSubscribe(subj, "notif-workers",
		buildHandler(cfg),
		nats.Durable("notif-consumer"),
		nats.AckExplicit(),
		nats.ManualAck(),
		nats.MaxDeliver(5),
		nats.AckWait(30*time.Second),
	)

	return err
}

func buildHandler(cfg notification.Config) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var event common.NotificationEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("JSON unmarshal failed: %v", err)
			_ = msg.Term()
			return
		}

		err := processJob(event, cfg)
		if err != nil {
			if err == util.ErrInvalidPhoneNumber {
				log.Printf("Invalid phone number: %v", event.NotificationId)
				_ = msg.Term()
				return
			}

			log.Printf("Retry job %v due to error: %v", event.NotificationId, err)
			_ = msg.Nak()
			return
		}

		log.Printf("Job processed: %v", event.NotificationId)
		_ = msg.Ack()
	}
}
