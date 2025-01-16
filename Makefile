# Detect platform and architecture
PLATFORM := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)

# Map architecture names for Go builds
ifeq ($(ARCH), x86_64)
    ARCH := amd64
endif
ifeq ($(ARCH), arm64)
    ARCH := arm64
endif

APP_NAME := goetl
BUILD_DIR := ./bin
INSTALL_DIR := $(HOME)/.local/bin

.PHONY: build clean install uninstall help

build: ## Build the Goetl binary for the current platform
	mkdir -p $(BUILD_DIR)
	GOOS=$(PLATFORM) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(APP_NAME) main.go

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

install: build ## Install the Goetl binary to ~/.local/bin
	install $(BUILD_DIR)/$(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)

uninstall: ## Remove the installed binary from ~/.local/bin
	rm -f $(INSTALL_DIR)/$(APP_NAME)

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
