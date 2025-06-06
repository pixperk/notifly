# Variables
GQLGEN = go run github.com/99designs/gqlgen
GO = go
SERVER = cmd/graphql/main.go

.PHONY: generate run 



generate:
	@echo "Generating Go code from GraphQL schema..."
	$(GQLGEN) generate

run:
	$(GO) run $(SERVER)

