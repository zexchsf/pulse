include .env
export

SHELL := /bin/bash

.DEFAULT_GOAL := help

# Starts the API service
.PHONY: api
api:
	@cd ./cmd/api && go run .

# Clean up dependencies
tidy:
	go mod tidy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  api   - Start the API service"
	@echo "  tidy  - Clean up dependencies"