.PHONY: all generate fmt build test

GO ?= go
GODOC ?= $(GO) run -v golang.org/x/tools/cmd/godoc@latest
GODOC_PORT ?= 6060
GOFMT ?= gofmt
GOFMT_FLAGS = -w -l -s
GOGENERATE_FLAGS = -v

all: generate fmt build

doc:
	$(GODOC) -http :$(GODOC_PORT)

fmt:
	@find . -name '*.go' | xargs -r $(GOFMT) $(GOFMT_FLAGS)
	$(GO) mod tidy || true

generate:
	@git grep -l '^//go:generate' | xargs -r -n1 dirname | sort -u | while read d; do \
		git grep -l '^//go:generate' "$$d"/*.go | xargs -r $(GO) generate $(GOGENERATE_FLAGS); \
	done

build:
	$(GO) get -v ./...

test:
	$(GO) test -v ./...
