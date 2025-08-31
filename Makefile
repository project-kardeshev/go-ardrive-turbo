.PHONY: test test-verbose test-coverage test-unit test-integration build clean lint fmt help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Build parameters
BINARY_NAME=turbo-sdk
BINARY_PATH=./bin/$(BINARY_NAME)

# Test parameters
TEST_PACKAGES=./pkg/... ./test/...
UNIT_TEST_PACKAGES=./pkg/...
INTEGRATION_TEST_PACKAGES=./test/...

help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the project
	$(GOBUILD) -o $(BINARY_PATH) -v ./examples/

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
	rm -f coverage.out coverage.html

test: ## Run all tests
	$(GOTEST) -v $(TEST_PACKAGES)

test-unit: ## Run unit tests only
	$(GOTEST) -v $(UNIT_TEST_PACKAGES)

test-integration: ## Run integration tests only
	$(GOTEST) -v $(INTEGRATION_TEST_PACKAGES)

test-verbose: ## Run all tests with verbose output
	$(GOTEST) -v -race $(TEST_PACKAGES)

test-coverage: ## Run tests with coverage report
	$(GOTEST) -cover -coverprofile=coverage.out $(TEST_PACKAGES)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-func: ## Run tests and show coverage by function
	$(GOTEST) -cover -coverprofile=coverage.out $(TEST_PACKAGES)
	$(GOCMD) tool cover -func=coverage.out

bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem $(TEST_PACKAGES)

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

fmt: ## Format Go code
	$(GOFMT) -s -w .

fmt-check: ## Check if code is formatted
	@test -z "$$($(GOFMT) -s -l . | tee /dev/stderr)"

mod-tidy: ## Tidy go modules
	$(GOMOD) tidy

mod-verify: ## Verify go modules
	$(GOMOD) verify

deps: ## Download dependencies
	$(GOMOD) download

update-deps: ## Update dependencies
	$(GOGET) -u ./...
	$(GOMOD) tidy

ci: fmt-check mod-verify test-coverage ## Run CI pipeline locally

install-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

dev-setup: install-tools deps ## Set up development environment
	@echo "Development environment setup complete"

# Example commands
run-example: build ## Build and run example
	$(BINARY_PATH)

# Docker commands (for future use)
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

docker-test: ## Run tests in Docker
	docker run --rm -v "$(PWD)":/app -w /app golang:1.23 make test

# Check if required tools are installed
check-tools: ## Check if required development tools are installed
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed. Run 'make install-tools'"; exit 1; }
	@echo "All required tools are installed"

# Security
sec-scan: ## Run security scan (requires gosec)
	gosec ./...

# Generate
generate: ## Run go generate
	$(GOCMD) generate ./...

# All-in-one commands
all: clean fmt mod-tidy test build ## Clean, format, test, and build

release-check: fmt-check mod-verify test-coverage ## Pre-release checks
	@echo "✅ All release checks passed"

.DEFAULT_GOAL := help
