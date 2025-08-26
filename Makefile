# Makefile
BINARY_NAME=rift

GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*' 2>/dev/null | head -1)
TEST_FILES := $(shell find . -name '*_test.go' -not -path './vendor/*' 2>/dev/null | head -1)
MAIN_GO_FILES := $(shell find ./cmd/rift -name '*.go' -not -path './vendor/*' 2>/dev/null | head -1)

.PHONY: build lint test clean fmt

build:
	@if [ -z '$(MAIN_GO_FILES)' ]; then \
	    echo "❌ Main application directory ./cmd/rift/ does not exist or contains no .go files. Cannot build."; \
	    echo "   Creating a basic main.go structure is recommended."; \
	    exit 1; \
	else \
	    echo "Building $(BINARY_NAME)..."; \
	    go build -o $(BINARY_NAME) ./cmd/rift/; \
	fi

lint:
	@if [ -z '$(GO_FILES)' ]; then \
	    echo "No Go files found, skipping linter."; \
	else \
	    echo "Running linter..."; \
	    golangci-lint run ./...; \
	fi

test:
	@if [ -z '$(TEST_FILES)' ]; then \
	    echo "No tests found, skipping test execution."; \
	else \
	    echo "Running tests..."; \
	    go test -v -race ./...; \
	fi

fmt:
	@echo "Formatting code..."
	gofumpt -l -w .

clean:
	@echo "Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)
