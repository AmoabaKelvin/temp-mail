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

# Run the API server
run-api:
	go run cmd/api/main.go

# Run the mail server
run-mail:
	go run cmd/mail_server/main.go

# Run the frontend (Next.js UI)
run-ui:
	cd ui && npm run dev

# Run everything (in separate terminals - requires multiple terminals)
run-all:
	@echo "Please run these commands in separate terminals:"
	@echo "make run-api"
	@echo "make run-mail"
	@echo "make run-ui"

# Build the UI
build-ui:
	cd ui && npm run build

# Install UI dependencies
ui-deps:
	cd ui && npm install
