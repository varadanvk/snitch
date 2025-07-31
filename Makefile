.PHONY: build clean install uninstall test run

# Build configuration
BINARY_NAME=snitch
BINARY_PATH=./$(BINARY_NAME)
BUILD_FLAGS=-ldflags="-s -w"

# Default target
all: build

# Build the application
build:
	@echo "Building snitch binary..."
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) .
	@echo "✓ Build complete: $(BINARY_PATH)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f $(BINARY_PATH)
	@echo "✓ Clean complete"

# Install the application globally
install:
	@echo "Installing snitch globally..."
	go install $(BUILD_FLAGS)
	@echo "✓ Installed to $(shell go env GOPATH)/bin/snitch"

# Uninstall the application
uninstall:
	rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run the application
run: build
	$(BINARY_PATH)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@echo "  → Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	@echo "  → macOS Intel..."
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	@echo "  → macOS Apple Silicon..."
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	@echo "  → Windows AMD64..."
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .
	@echo "✓ Cross-platform builds complete"

# Create distribution directory
dist:
	mkdir -p dist

# Package releases
package: dist build-all
	cd dist && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd dist && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd dist && tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd dist && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe