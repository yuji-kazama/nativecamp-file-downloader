.PHONY: build run test clean install lint fmt vet

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=ncfiledownloader
BINARY_PATH=./$(BINARY_NAME)

# Main package
MAIN_PATH=./cmd/ncfiledownloader

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Run the project
run:
	$(GOCMD) run $(MAIN_PATH)/main.go $(ARGS)

# Test the project
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf ./out

# Install dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

# Format code
fmt:
	$(GOFMT) ./...

# Run go vet
vet:
	$(GOVET) ./...

# Run linters
lint: fmt vet
	@echo "Running linters..."

# Install the binary
install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run with example URLs
example:
	./run.sh

# Help
help:
	@echo "Available commands:"
	@echo "  make build         - Build the binary"
	@echo "  make run ARGS=...  - Run the application with arguments"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run linters"
	@echo "  make install       - Install binary to GOPATH/bin"
	@echo "  make build-all     - Build for multiple platforms"
	@echo "  make example       - Run with example URLs using run.sh"