# Makefile for performance-metrics
# Visualize performance metrics for EigenDB
# Owner: @ryan-awad
# Maintainers: @ryan-awad

.PHONY: all setup_venv run clean help

# Default target
all: run

# Set up the virtual environment
setup_venv:
	python3 -m virtualenv .venv && source .venv/bin/activate && pip3 install -r requirements.txt

# Run the dev server (depends on setup_venv)
run: setup_venv
	source .venv/bin/activate && python3 -m streamlit run main.py

# Clean cache artifacts
clean:
	rm -rf **/__pycache__ || true

# Help target
help:
	@echo "Available targets:"
	@echo "  setup_venv : Set up the virtual environment"
	@echo "  run        : Run the app (depends on setup_venv)"
	@echo "  clean      : Remove __pycache__ artifacts"
	@echo "  all        : Default target (run)"
