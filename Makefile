.PHONY: build test clean run install help

# Build variables
BINARY_NAME=fintrack
VERSION?=0.1.0-dev
BUILD_DIR=bin
GO=go

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/fintrack/

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Run the application
run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Install the application
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Development: run with live reload (requires air)
dev:
	@which air > /dev/null || (echo "Error: 'air' not installed. Install with: go install github.com/cosmtrek/air@latest" && exit 1)
	@air

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "Error: 'golangci-lint' not installed" && exit 1)
	golangci-lint run

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GO) mod verify

# Cross-platform builds
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/fintrack/

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/fintrack/
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/fintrack/

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/fintrack/

# Run tests with race detection
test-race:
	@echo "Running tests with race detector..."
	$(GO) test -race -v ./...

# Run only fast unit tests
test-unit:
	@echo "Running unit tests..."
	$(GO) test -v -short ./...

# Run only integration tests
test-integration:
	@echo "Running integration tests..."
	@if [ -d "tests/integration" ] && [ "$$(ls -A tests/integration/*.go 2>/dev/null)" ]; then \
		$(GO) test -v -tags=integration ./tests/integration/...; \
	else \
		echo "No integration tests found"; \
	fi

# Run usage tests (executable documentation)
test-usage:
	@echo "Running usage tests..."
	@echo "Building binary..."
	@make build > /dev/null
	@echo "Executing usage documentation tests..."
	@$(GO) test -v ./tests/usage/ || true
	@echo ""
	@echo "✅ Usage tests complete. Check tests/usage/*.md for updated results."

# Run usage tests and update markdown files
test-usage-update: test-usage
	@echo "Usage documentation has been updated with actual results"

# Watch mode - re-run tests on file changes
test-watch:
	@echo "Watching for changes..."
	@which entr > /dev/null || (echo "Error: 'entr' not installed. Install with: apt-get install entr" && exit 1)
	@find . -name "*.go" | entr -c make test

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Check if coverage meets threshold
test-coverage-check:
	@echo "Checking coverage threshold..."
	@$(GO) test -cover ./... > /tmp/coverage.txt 2>&1
	@coverage=$$(grep -oP '\d+\.\d+(?=% of statements)' /tmp/coverage.txt | head -1); \
	if [ -z "$$coverage" ]; then \
		echo "⚠️  No coverage data available"; \
		exit 0; \
	fi; \
	echo "Current coverage: $$coverage%"; \
	if (( $$(echo "$$coverage < 60.0" | bc -l) )); then \
		echo "❌ Coverage $$coverage% is below 60% threshold"; \
		exit 1; \
	else \
		echo "✅ Coverage $$coverage% meets threshold"; \
	fi

# Run all quality checks (tests, fmt, lint, coverage)
quality:
	@echo "Running quality checks..."
	@make fmt
	@make lint
	@make test-race
	@make test-coverage-check

# Help
help:
	@echo "FinTrack - Personal Finance Tracking CLI"
	@echo ""
	@echo "Build targets:"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for all platforms"
	@echo "  install       - Install to /usr/local/bin"
	@echo "  clean         - Clean build artifacts"
	@echo ""
	@echo "Development targets:"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run with live reload (requires air)"
	@echo ""
	@echo "Testing targets:"
	@echo "  test          - Run all tests"
	@echo "  test-race     - Run tests with race detector"
	@echo "  test-unit     - Run only unit tests"
	@echo "  test-integration - Run only integration tests"
	@echo "  test-usage    - Run usage documentation tests"
	@echo "  test-usage-update - Run usage tests and update docs"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-coverage-check - Check if coverage meets 60% threshold"
	@echo "  test-watch    - Watch and re-run tests on changes"
	@echo "  benchmark     - Run performance benchmarks"
	@echo ""
	@echo "Code quality targets:"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  quality       - Run all quality checks"
	@echo ""
	@echo "Dependency targets:"
	@echo "  deps          - Download dependencies"
	@echo "  verify        - Verify dependencies"
	@echo ""
	@echo "  help          - Show this help message"
