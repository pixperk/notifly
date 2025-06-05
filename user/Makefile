# --- Config ---
SERVICE := user
DB_URL ?= postgres://user:user@localhost:5432/user?sslmode=disable
MIGRATIONS_DIR := db/migrations
SQLC_CONFIG := sqlc.yaml
MIGRATE := migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)"

# --- Tasks ---
sqlc:
	sqlc generate -f $(SQLC_CONFIG)

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-drop:
	$(MIGRATE) drop -f

run:
	go run cmd/user/main.go


.PHONY: sqlc migrate-create migrate-up migrate-down migrate-drop run
