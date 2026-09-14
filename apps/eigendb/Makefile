# Makefile for eigendb
# A blazing fast in-memory vector database
# Owner: @ryan-awad
# Maintainers: @ryan-awad

.PHONY: all build run clean test e2e lint

# Default target
all: build

# Build the project
build:
	@mkdir -p dist
	go build -o dist/eigen_db main.go

# Run the project
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf dist/
	rm -rf eigen/* || true

# Run tests
test:
	go test ./... -count=1

# Run end-to-end tests with retry logic
e2e: clean build
	./run_e2e.sh

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
	@echo "  e2e    : Run end-to-end tests"
	@echo "  lint   : Run linter"
	@echo "  all    : Default target (build)"

