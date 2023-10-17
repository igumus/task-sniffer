PROJECT_BINARY=tds
PROJECT_BINARY_OUTPUT=out
PROJECT_RELEASER_OUTPUT=dist

.PHONY: all

all: help

## Build:
tidy: ## Tidy project
	@go mod tidy

clean: ## Cleans temporary folder
	@rm -rf ${PROJECT_BINARY_OUTPUT}

build: clean tidy ## Builds project
	@go mod tidy
	@GO111MODULE=on CGO_ENABLED=0 go build -ldflags="-w -s" -o ${PROJECT_BINARY_OUTPUT}/bin/${PROJECT_BINARY} main.go

test-clean: build 
	@go clean -testcache

test-all: test-clean ## Runs all tests 
	@go test -v ./... -race -count=1

pre-commit: test-all ## Checks everything is allright
	@echo "Commit Status: OK"

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  %s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
