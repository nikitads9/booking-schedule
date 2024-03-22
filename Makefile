include .env
BIN_SCHEDULER := "./bin/bookings"
BIN_AUTH := "./bin/auth"
BIN_NOTIFIER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"

#GIT_HASH := $(shell git log --format="%h" -n 1)
#LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

env:
	set -o allexport && source ./.env && set +o allexport

migrate-up:
	export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSL}"
	sleep 2 && goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" up -v
migrate-down:
	export PG_DSN="host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSL}"
	sleep 2 && goose -dir ${MIGRATION_DIR} postgres "${PG_DSN}" up -v

build: build-bookings build-auth build-scheduler build-sender
build-bookings:
	go build -v -ldflags "-w -s" -o $(BIN_SCHEDULER) ./cmd/bookings/bookings.go
build-auth:
	go build -v -ldflags "-w -s" -o $(BIN_AUTH) ./cmd/auth/auth.go
build-scheduler:
	go build -v -ldflags "-w -s" -o $(BIN_NOTIFIER) ./cmd/scheduler/scheduler.go
build-sender:
	go build -v -ldflags "-w -s" -o $(BIN_SENDER) ./cmd/sender/sender.go

.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init
			go install -v golang.org/x/tools/gopls@latest
			go install -v github.com/swaggo/swag/cmd/swag@latest
			go mod tidy

.PHONY: generate-swag
generate-swag:
	swag init --generalInfo cmd/bookings/bookings.go --parseDependency --parseInternal --tags users,bookings --output ./docs/bookings/
	swag init --generalInfo cmd/auth/auth.go --parseDependency --parseInternal --tags auth --output ./docs/auth/

.PHONY: coverage
coverage:
	go test -race -coverprofile="coverage.out" -covermode=atomic ./...
	go tool cover -html="coverage.out"

PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out