BINARY_NAME ?= yunxiao-mcp-server
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Version=$(VERSION) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Commit=$(COMMIT) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Date=$(DATE)

.PHONY: build test format tidy clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/yunxiao-mcp-server

test:
	go test ./...

format:
	gofmt -w $$(find . -path './third-party-projects' -prune -o -name '*.go' -print)

tidy:
	go mod tidy

clean:
	rm -rf bin coverage.out

