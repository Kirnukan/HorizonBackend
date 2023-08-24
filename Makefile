include .env
export

MIGRATE := $(shell go env GOPATH)/bin/migrate
MIGRATIONS_PATH = /home/hako/HorizonPlugin/backend/db/migrations
PG_CONN_STR="postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_DBNAME)?sslmode=$(PG_SSLMODE)"



start: start-backend

start-backend:
	@echo "Запуск бэкенда..."
	@go run cmd/server/main.go


migrate-up:
	@echo "Запуск миграций..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(PG_CONN_STR) up

migrate-down:
	@echo "Откат миграций..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(PG_CONN_STR) down

.PHONY:  start-backend migrate-up migrate-down

