BINARY_NAME=sterm
BUILD_DIR=build

GOOS        ?= $(shell go env GOOS)
GOARCH      ?= $(shell go env GOARCH)
CGO_ENABLED ?= $(shell go env CGO_ENABLED)
BINARY_EXT  ?=

export GOOS
export GOARCH
export CGO_ENABLED

PUB     ?= false
PUB_DIR ?= $(BUILD_DIR)/publish

# Set different command styles based on OS
ifeq ($(OS),Windows_NT)
    MKDIR     := powershell -Command New-Item -ItemType Directory -Force
    CP        := copy /y
    RM        := del /f /q
    FIX_PATH   = $(subst /,\\,$(1))
	BINARY_EXT = .exe
else
    MKDIR    := mkdir -p
    CP       := cp
    RM       := rm -f
    FIX_PATH  = $(1)
endif

.PHONY: publish delete \
        build build-all \
        build-binary build-all-binaries \
        build-linux-amd64 build-linux-arm64 \
        build-windows-amd64 \
        build-darwin-amd64 build-darwin-arm64

# Single build
ifeq ($(PUB),true)
build: delete build-binary publish
build-all: delete build-all-binaries
else
build: delete build-binary
build-all: delete build-all-binaries
endif

# Single build (uses GOOS/GOARCH/CGO_ENABLED from environment)
build-binary:
	go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(BINARY_EXT)

# Build for multiple platforms — works on Windows (cmd.exe) and Unix
# MacOS not support currently.
build-all-binaries: build-linux-amd64 build-linux-arm64 build-windows-amd64 #build-darwin-amd64 build-darwin-arm64

# Per-platform targets call build-binary directly, NOT build, to avoid re-running delete
build-linux-amd64:
	$(MAKE) build-binary GOOS=linux GOARCH=amd64 CGO_ENABLED=0 BINARY_EXT=

build-linux-arm64:
	$(MAKE) build-binary GOOS=linux GOARCH=arm64 CGO_ENABLED=0 BINARY_EXT=

build-windows-amd64:
	$(MAKE) build-binary GOOS=windows GOARCH=amd64 CGO_ENABLED=0 BINARY_EXT=.exe

build-darwin-amd64:
	$(MAKE) build-binary GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 BINARY_EXT=

build-darwin-arm64:
	$(MAKE) build-binary GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 BINARY_EXT=

publish:
	@$(MKDIR) $(call FIX_PATH,$(PUB_DIR))
	$(CP) $(call FIX_PATH,$(BUILD_DIR)/*) $(call FIX_PATH,$(PUB_DIR)/)

delete:
	$(RM) $(call FIX_PATH,$(BUILD_DIR)/*)