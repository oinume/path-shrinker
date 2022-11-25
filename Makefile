NAME = path-shrinker

LINT_TOOLS=\
	golang.org/x/lint/golint \
	golang.org/x/tools/cmd/goimports \
	github.com/client9/misspell \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck \
	github.com/golangci/golangci-lint/cmd/golangci-lint
LINT_PACKAGES = $(shell go list ./...)
FORMAT_PACKAGES = $(foreach pkg,$(LINT_PACKAGES),$(shell go env GOPATH)/$(pkg))

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
	go test -race -v ./

.PHONY: lint
lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.49.0 golangci-lint run /app/...

.PHONY: fmt
fmt:
	tools/bin/goimports -l . | grep -E '.'; test $$? -eq 1
	gofmt -l . | grep -E '.'; test $$? -eq 1

.PHONY: vet
vet:
	go vet -v $(LINT_PACKAGES)

.PHONY: staticcheck
staticcheck:
	tools/bin/staticcheck $(LINT_PACKAGES)

.PHONY: errcheck
errcheck:
	tools/bin/errcheck -ignore 'fmt:[FS]?[Pp]rint*' $(LINT_PACKAGES)

.PHONY: clean
clean:
	${RM} bin/$(NAME)
	${RM} -fr dist/*
