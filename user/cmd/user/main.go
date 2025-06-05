package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pixperk/notifly/user"
)

func main() {
	cfg, err := user.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	store := user.NewStore(conn)
	service, err := user.NewService(*store, &cfg)
	if err != nil {
		log.Fatalf("Cannot create service: %v", err)
	}

	log.Printf("User service created successfully, starting gRPC server on port %d", cfg.Port)

	if err := user.ListenGRPC(service, cfg.Port); err != nil {
		log.Fatalf("Cannot start gRPC server: %v", err)
	}

}
