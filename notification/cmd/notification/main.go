package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pixperk/notifly/notification"
	"github.com/pixperk/notifly/notification/util"
)

func main() {

	cfg, err := notification.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	nc, err := util.ConnectNats(cfg.NatsURL, cfg.NatsClientID)
	if err != nil {
		log.Fatalf("Cannot connect to NATS: %v", err)
	}

	err = util.SubscribeToNotifications(nc, cfg)
	if err != nil {
		log.Fatalf("NATS subscription failed: %v", err)
	}

	log.Printf("Started...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
