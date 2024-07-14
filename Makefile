NAME = path-shrinker
GO_TEST ?= go test -v -race -p=1
GOLANGCI_LINT_VERSION = v1.59.1

.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/$(NAME) ./cmd/$(NAME)/main.go

.PHONY: test
test:
	$(GO_TEST) ./...

lint: ## Run golangci-lint
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run --modules-download-mode=readonly /app/...
.PHONY: lint

lint/fix: ## Run golangci-lint with --fix
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run --fix --modules-download-mode=readonly /app/...
.PHONY: lint/fix

lint/version: ## Show golangci-lint version
	@echo $(GOLANGCI_LINT_VERSION)

.PHONY: clean
clean:
	${RM} bin/$(NAME)
	${RM} -fr dist/*
