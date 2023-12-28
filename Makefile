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
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go get -u github.com/pressly/goose/v3

## migration: новая миграция
.PHONY: migration
migration:
	@read -p "Введите название миграции: " migration_name; \
	goose -dir=./scripts/migrations create $$migration_name sql

## lint: статический анализ
.PHONY: lint
lint:
	golangci-lint run ./...

## test: запуск тестов
.PHONY: test
test:
	go test -v -race ./...

## test-cover: запуск тестов с покрытием кода
.PHONY: test-cover
test-cover:
	go test -v -race -buildvcs -coverprofile=./coverage.out ./...
	go tool cover -html=./coverage.out

## tidy: обновление зависимостей
.PHONY: tidy
tidy:
	go mod tidy -v

## audit: проверка зависимостей
.PHONY: audit
audit:
	go mod download
	go mod verify

## build: сборка проекта
.PHONY: build
build: audit
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o ./sso-server ./cmd/server
