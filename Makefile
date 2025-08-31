BIN_FILE := cmd/bin/autodoc
GO_PATH=$(shell go env GOPATH)

# Get git information for version injection
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Build flags with version information
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o ${BIN_FILE} cmd/main.go

.PHONY: test
test:
	go test `go list ./... | grep -v test/` -v -coverprofile cover.out -covermode=atomic

.PHONY: coverage
coverage:
	go tool cover -html=cover.out