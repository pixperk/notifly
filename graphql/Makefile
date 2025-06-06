# Variables
GQLGEN = $(shell go env GOPATH)/bin/gqlgen
GO = go
SERVER = main.go

.PHONY: all install init generate run clean

all: install generate run

install:
	@echo "Installing gqlgen CLI..."
	$(GO) install github.com/99designs/gqlgen@latest

init:
	@echo "Initializing gqlgen schema and config..."
	$(GQLGEN) init

generate:
	@echo "Generating Go code from GraphQL schema..."
	$(GQLGEN) generate

run:
	$(GO) run $(SERVER)
