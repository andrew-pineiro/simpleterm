BINARY_NAME=sterm
BUILD_DIR=build
OLD_DIR=$(BUILD_DIR)/.old

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: build build-all

# Single build with specified GOOS/GOARCH
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe

# Build for multiple platforms
build-all:
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
        GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
        GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
        GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
        GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
