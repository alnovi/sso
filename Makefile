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

## init: установка требуемых утилит
.PHONY: init
init:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go get -u github.com/pressly/goose/v3

## migration: новая миграция
.PHONY: migration
migration:
	@read -p "Enter migration name: " migration_name; \
	goose -dir=./scripts/migrations create $$migration_name go

## swag: генерация документации
.PHONY: swag
swag:
	swag init -g ./cmd/server/main.go -p snakecase

## lint: статический анализ
.PHONY: lint
lint:
	golangci-lint run ./...

## test: запуск всех тестов
.PHONY: test
test:
	go test -v -count=1 -coverpkg=./... -coverprofile=./coverage.out ./internal/... ./pkg/... ./tests/...
	go tool cover -html=./coverage.out
	rm ./coverage.out

## test-unit: запуск unit тестов
.PHONY: test-unit
test-unit:
	go test -v -race -count=1 --covermode=atomic -coverpkg=./... -coverprofile=./coverage.out ./internal/... ./pkg/...
	#go tool cover -func=./coverage.out
	go tool cover -html=./coverage.out
	rm ./coverage.out

## test-integration: запуск integration тестов
.PHONY: test-integration
test-integration:
	go test -v -count=1 -coverpkg=./... -coverprofile=./coverage.out ./tests/...
	go tool cover -html=./coverage.out
	rm ./coverage.out
