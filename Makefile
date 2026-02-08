.PHONY: build clean install test run help

# Binary name
BINARY_NAME=BackItUp

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)
	@echo "Clean complete!"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run the application
run: build
	./$(BINARY_NAME)

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run || go vet ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install dependencies"
	@echo "  test     - Run tests"
	@echo "  run      - Build and run the application"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  help     - Show this help message"
