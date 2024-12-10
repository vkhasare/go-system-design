# Makefile for building a statically linked Go executable with CGO dependencies

# Output binary name
BINARY_NAME := go-url-shortner

.PHONY: all clean build run

# Default target
all: clean build

# Clean up generated files
clean:
	rm -f $(BINARY_NAME)

# Build the binary
deps:
	go mod tidy

build: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) -ldflags '-extldflags "-static"' cmd/main.go
	chmod +x $(BINARY_NAME)

# Run the application
run:
	$(ENV_VARS) go run $(LDFLAGS) cmd/main.go
