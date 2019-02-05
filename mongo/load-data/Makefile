SHELL=bash

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?=.

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)

MONGO_URI?='localhost:27017'
RELATIVE_FILE_LOCATION?='files/'

build:
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build -o $(BUILD_ARCH)/$(BIN_DIR)/general-data-builder general-data-builder/main.go
	go build -o $(BUILD_ARCH)/$(BIN_DIR)/institution-builder institution-builder/main.go
	go build -o $(BUILD_ARCH)/$(BIN_DIR)/course-builder course-builder/main.go
debug:
	HUMAN_LOG=1 go run general-data-builder/main.go -mongo-uri=$(MONGO_URI) -relative-file-location=$(RELATIVE_FILE_LOCATION)
	HUMAN_LOG=1 go run institution-builder/main.go -mongo-uri=$(MONGO_URI) -auth-token=$(AUTH_TOKEN) -relative-file-location=$(RELATIVE_FILE_LOCATION)
	HUMAN_LOG=1 go run course-builder/main.go -mongo-uri=$(MONGO_URI) -relative-file-location=$(RELATIVE_FILE_LOCATION)

.PHONEY: build debug
