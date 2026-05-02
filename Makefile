BINARY_NAME ?= yunxiao-mcp-server
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Version=$(VERSION) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Commit=$(COMMIT) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Date=$(DATE)

.PHONY: build test lint format tidy ci smoke clean coverage

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/yunxiao-mcp-server

test:
	go test ./...

lint:
	go vet ./...
	test -z "$$(gofmt -l cmd internal pkg)"
	@which gocyclo >/dev/null 2>&1 && gocyclo -over 15 $$(find cmd internal pkg -name '*.go' -not -name '*_test.go') || echo "gocyclo not installed, skipping complexity check"

format:
	gofmt -w $$(find . -path './third-party-projects' -prune -o -name '*.go' -print)

tidy:
	go mod tidy

ci: lint
	go mod verify
	go test -race ./...
	$(MAKE) build

smoke: build
	./scripts/smoke.sh

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

clean:
	rm -rf bin coverage.out
