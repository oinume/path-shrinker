NAME = path-shrinker
GO_TEST ?= go test -v -race -p=1

.PHONY: all
all: build

.PHONY: install-tools
install-tools:
	@go list -f='{{ join .Imports "\n" }}' ./tools.go | tr -d [ | tr -d ] | xargs -I{} go install {}

.PHONY: bootstrap-lint-tools
bootstrap-lint-tools:
	@cd tools && for tool in $(LINT_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		GO111MODULE=on GOBIN=$(PWD)/tools/bin go install $$tool; \
	done

.PHONY: build
build:
	go build -o bin/$(NAME) ./cmd/$(NAME)/main.go

.PHONY: test
test:
	$(GO_TEST) ./...

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run /app/...

.PHONY: clean
clean:
	${RM} bin/$(NAME)
	${RM} -fr dist/*
