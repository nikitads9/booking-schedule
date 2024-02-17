BIN_SCHEDULER := "./bin/schedule"
BIN_NOTIFIER := "./bin/notifier"
BIN_SENDER := "./bin/sender"

DOCKER_IMG="schedule:develop"

#GIT_HASH := $(shell git log --format="%h" -n 1)
#LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init
			go install -v golang.org/x/tools/gopls@latest
			go get -u github.com/go-chi/chi/v5
			go get -u github.com/go-chi/render
			go install github.com/swaggo/swag/cmd/swag@latest
			go mod tidy

.PHONY: generate-swag
generate-swag:
	swag init --generalInfo cmd/server/schedule.go --parseDependency --parseInternal

.PHONY: coverage
coverage:
	go test -race -coverprofile="coverage.out" -covermode=atomic ./...
	go tool cover -html="coverage.out"

PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out