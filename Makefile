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
lint: fmt vet staticcheck errcheck

.PHONY: fmt
fmt:
	tools/bin/goimports -l . | grep -E '.'; test $$? -eq 1
	#tools/bin/gofmt -w $(LINT_PACKAGES) | grep -E '.'; test $$? -eq 1

.PHONY: vet
vet:
	go vet -v $(LINT_PACKAGES)

.PHONY: staticcheck
staticcheck:
	tools/bin/staticcheck $(LINT_PACKAGES)

.PHONY: errcheck
errcheck:
	tools/bin/errcheck -ignore 'fmt:[FS]?[Pp]rint*' $(LINT_PACKAGES)
