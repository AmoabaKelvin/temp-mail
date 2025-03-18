include .env

MIGRATIONS_PATH ?= ./migrations

.PHONY: start-api
start-api:
	go run cmd/api/main.go

.PHONY: start-mail
start-mail:
	go run cmd/mail_server/main.go

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) "$$name"

.PHONY: migrate-up
migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database $(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database $(DB_ADDR) down
