# Makefile for Golang dev.

# Project-specific settings
BINARY_NAME := $(shell basename "$$PWD")
MAIN_GO := ./cmd/main.go

# Define phony targets to avoid conflicts with files of the same name & improve performance
.PHONY: init main-init ez-init dogo-init clean build run test

# Initial setup
init: main-init ez-init dogo-init clean build

# Install the main.go path
main-init:
	@if [ ! -d $(dir $(MAIN_GO)) ]; then \
		echo 'Creating directory: $(dir $(MAIN_GO))'; \
		mkdir -p $(dir $(MAIN_GO)); \
	fi
	@if [ ! -e $(MAIN_GO) ]; then \
		echo 'Creating default $(MAIN_GO) configuration file...'; \
		echo 'package main\n' > $(MAIN_GO); \
		echo 'import "fmt"\n' >> $(MAIN_GO); \
		echo 'func main() {' >> $(MAIN_GO); \
		echo '    fmt.Println("hello world")' >> $(MAIN_GO); \
		echo '}' >> $(MAIN_GO); \
	fi

# Install the James-Bond utilities toolbox pkg
ez-init:
	go get github.com/bondzai/goez@v0.1.0

# Install the dogo compiler for automatic rebuilds. Create a dogo.json configuration file if it doesn't exist
dogo-init:
	go get github.com/liudng/dogo
	@if [ ! -e dogo.json ]; then \
		echo 'Creating default dogo.json configuration file...'; \
		echo '{' > dogo.json; \
		echo '    "WorkingDir": ".",' >> dogo.json; \
		echo '    "SourceDir": ["."],' >> dogo.json; \
		echo '    "SourceExt": [".c", ".cpp", ".go", ".h"],' >> dogo.json; \
		echo '    "BuildCmd": "go build -o bin/$(BINARY_NAME) $(MAIN_GO)",' >> dogo.json; \
		echo '    "RunCmd": "./bin/$(BINARY_NAME)",' >> dogo.json; \
		echo '    "Decreasing": 1' >> dogo.json; \
		echo '}' >> dogo.json; \
	fi

# Clean the project: remove binary and clean Go cache
clean:
	@echo "  >  Cleaning build cache...\n"
	go clean
	rm -f $(BINARY_NAME)
	go mod tidy

# Build the application: compile the Go code in $(MAIN_GO) into a binary
build:
	@echo "  >  Building binary file...\n"
	go build -o bin/$(BINARY_NAME) $(MAIN_GO)

# Run the application: use dogo for automatic rebuilds on file changes
run:
	@echo "  >  Running application...\n"
	dogo -c dogo.json

# Run tests
test:
	@echo "  >  Running tests...\n"
	go test -v ./...
