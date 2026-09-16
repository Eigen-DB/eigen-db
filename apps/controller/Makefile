# Makefile for controller
# Kubernetes controller for managing customer instances on Eigen Cloud
# Owner: @ryan-awad
# Maintainers: @ryan-awad

.PHONY: all build run clean test lint help

# Default target
all: build

# Build the project
build:
	@mkdir -p dist
	go build -o dist/controller main.go

# Run the project
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf dist/

# Run tests
test:
	go test ./... -count=1

# Run linter
lint:
	golangci-lint run --allow-parallel-runners

# Help target
help:
	@echo "Available targets:"
	@echo "  build  : Build the project"
	@echo "  run    : Run the application"
	@echo "  clean  : Remove build artifacts"
	@echo "  test   : Run unit tests"
	@echo "  lint   : Run linter"
	@echo "  all    : Default target (build)"
