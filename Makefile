BINARY_NAME ?= yunxiao-mcp-server
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Version=$(VERSION) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Commit=$(COMMIT) \
	-X github.com/futuretea/yunxiao-mcp-server/pkg/core/version.Date=$(DATE)

.PHONY: build test lint format tidy ci smoke clean coverage docs

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

COVERAGE_PKGS := $(shell go list ./... | grep -v -E '/scripts$$|/cmd/yunxiao-mcp-server$$')

coverage:
	go test -coverprofile=coverage.out $(COVERAGE_PKGS)
	go tool cover -func=coverage.out | tail -1

coverage-check:
	go test -coverprofile=coverage.out $(COVERAGE_PKGS)
	@go tool cover -func=coverage.out | awk 'END {print $$3}' | tr -d '%' | awk '{if ($$1 < 98.0) {print "Coverage " $$1 "% is below 98% threshold"; exit 1} else {print "Coverage " $$1 "% meets threshold"}}'

docs:
	go run scripts/gen-tool-docs.go

OSES = darwin linux windows
ARCHS = amd64 arm64

NPM_VERSION ?= $(shell echo $(shell git describe --tags --always) | sed 's/^v//')

CLEAN_TARGETS += $(foreach os,$(OSES),$(foreach arch,$(ARCHS),bin/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,)))
CLEAN_TARGETS += $(foreach os,$(OSES),$(foreach arch,$(ARCHS),./npm/$(BINARY_NAME)-$(os)-$(arch)/bin/))
CLEAN_TARGETS += ./npm/$(BINARY_NAME)/.npmrc ./npm/$(BINARY_NAME)/LICENSE ./npm/$(BINARY_NAME)/README.md
CLEAN_TARGETS += $(foreach os,$(OSES),$(foreach arch,$(ARCHS),./npm/$(BINARY_NAME)-$(os)-$(arch)/.npmrc))

build-all-platforms:
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		GOOS=$(os) GOARCH=$(arch) go build -ldflags "-s -w $(LDFLAGS)" -o bin/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,) ./cmd/yunxiao-mcp-server; \
	))

npm-copy-binaries: build-all-platforms
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		mkdir -p npm/$(BINARY_NAME)-$(os)-$(arch)/bin; \
		cp bin/$(BINARY_NAME)-$(os)-$(arch)$(if $(findstring windows,$(os)),.exe,) npm/$(BINARY_NAME)-$(os)-$(arch)/bin/; \
	))

npm-publish: npm-copy-binaries
	$(foreach os,$(OSES),$(foreach arch,$(ARCHS), \
		DIRNAME="$(BINARY_NAME)-$(os)-$(arch)"; \
		cd npm/$$DIRNAME; \
		echo '//registry.npmjs.org/:_authToken=$(NPM_TOKEN)' >> .npmrc; \
		jq '.version = "$(NPM_VERSION)"' package.json > tmp.json && mv tmp.json package.json; \
		npm publish --access=public; \
		cd ../..; \
	))
	cp README.md LICENSE ./npm/$(BINARY_NAME)/
	echo '//registry.npmjs.org/:_authToken=$(NPM_TOKEN)' >> ./npm/$(BINARY_NAME)/.npmrc
	jq '.version = "$(NPM_VERSION)"' ./npm/$(BINARY_NAME)/package.json > tmp.json && mv tmp.json ./npm/$(BINARY_NAME)/package.json
	jq '.optionalDependencies |= with_entries(.value = "$(NPM_VERSION)")' ./npm/$(BINARY_NAME)/package.json > tmp.json && mv tmp.json ./npm/$(BINARY_NAME)/package.json
	cd npm/$(BINARY_NAME) && npm publish --access=public

clean:
	rm -rf bin coverage.out $(CLEAN_TARGETS)
