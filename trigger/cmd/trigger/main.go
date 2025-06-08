package main

import (
	"log"

	"github.com/pixperk/notifly/common/auth"
	"github.com/pixperk/notifly/trigger"
	"github.com/pixperk/notifly/trigger/util"
)

func main() {
	cfg, err := trigger.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	nc, err := util.ConnectNats(cfg.NatsURL, cfg.NatsClientID)
	if err != nil {
		log.Fatalf("Cannot connect to NATS: %v", err)
	}

	tokenMaker, err := auth.NewPasetoMaker([]byte(cfg.TokenSymmetricKey))
	if err != nil {
		log.Fatalf("Cannot create token maker: %v", err)
	}

	js, err := util.InitJetStream(nc)
	if err != nil {
		log.Fatalf("JetStream init failed: %v", err)
	}

	service := trigger.NewService(js)

	log.Printf("Trigger service is starting on port %d", cfg.Port)

	if err := trigger.ListenGRPC(service, cfg.Port, tokenMaker); err != nil {
		log.Fatalf("Cannot start gRPC server: %v", err)
	}

}
