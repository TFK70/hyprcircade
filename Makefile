ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
MAIN_GO := $(ROOT_DIR)/cmd/hyprcircade/main.go
OUT := $(ROOT_DIR)/build/hyprcircade

.SILENT:
.ONESHELL:
build:
	$(info ==== building $(MAIN_GO) ====)
	go build -o $(OUT) $(MAIN_GO)

.SILENT:
.ONESHELL:
test:
	$(info ==== Running tests ====)
	go test ./...
