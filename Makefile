# Makefile
BIN_DIR     := bin
BINARY_NAME := rift
BUILD_OUT   := $(BIN_DIR)/$(BINARY_NAME)

GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*' 2>/dev/null | head -1)
TEST_FILES := $(shell find . -name '*_test.go' -not -path './vendor/*' 2>/dev/null | head -1)
MAIN_GO_FILES := $(shell find ./cmd/rift -name '*.go' -not -path './vendor/*' 2>/dev/null | head -1)

.PHONY: help build lint test clean fmt validate

help:  ## Show this help message
	@echo "Rift Build System"
	@echo "================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

build:  ## Build the rift binary
	@if [ -z '$(MAIN_GO_FILES)' ]; then \
	    echo "xxx Main application directory ./cmd/rift/ does not exist or contains no .go files. Cannot build."; \
	    echo "    Creating a basic main.go structure is recommended."; \
	    exit 1; \
	else \
	    echo ">>> Building $(BINARY_NAME)..."; \
			mkdir -p $(BIN_DIR); \
	    go build -ldflags "$(LDFLAGS)" -o $(BUILD_OUT) ./cmd/rift; \
	fi

lint:  ## Run linter (golangci-lint) - skips if no Go files
	@if [ -z '$(GO_FILES)' ]; then \
	    echo "=== No Go files found, skipping linter."; \
	else \
	    echo ">>> Running linter..."; \
	    golangci-lint run ./...; \
	fi

test:  ## Run tests - skips if no test files
	@if [ -z '$(TEST_FILES)' ]; then \
	    echo "=== No tests found, skipping test execution."; \
	else \
	    echo ">>> Running tests..."; \
	    go test -v -race ./...; \
	fi

fmt:  ## Format code with gofumpt
	@echo ">>> Formatting code..."
	gofumpt -l -w .

clean:  ## Remove built artifacts and clean up
	@echo ">>> Cleaning up..."
	go clean
	rm -rf $(BIN_DIR)

validate: ## Validate project structure requirements
	@echo ">>> Validating project structure..."
	@go run scripts/validate_structure.go
