# Makefile for Golang dev.

BINARY_NAME := $(shell basename "$$PWD")
MAIN_GO := ./cmd/main.go

.PHONY: init main-init ez-init dogo-init clean build run gen-gitignore test up_build dev

init: gen-gitignore main-init ez-init dogo-init clean build

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

ez-init:
	go get golang.org/x/tools/gopls@latest
	go get github.com/bondzai/goez@v0.1.0
	go get github.com/robfig/cron/v3@v3.0.0
	go get -u github.com/streadway/amqp
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promauto
	go get github.com/prometheus/client_golang/prometheus/promhttp

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

clean:
	@echo "  >  Cleaning build cache...\n"
	go clean
	rm -f $(BINARY_NAME)
	go mod tidy

build:
	@echo "  >  Building binary file...\n"
	go build -o bin/$(BINARY_NAME) $(MAIN_GO)

run:
	@echo "  >  Running application...\n"
	dogo -c dogo.json

gen-gitignore:
	@echo "  >  Generating .gitignore...\n"
	@echo "# Binaries" > .gitignore
	@echo "bin/" >> .gitignore
	@echo "# OS-specific files" >> .gitignore
	@echo "*.exe" >> .gitignore
	@echo "*.exe~" >> .gitignore
	@echo "*.dll" >> .gitignore
	@echo "*.so" >> .gitignore
	@echo "*.dylib" >> .gitignore
	@echo "# Test binary, built with 'go test -c'" >> .gitignore
	@echo "*.test" >> .gitignore
	@echo "# Output of the go coverage tool, specifically when used with LiteIDE" >> .gitignore
	@echo "*.out" >> .gitignore
	@echo "# Dependency directories (remove the comment below to include it)" >> .gitignore
	@echo "# vendor/" >> .gitignore
	@echo "# custom ignore" >> .gitignore
	@echo "*.env" >> .gitignore
	@echo "*.zip" >> .gitignore

test:
	@echo "  >  Running tests...\n"
	go test -v ./...

up_build:
	@echo "  >  Building binary file...\n"
	docker compose up --build --scale invoker=3

dev:
	@echo "  >  Running application...\n"
	go run $(MAIN_GO)
