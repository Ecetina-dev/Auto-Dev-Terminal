# Makefile for Auto-Dev-Terminal
# Provides convenient commands for building, testing, and development

# Variables
BINARY_NAME=auto-dev-terminal
BIN_DIR=bin
BUILD_DIR=$(BIN_DIR)
GO=go
GOFLAGS=-v
PACKAGES=$(shell $(GO) list ./...)
TESTFLAGS=-race -coverprofile=coverage.out -covermode=atomic

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build test test-coverage lint fmt vet clean install uninstall help

# Default target
all: lint test build

## Build
build: ## Build the binary
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cli
	@echo "$(GREEN)Built successfully!$(NC)"

build-all: ## Build for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cli
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/cli
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/cli
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cli
	@echo "$(GREEN)All platforms built!$(NC)"

## Test
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test $(GOFLAGS) $(TESTFLAGS) ./...
	@echo "$(GREEN)Tests passed!$(NC)"

test-unit: ## Run unit tests only
	@echo "$(BLUE)Running unit tests...$(NC)"
	$(GO) test $(GOFLAGS) -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GO) test $(TESTFLAGS) ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## Lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	golangci-lint run ./...

lint-fix: ## Run linter with auto-fix
	@echo "$(BLUE)Running linter with auto-fix...$(NC)"
	golangci-lint run ./... --fix

## Format
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO) vet $(PACKAGES)
	@echo "$(GREEN)Vet passed!$(NC)"

## Clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning...$(NC)"
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Cleaned!$(NC)"

## Install
install: ## Install binary to GOPATH/bin
	$(GO) install ./cmd/cli

## Development
dev: ## Run in development mode
	$(GO) run ./cmd/cli

deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

tidy: ## Tidy dependencies
	$(GO) mod tidy

## CI/CD
ci: lint test build ## Run full CI pipeline locally

## Help
help: ## Show this help message
	@echo "$(BLUE)Auto-Dev-Terminal Makefile$(NC)"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
