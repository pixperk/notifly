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
		nats.MaxReconnects(5),
		nats.ReconnectWait(2 * time.Second),
	}

	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return nc, nil
}

func SubscribeToNotifications(nc *nats.Conn, cfg notification.Config) error {
	stream := "NOTIF_STREAM"
	consumer := "notif-consumer"

	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}
	// Ensure the stream exists
	streamInfo, err := js.StreamInfo(stream)
	if err != nil {
		log.Printf("Stream %s not found, creating it", stream)
		// Create the stream if it doesn't exist
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      stream,
			Subjects:  []string{"notifications.*"},
			Storage:   nats.FileStorage,
			Retention: nats.InterestPolicy,
			MaxAge:    7 * 24 * time.Hour, // Keep messages for 7 days max
		})
		if err != nil {
			return fmt.Errorf("failed to create stream: %w", err)
		}
	} else {
		log.Printf("Stream %s found with %d messages", streamInfo.Config.Name, streamInfo.State.Msgs)
	}

	consumerInfo, err := js.ConsumerInfo(stream, consumer)
	if err != nil {
		log.Printf("Consumer %s not found, creating it", consumer)
		// Create a pull-based consumer
		_, err = js.AddConsumer(stream, &nats.ConsumerConfig{
			Durable:           consumer,
			AckPolicy:         nats.AckExplicitPolicy,
			MaxDeliver:        5,
			AckWait:           30 * time.Second,
			FilterSubject:     "notifications.*",
			DeliverPolicy:     nats.DeliverAllPolicy,
			ReplayPolicy:      nats.ReplayInstantPolicy,
			MemoryStorage:     false,
			InactiveThreshold: 24 * time.Hour, // Keep consumer state for 24 hours
		})
		if err != nil {
			return fmt.Errorf("failed to add consumer: %w", err)
		}
	} else {
		log.Printf("Consumer %s found with %d pending, %d ack pending",
			consumerInfo.Name, consumerInfo.NumPending, consumerInfo.NumAckPending)
	}

	// Start a goroutine to pull messages from the consumer
	go processPullConsumer(js, stream, consumer, cfg)

	return nil
}

func processPullConsumer(js nats.JetStreamContext, stream, consumer string, cfg notification.Config) {
	// Create a pull subscription
	sub, err := js.PullSubscribe("notifications.*", consumer, nats.BindStream(stream))
	if err != nil {
		log.Printf("Error creating pull subscription: %v", err)
		return
	}

	// Process messages in a loop
	for {
		// Fetch a batch of messages
		msgs, err := sub.Fetch(10, nats.MaxWait(5*time.Second))
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			log.Printf("Error fetching messages: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, msg := range msgs {

			meta, err := msg.Metadata()
			if err != nil {
				log.Printf("Failed to get message metadata: %v", err)
				_ = msg.Term()
				continue
			}
			deliveryCount := meta.NumDelivered
			log.Printf("Processing message %s, delivery attempt %d/%d",
				msg.Subject, deliveryCount, 5)

			var event common.NotificationEvent
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("JSON unmarshal failed: %v", err)
				_ = msg.Term()
				continue
			}

			err = processJob(event, cfg)
			if err != nil {
				if err == util.ErrInvalidPhoneNumber {
					log.Printf("Invalid phone number for job %v - terminating message", event.NotificationId)
					_ = msg.Term()
				} else if deliveryCount >= 5 {
					log.Printf("Max delivery attempts reached for job %v - terminating message", event.NotificationId)
					_ = msg.Term()
				} else {
					log.Printf("Error processing job %v (attempt %d): %v - will retry",
						event.NotificationId, deliveryCount, err)
					// For retryable errors, we need to explicitly NAK the message
					_ = msg.Nak()
				}
			} else {
				log.Printf("Successfully processed job %v", event.NotificationId)
				// Important: Always explicitly acknowledge successful processing
				err = msg.Ack()
				if err != nil {
					log.Printf("Error acknowledging message: %v", err)
				}
			}
		}
	}
}
