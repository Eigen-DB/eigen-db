# Makefile for hnswgo
# A Go wrapper for hnswlib
# Owner: @ryan-awad
# Maintainers: @ryan-awad

.PHONY: all test lint help

# Default target
all: test

# Run tests
test:
	go test ./... -count=1

# Run linter
lint:
	golangci-lint run --allow-parallel-runners

# Help target
help:
	@echo "Available targets:"
	@echo "  test : Run unit tests"
	@echo "  lint : Run linter"
	@echo "  all  : Default target (test)"
