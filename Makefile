.PHONY: all build run clean test docker
BINARY_NAME=./myapp
BINARY_PATH=./cmd/main.go
BUILD_FLAGS=
BUILD_DEPS=
TEST_DEPS=
RUN_DEPS=

build:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(BINARY_PATH) $(BUILD_DEPS)

test:
	go test $(TEST_DEPS) ./...

tidy:
	go mod tidy

run: build
	$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_PATH)

all: test build

docker:
	docker-compose -f ./build/docker-compose.yml up

env:
	export $(cat .env | xargs) && rails c