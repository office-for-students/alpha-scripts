SHELL=bash

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?='mongo'

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

MONGO_URI?='localhost:27017'

build:
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build -o $(BUILD_ARCH)/$(BIN_DIR) main.go
debug:
	HUMAN_LOG=1 go run main.go -mongo-uri=$(MONGO_URI)

.PHONEY: build debug
