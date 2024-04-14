PROJECT_DIR = $(shell pwd)
APP_NAME=myapp
APP_PATH=$(PROJECT_DIR)/$(APP_NAME)
GO_MAIN_PATH=$(PROJECT_DIR)/cmd/myapp/main.go

BUILD_FLAGS=
BUILD_DEPS=
TEST_DEPS=

DOCKER_PATH=$(PROJECT_DIR)/docker/server
DOCKER_NAME=avito-backend
ENV_FILE=$(DOCKER_PATH)/.env

DOCKER_TEST_PATH=$(PROJECT_DIR)/docker/tests
DOCKER_TEST_NAME=avito-backend-test
ENV_TEST_FILE=$(DOCKER_TEST_PATH)/.env

DOCKER_CMD=docker-compose --env-file $(ENV_FILE) -p $(DOCKER_NAME) --project-directory $(PROJECT_DIR) -f $(DOCKER_PATH)/docker-compose.yaml
DOCKER_TEST_CMD=docker-compose --env-file $(ENV_TEST_FILE) -p $(DOCKER_TEST_NAME) --project-directory $(PROJECT_DIR) -f $(DOCKER_TEST_PATH)/docker-compose.yaml
PROJECT_BIN = $(PROJECT_DIR)/bin
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

# Local run

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o $(APP_PATH) $(GO_MAIN_PATH) $(BUILD_DEPS)

.PHONY: test-fast
test:
	go test $(PROJECT_DIR)/...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: build
	$(APP_PATH)

.PHONY: clean
clean:
	go clean
	rm -f $(APP_PATH)

.PHONY: all
all: test build

# Docker-compose

.PHONY: docker
docker:
	$(DOCKER_CMD) up 

.PHONY: docker-clear
docker-clear:
	$(DOCKER_CMD) down --volumes

.PHONY: docker-rebuild
docker-rebuild:
	$(DOCKER_CMD) build
	$(DOCKER_CMD) up

.PHONY: docker-test
docker-test:
	$(DOCKER_TEST_CMD) build
	$(DOCKER_TEST_CMD) up --exit-code-from test_go-app
	$(DOCKER_TEST_CMD) down --volumes

.PHONY: docker-test-clear
docker-test-clear:
	$(DOCKER_TEST_CMD) down --volumes


# Linter

.PHONY: .install-linter
.install-linter:
	$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
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
	export $(grep -v '^#' docker/server/.env | xargs)