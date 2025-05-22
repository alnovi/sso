.SILENT:
.DEFAULT-GOAL:= help

## help: справка
.PHONY: help
help:
	@echo 'Single Sign-On (SSO)'
	@echo ''
	@echo 'Usage:'
	@echo '  make <command>'
	@echo ''
	@echo 'The commands are:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## migration: новая миграция
.PHONY: migration
migration:
	@read -p "Enter migration name: " migration_name; \
	go tool goose -dir=./scripts/migrations create $$migration_name go

## swag: генерация документации
.PHONY: swag
swag:
	@go tool swag init -g ./cmd/server/main.go -p snakecase

## lint: статический анализ
.PHONY: lint
lint:
	go tool golangci-lint run ./...

## lint-fix: статический анализ (авто-исправление)
.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run ./... --fix --timeout 650s

## test: запуск всех тестов
.PHONY: test
test:
	@go test -v -count=1 -coverpkg=./... -coverprofile=./coverage.out ./internal/... ./pkg/... ./tests/...
	@go tool cover -html=./coverage.out
	@rm ./coverage.out

## test-unit: запуск unit тестов
.PHONY: test-unit
test-unit:
	@go test -v -count=1 -coverpkg=./... -coverprofile=./coverage.out ./internal/... ./pkg/...
	@go tool cover -html=./coverage.out
	@rm ./coverage.out

## test-integration: запуск integration тестов
.PHONY: test-integration
test-integration:
	@go test -v -count=1 -coverpkg=./... -coverprofile=./coverage.out ./tests/...
	@go tool cover -html=./coverage.out
	@rm ./coverage.out

## dc-up: start docker containers
.PHONY: dc-up
dc-up:
	docker compose up -d -t=60

## dc-stop: stop docker containers
.PHONY: dc-stop
dc-stop:
	docker compose stop

## dc-exec: enter to node container
.PHONY: dc-exec
dc-exec:
	docker compose exec node sh

## dc-web-dev: start npm run dev
.PHONY: dc-web-dev
dc-web-dev:
	docker compose exec -it node npm run dev

## dc-web-build: start npm run build
.PHONY: dc-web-build
dc-web-build:
	docker compose exec -it node npm run build
