package util

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func InitJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      "NOTIF_STREAM",
		Subjects:  []string{"notifications.*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	})

	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return nil, fmt.Errorf("stream creation failed: %w", err)
	}

	log.Println("JetStream initialized with stream NOTIF_STREAM")
	return js, nil
}
