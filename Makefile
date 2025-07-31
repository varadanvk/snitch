.PHONY: build clean install uninstall test run

# Build configuration
BINARY_NAME=snitch
BINARY_PATH=./$(BINARY_NAME)
BUILD_FLAGS=-ldflags="-s -w"

# Default target
all: build

# Build the application
build:
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) .

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_PATH)

# Install the application globally
install:
	go install $(BUILD_FLAGS)

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
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Create distribution directory
dist:
	mkdir -p dist

# Package releases
package: dist build-all
	cd dist && tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd dist && tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd dist && tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd dist && zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe