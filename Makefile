BINARY_NAME=code-tools-mcp.exe
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go' -type f)
CONFIG_FILE=./config.json

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd

.PHONY: run
run: build
	@echo "Running $(BINARY_NAME) with stdio transport..."
	./$(BUILD_DIR)/$(BINARY_NAME) stdio --config $(CONFIG_FILE)

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run
