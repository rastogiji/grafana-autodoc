BIN_FILE := cmd/bin/autodoc
GO_PATH=$(shell go env GOPATH)
.PHONY: build
build:
	go build -o ${BIN_FILE}  cmd/main.go


.PHONY: test
test:
	go test `go list ./... | grep -v test/` -v -coverprofile cover.out -covermode=atomic

.PHONY: coverage
coverage:
	go tool cover -html=cover.out