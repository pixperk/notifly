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
		nats.MaxReconnects(-1),       // Infinite reconnects
		nats.ReconnectWait(2 * 1000), // 2 seconds wait between reconnects
	}

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return nc, nil
}

func SubscribeToNotifications(nc *nats.Conn, cfg notification.Config) error {
	subj := "notifications.*"
	var maxRetries int
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}

	_, err = js.QueueSubscribe(subj, "notif-workers",
		func(msg *nats.Msg) {
			var event common.NotificationEvent
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				msg.Ack() //Ack to discard bad message
				return
			}

			maxRetries = getMaxRetriesByType(event)

			err := processJob(event, cfg)
			if err != nil {
				if err == util.ErrInvalidPhoneNumber {
					log.Printf("Invalid phone number for job %v: %v", event.NotificationId, err)
					msg.Ack() // Ack to discard bad message
					return
				}
				log.Printf("Error processing job %v: %v", event.NotificationId, err)
				msg.Nak() // Nack to retry the message
				return
			} else {
				log.Printf("Successfully processed job %v", event.NotificationId)
				msg.Ack() // Acknowledge successful processing
			}

		},
		nats.Durable("notif-consumer"),
		nats.AckExplicit(),
		nats.MaxDeliver(maxRetries),
		nats.AckWait(30*time.Second),
		nats.ManualAck())

	return err

}
