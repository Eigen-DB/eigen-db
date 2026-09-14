# Makefile for eigen-py
# EigenDB's Official Python API
# Owner: @ryan-awad
# Maintainers: @ryan-awad

.PHONY: all setup_venv build clean help

# Default target
all: build

# Set up the virtual environment
setup_venv:
	python3 -m virtualenv .venv && source .venv/bin/activate && pip3 install setuptools && python3 setup.py install

# Build the project (depends on setup_venv)
build: setup_venv
	python3 setup.py build

# Clean build artifacts
clean:
	rm -rf build/ && rm -rf dist/ && rm -rf **/*.egg-info && rm -rf **/__pycache__ || true

# Help target
help:
	@echo "Available targets:"
	@echo "  setup_venv : Set up the virtual environment"
	@echo "  build      : Build the project (depends on setup_venv)"
	@echo "  clean      : Remove build artifacts"
	@echo "  all        : Default target (build)"
