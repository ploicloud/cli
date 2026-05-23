SHELL := /bin/bash

BIN_DIR := bin
SPEC_URL ?= https://ploi.cloud/docs/api.json
SPEC_FILE := gen/spec.json
GEN_OUT := internal/commands/zz_generated.go

VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
CLIENT_ID ?= 019e5393-312b-732f-8292-7a157d677c1e

LDFLAGS := -X github.com/ploicloud/cli/internal/client.Version=$(VERSION) \
           -X github.com/ploicloud/cli/internal/client.Commit=$(COMMIT) \
           -X github.com/ploicloud/cli/internal/auth.ClientID=$(CLIENT_ID)

.PHONY: sync-spec generate build build-pcctl test lint tidy clean release

sync-spec:
	curl -fsSL $(SPEC_URL) -o $(SPEC_FILE)
	@echo "Updated $(SPEC_FILE) ($$(wc -c < $(SPEC_FILE)) bytes)"

generate:
	go run ./gen/cmdgen $(SPEC_FILE) $(GEN_OUT)

build: generate
	mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/ploicloud ./cmd/ploicloud
	cp $(BIN_DIR)/ploicloud $(BIN_DIR)/pcctl

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR) dist $(GEN_OUT)

release:
	./scripts/release.sh
