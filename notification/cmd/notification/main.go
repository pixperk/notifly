package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pixperk/notifly/common"
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

	queue := make(chan common.NotificationEvent, 100)

	err = util.SubscribeToNotifications(nc, queue)
	if err != nil {
		log.Fatalf("NATS subscription failed: %v", err)
	}

	log.Printf("Started...")
	util.StartWorkerPool(queue, 5, cfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
