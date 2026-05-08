# Makefile — mediaplatform-datasource
# Usage: `make help` untuk list semua target

SHELL := /bin/bash

# === Project ===
BINARY        := datasource
MODULE        := github.com/infraLinkit/mediaplatform-datasource
VERSION       ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# === Docker ===
DOCKER_REGISTRY := infralinkit
IMAGE_SERVER    := $(DOCKER_REGISTRY)/mediaplatform-datasource-server
IMAGE_MIGRATE   := $(DOCKER_REGISTRY)/mediaplatform-datasource-migrate
DOCKERFILE_SERVER  := Dockerfile.datasource.server
DOCKERFILE_MIGRATE := Dockerfile.datasource.migrate

# === Go ===
GO          := go
GOFLAGS     := -trimpath
LDFLAGS     := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)
TEST_FLAGS  := -race -count=1
COVERAGE    := coverage.out

# === Tools ===
STATICCHECK := $(shell which staticcheck 2>/dev/null)
GOLINT      := $(shell which golangci-lint 2>/dev/null)

.DEFAULT_GOAL := help

# ============================================================
# Help
# ============================================================
.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ============================================================
# Build
# ============================================================
.PHONY: build
build: ## Build binary ke ./bin/datasource
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o bin/$(BINARY) .

.PHONY: build-linux
build-linux: ## Build binary linux/amd64 (untuk container)
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o bin/$(BINARY)-linux .

.PHONY: install
install: ## Install ke $$GOPATH/bin
	$(GO) install $(GOFLAGS) -ldflags '$(LDFLAGS)' .

# ============================================================
# Run (local)
# ============================================================
.PHONY: run
run: ## Run server (./datasource server)
	$(GO) run . server

.PHONY: run-migrate
run-migrate: ## Run migrate (./datasource migrate)
	$(GO) run . migrate

# ============================================================
# Test
# ============================================================
.PHONY: test
test: ## Run all unit tests
	$(GO) test $(TEST_FLAGS) ./...

.PHONY: test-short
test-short: ## Run short tests (skip slow / integration)
	$(GO) test $(TEST_FLAGS) -short ./...

.PHONY: test-verbose
test-verbose: ## Run tests verbose
	$(GO) test $(TEST_FLAGS) -v ./...

.PHONY: test-pkg
test-pkg: ## Test single package: make test-pkg PKG=./src/helper
	$(GO) test $(TEST_FLAGS) -v $(PKG)

.PHONY: cover
cover: ## Run tests + coverage report
	$(GO) test $(TEST_FLAGS) -coverprofile=$(COVERAGE) ./...
	$(GO) tool cover -func=$(COVERAGE) | tail -1

.PHONY: cover-html
cover-html: cover ## Open HTML coverage report
	$(GO) tool cover -html=$(COVERAGE)

.PHONY: bench
bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem -run=^$$ ./...

# ============================================================
# Lint / static analysis
# ============================================================
.PHONY: vet
vet: ## go vet
	$(GO) vet ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck (install: go install honnef.co/go/tools/cmd/staticcheck@latest)
ifndef STATICCHECK
	@echo "staticcheck not installed. Run: go install honnef.co/go/tools/cmd/staticcheck@latest"
	@exit 1
else
	$(STATICCHECK) ./...
endif

.PHONY: lint
lint: ## Run golangci-lint (install: brew install golangci-lint)
ifndef GOLINT
	@echo "golangci-lint not installed."
	@exit 1
else
	$(GOLINT) run ./...
endif

.PHONY: check
check: vet staticcheck ## Run vet + staticcheck

# ============================================================
# Format
# ============================================================
.PHONY: fmt
fmt: ## go fmt
	$(GO) fmt ./...

.PHONY: tidy
tidy: ## go mod tidy
	$(GO) mod tidy

.PHONY: deps
deps: ## Download deps
	$(GO) mod download

# ============================================================
# Docker
# ============================================================
.PHONY: docker-build
docker-build: docker-build-server docker-build-migrate ## Build kedua image

.PHONY: docker-build-server
docker-build-server: ## Build server image
	docker build --platform=linux/amd64 -t $(IMAGE_SERVER):$(VERSION) -f $(DOCKERFILE_SERVER) .

.PHONY: docker-build-migrate
docker-build-migrate: ## Build migrate image
	docker build --platform=linux/amd64 -t $(IMAGE_MIGRATE):$(VERSION) -f $(DOCKERFILE_MIGRATE) .

.PHONY: docker-push
docker-push: docker-push-server docker-push-migrate ## Push kedua image

.PHONY: docker-push-server
docker-push-server: ## Push server image
	docker push $(IMAGE_SERVER):$(VERSION)

.PHONY: docker-push-migrate
docker-push-migrate: ## Push migrate image
	docker push $(IMAGE_MIGRATE):$(VERSION)

.PHONY: docker-run-server
docker-run-server: ## Run server container local
	docker run --rm -it --env-file .env -p 8081:81 $(IMAGE_SERVER):$(VERSION)

.PHONY: docker-run-migrate
docker-run-migrate: ## Run migrate container local
	docker run --rm --env-file .env $(IMAGE_MIGRATE):$(VERSION)

# ============================================================
# Cleanup
# ============================================================
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/ $(COVERAGE) coverage.html

.PHONY: clean-all
clean-all: clean ## Clean + module cache
	$(GO) clean -cache -testcache -modcache

# ============================================================
# Convenience
# ============================================================
.PHONY: all
all: fmt tidy check test build ## fmt + tidy + check + test + build

.PHONY: ci
ci: deps check test build ## CI pipeline (no fmt/tidy modifications)

.PHONY: version
version: ## Print version info
	@echo "version:    $(VERSION)"
	@echo "build_time: $(BUILD_TIME)"
	@echo "module:     $(MODULE)"
