# Makefile for faissgo
# A Go wrapper for Faiss
# Owner: @ryan-awad
# Maintainers: @ryan-awad

.PHONY: all test lint help

# Default target
all: test	

# Run tests
test:
	go test ./... -count=1 -v

# Run linter (intended for separate workflow, optional for PRs)
lint:
	golangci-lint run --allow-parallel-runners

# Help target
help:
	@echo "Available targets:"
	@echo "  test  : Run all tests (verbose)"
	@echo "  lint  : Run linter"
	@echo "  all   : Default target (test)"
