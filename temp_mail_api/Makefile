include .env

MIGRATIONS_PATH ?= ./migrations

.PHONY: start
start:
	go run cmd/api/main.go

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
