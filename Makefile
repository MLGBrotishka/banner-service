APP_PATH=./myapp
GO_MAIN_PATH=./cmd/main.go
BUILD_FLAGS=
BUILD_DEPS=
TEST_DEPS=

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

# Local run

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o $(APP_PATH) $(GO_MAIN_PATH) $(BUILD_DEPS)

.PHONY: test
test:
	go test $(TEST_DEPS) ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: build
	$(APP_PATH)

.PHONY: clen
clean:
	go clean
	rm -f $(GO_MAIN_PATH)

.PHONY: all
all: test build

# Docker-compose

.PHONY: docker
docker:
	docker-compose -f docker-compose.yml up

.PHONY: docker-clear
docker-clear:
	docker-compose down --volumes

.PHONY: docker-rebuild
docker-rebuild:
	docker-compose -f docker-compose.yml build
	docker-compose -f docker-compose.yml up

# Linter

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.46.2

.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml

# Get env variables

env:
	export $(cat .env | xargs) && rails c