include .env
BIN_SCHEDULER := "./bin/events"
BIN_NOTIFIER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"

DOCKER_IMG="schedule:develop"

#GIT_HASH := $(shell git log --format="%h" -n 1)
#LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
prepare-env:
	set -o allexport && source ./.env && set +o allexport

migrate-up:
	export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSL}"
	sleep 2 && goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" up -v
migrate-down:
	export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSL}"
	sleep 2 && goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" up -v

build: build-events build-scheduler build-sender
build-events:
	go build -v -o $(BIN_SCHEDULER) ./cmd/events/events.go
build-scheduler:
	go build -v -o $(BIN_NOTIFIER) ./cmd/scheduler/scheduler.go
build-sender:
	go build -v -o $(BIN_SENDER) ./cmd/sender/sender.go

.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init
			go install -v golang.org/x/tools/gopls@latest
			go install -v github.com/swaggo/swag/cmd/swag@latest
#go install github.com/joho/godotenv/cmd/godotenv@latest TODO: move all configs to env
			go mod tidy

.PHONY: generate-swag
generate-swag:
	swag init --generalInfo cmd/events/events.go --parseDependency --parseInternal

.PHONY: coverage
coverage:
	go test -race -coverprofile="coverage.out" -covermode=atomic ./...
	go tool cover -html="coverage.out"

PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out