NAME = path-shrinker
GO_TEST ?= go test -v -race -p=1

.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/$(NAME) ./cmd/$(NAME)/main.go

.PHONY: test
test:
	$(GO_TEST) ./...

lint: ## Run golangci-lint
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:v1.59.0 golangci-lint run --modules-download-mode=readonly /app/...
.PHONY: lint

lint/fix: ## Run golangci-lint with --fix
	docker run --rm -v ${GOPATH}/pkg/mod:/go/pkg/mod -v $(shell pwd):/app -v $(shell go env GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go -w /app golangci/golangci-lint:v1.59.0 golangci-lint run --fix --modules-download-mode=readonly /app/...
.PHONY: lint/fix

.PHONY: clean
clean:
	${RM} bin/$(NAME)
	${RM} -fr dist/*
