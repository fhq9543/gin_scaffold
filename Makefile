.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

RELEASE_VERSION = v1.0.0

APP_NAME 		= app
BIN_SERVER  	= ./build/${APP_NAME}
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(RELEASE_VERSION).$(GIT_COUNT).$(GIT_HASH)

all: run

build: clean build-app

build-app:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(BIN_SERVER) ./cmd/${APP_NAME}

run:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./cmd/${APP_NAME}/main.go

migrate: wire
	@go run ./cmd/migrate/main.go # 数据库迁移

wire:
	@wire gen ./app

clean:
	@rm -rf $(BIN_SERVER)
