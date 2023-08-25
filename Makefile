include .env
export

MIGRATE := $(shell go env GOPATH)/bin/migrate
MIGRATIONS_PATH = $(DB_PATH)/migrations
PG_CONN_STR="postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_DBNAME)?sslmode=$(PG_SSLMODE)"


start: start-db start-backend

start-db:
	@echo "Запуск базы данных..."
	@sudo service postgresql start

start-backend:
	@echo "Запуск бэкенда..."
	@go run cmd/server/main.go


migrate-up:
	@echo "Запуск миграций..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(PG_CONN_STR) up

migrate-down:
	@echo "Откат миграций..."
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database $(PG_CONN_STR) down

.PHONY:  start-db start-backend migrate-up migrate-down

