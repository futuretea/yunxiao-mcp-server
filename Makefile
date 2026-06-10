BINARY_NAME ?= yunxiao
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Version=$(VERSION) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Commit=$(COMMIT) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Date=$(DATE)

.PHONY: build test lint format tidy ci smoke clean coverage docs build-all-platforms npm-copy-binaries npm-publish

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME) ./cmd/yunxiao

test:
	go test ./...

lint:
	go vet ./...
	test -z "$$(gofmt -l cmd internal pkg)"
	@which golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed, skipping"
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

COVERAGE_PKGS := $(shell go list ./... | grep -v -E '/scripts$$|/cmd/yunxiao$$|/internal/cmd$$')

coverage:
	go test -coverprofile=coverage.out $(COVERAGE_PKGS)
	go tool cover -func=coverage.out | tail -1

coverage-check:
	go test -coverprofile=coverage.out $(COVERAGE_PKGS)
	@go tool cover -func=coverage.out | awk 'END {print $$3}' | tr -d '%' | awk '{if ($$1 < 97.9) {print "Coverage " $$1 "% is below 98% threshold"; exit 1} else {print "Coverage " $$1 "% meets threshold"}}'

docs:
	go run scripts/gen-tool-docs.go

OSES = darwin linux windows
ARCHS = amd64 arm64

NPM_VERSION ?= $(shell echo $(shell git describe --tags --always) | sed 's/^v//')
NPM_PUBLISH_FLAGS ?= --access=public

CLEAN_TARGETS += $(foreach os,$(OSES),$(foreach arch,$(ARCHS),bin/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,)))
CLEAN_TARGETS += $(foreach os,$(OSES),$(foreach arch,$(ARCHS),./npm/yunxiao-mcp-server-$(os)-$(arch)/bin/))
CLEAN_TARGETS += ./npm/yunxiao-mcp-server/.npmrc ./npm/yunxiao-mcp-server/LICENSE ./npm/yunxiao-mcp-server/README.md
CLEAN_TARGETS += $(foreach os,$(OSES),$(foreach arch,$(ARCHS),./npm/yunxiao-mcp-server-$(os)-$(arch)/.npmrc))

build-all-platforms:
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		GOOS=$(os) GOARCH=$(arch) go build -ldflags "-s -w $(LDFLAGS)" -o bin/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,) ./cmd/yunxiao; \
	))

npm-copy-binaries: build-all-platforms
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		rm -rf npm/yunxiao-mcp-server-$(os)-$(arch)/bin; \
		mkdir -p npm/yunxiao-mcp-server-$(os)-$(arch)/bin; \
		cp bin/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,) npm/yunxiao-mcp-server-$(os)-$(arch)/bin/$(BINARY_NAME)$(if $(findstring windows,$(os)),.exe,); \
	))

npm-publish: npm-copy-binaries
	test -f README.md
	test -f LICENSE
	@test -n "$$NPM_TOKEN" || (echo "NPM_TOKEN is required"; exit 1)
	@set -e; $(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		DIRNAME="yunxiao-mcp-server-$(os)-$(arch)"; \
		cd npm/$$DIRNAME; \
		printf '%s\n' "//registry.npmjs.org/:_authToken=$$NPM_TOKEN" >> .npmrc; \
		jq '.version = "$(NPM_VERSION)"' package.json > tmp.json; \
		mv tmp.json package.json; \
		npm publish $(NPM_PUBLISH_FLAGS); \
		cd ../..; \
	))
	cp README.md LICENSE ./npm/yunxiao-mcp-server/
	@printf '%s\n' "//registry.npmjs.org/:_authToken=$$NPM_TOKEN" >> ./npm/yunxiao-mcp-server/.npmrc
	jq '.version = "$(NPM_VERSION)"' ./npm/yunxiao-mcp-server/package.json > tmp.json && mv tmp.json ./npm/yunxiao-mcp-server/package.json
	jq '.optionalDependencies |= with_entries(.value = "$(NPM_VERSION)")' ./npm/yunxiao-mcp-server/package.json > tmp.json && mv tmp.json ./npm/yunxiao-mcp-server/package.json
	cd npm/yunxiao-mcp-server && npm publish $(NPM_PUBLISH_FLAGS)

clean:
	rm -rf bin coverage.out $(CLEAN_TARGETS)
