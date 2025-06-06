package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/pixperk/notifly/graphql"
	"github.com/pixperk/notifly/graphql/middleware"
)

func main() {
	cfg, err := graphql.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	server, err := graphql.NewGraphQLServer(cfg.UserURL, cfg.TriggerURL)
	if err != nil {
		log.Fatalf("Error starting GraphQL Server : %v", err)
	}

	http.Handle("/graphql", middleware.AuthFromCookie(middleware.InjectResponseWriter(handler.GraphQL(server.ToExecutableSchema()))))
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Println("GraphQL server started on port 3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
