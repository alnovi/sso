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
	go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	go install github.com/vektra/mockery/v2@v2.43.0

## migration: новая миграция
.PHONY: migration
migration:
	@read -p "Введите название миграции: " migration_name; \
	goose -dir=./scripts/migrations create $$migration_name go

## lint: статический анализ
.PHONY: lint
lint:
	golangci-lint run ./...

## mockery: генерация mocks
.PHONY: mockery
mockery:
	mockery

## test: запуск всех тестов
.PHONY: test
test:
	go test -v -race -count=1 -coverpkg=./... -coverprofile=./coverage.out ./...
	#go tool cover -html=./coverage.out
	rm ./coverage.out

## test-unit: запуск unit тестов
.PHONY: test-unit
test-unit:
	go test -v -race -count=1 -coverprofile=./coverage.out ./internal/... ./pkg/...
	#go tool cover -func=./coverage.out
	#go tool cover -html=./coverage.out
	rm ./coverage.out

## test-integration: запуск integration тестов
.PHONY: test-integration
test-integration:
	go test -v -coverpkg=./... -count=1 -coverprofile=./coverage.out ./tests/...
	#go tool cover -html=./coverage.out
	rm ./coverage.out

## audit: проверка зависимостей
.PHONY: audit
audit:
	go mod download
	go mod verify

## build: сборка проекта
.PHONY: build
build: audit
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o ./sso-server ./cmd/server