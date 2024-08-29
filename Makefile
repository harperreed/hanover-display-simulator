# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=hanover-simulator
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
BUILD_FLAGS=-v

# Test flags
TEST_FLAGS=-v

# Main package path
MAIN_PACKAGE=.

all: test build

build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

test:
	$(GOTEST) $(TEST_FLAGS) ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)
	./$(BINARY_NAME)

deps:
	$(GOGET) -v -t -d ./...
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_UNIX) -v $(MAIN_PACKAGE)

docker-build:
	docker build -t $(BINARY_NAME) .

# Generate test coverage
cover:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Run golangci-lint
lint:
	golangci-lint run

# Format code
fmt:
	gofmt -s -w .

# Generate documentation
doc:
	godoc -http=:6060

# Perform all pre-commit checks
check: fmt lint test

# Install project binary
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_PACKAGE)

.PHONY: all build test clean run deps build-linux docker-build cover lint fmt doc check install
