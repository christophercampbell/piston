ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
	ARCH = amd64
else
	ifeq ($(ARCH),aarch64)
		ARCH = arm64
	endif
endif

GOBASE := $(shell pwd)
GOOS=$(shell uname -s  | tr '[:upper:]' '[:lower:]')
GOBIN := $(GOBASE)/build
GOCMD := $(GOBASE)/main.go
GOBINARY := piston

GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=$(GOOS)

build: ## build the binary
	$(GOENVVARS)  go build -o $(GOBIN)/$(GOBINARY) $(GOCMD)
.PHONY: build

test: ## run tests
	go test ./...
.PHONY: test

clean: ## clean build artifacts
	rm -rf $(GOBIN)
.PHONY: clean

help: ## prints this help
		@grep -h -E '^[a-zA-Z0-9_-]*:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help

