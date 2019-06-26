NAME = path-shrinker

LINT_TOOLS=\
	golang.org/x/lint/golint \
	golang.org/x/tools/cmd/goimports \
	github.com/client9/misspell \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck

.PHONY: all
all: build

.PHONY: bootstrap-lint-tools
bootstrap-lint-tools:
	@cd tools && for tool in $(LINT_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		GO111MODULE=on go install $$tool; \
	done

.PHONY: build
build:
	go build -o bin/$(NAME) ./cmd/$(NAME)/main.go
