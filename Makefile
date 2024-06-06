# Variable containing the name of the directory we want to clean
CACHE_DIR := ./cache

# List all subdirectories inside 'cache'
SUBDIRS := $(wildcard $(CACHE_DIR)/*)

# Default target
.PHONY: clean
clean:
	@echo "Removing directories in $(CACHE_DIR)"
	@for dir in $(SUBDIRS); do \
		if [ -d "$$dir" ]; then \
			echo "Removing $$dir"; \
			rm -rf "$$dir"; \
		fi; \
	done

.PHONY: run
run:
	@go run ./main.go

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: build
build:
	@CGO_ENABLED=1 go build -v -o bin/clickup-tui
