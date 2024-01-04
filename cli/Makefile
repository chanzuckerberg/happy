
VERSION = `cat ../.release-please-manifest.json | jq .cli -r`

default:
		@make help

.PHONY: help
help: ## Display makefile target descriptions.
	@printf -- "\Happy Path CLI - Built to improve your Developer Experience\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

# Install
install: install-coverage ## Install coverage tool for happy
.PHONY: install

install-coverage: ## Install coverage tool
	@go install github.com/blend/go-sdk/cmd/coverage@latest
.PHONY: install-coverage

# Test
test: ## Run test coverage
	@ go test -timeout 10s ./...

coverage: install-coverage ## Run code coverage tooling and ouput
	@coverage --keep-coverage-out --exclude="., *mock*, pkg/backend/aws/interfaces/*" --covermode=atomic --coverprofile=coverage.txt --enforce
.PHONY: coverage

coverage-update: install-coverage ## Run code coverage update only
	@coverage --update --keep-coverage-out --exclude="., *mock*, pkg/backend/aws/interfaces/*" --covermode=atomic --coverprofile=coverage.txt
.PHONY: coverage-update

coverage-review: install-coverage ## Coverage review convert to HTML file
	@go tool cover -html=coverage.txt
.PHONY: coverage-review

lint: ## Run lint tooling and check imports
	@ golangci-lint run --verbose
	@ go vet ./...
	@ goimports -w .
.PHONY: lint

fmt: ## Run golang lint tool and autofix those that are possible
	@ golangci-lint run --fix

.PHONY: build
build: ## Build only local OS for testing compilation and running locally
	@go build -o dist/happy main.go

build-all: clean-build build-mac-arm build-windows build-mac-intel
	@echo "Built multiple packages, use -j with this command to do in parallel"

build-mac-arm:
	@echo Building mac arm version
	GOARCH=arm64 GOOS=darwin go build -o dist/happy-mac-arm-${VERSION} main.go

build-windows:
	@echo Building windows version
	GOARCH=amd64 GOOS=windows go build -o dist/happy-windows-${VERSION} main.go

build-mac-intel:  ## Build the CLI tool for various platforms
	@echo Building mac intel version
	GOARCH=amd64 GOOS=darwin go build -o dist/happy-mac-intel-${VERSION} main.go


clean-all: clean-build clean-vendor ## Clean `build` and `vendor` directories

clean-vendor: ## Clean vendor directory from `go mod vendor` invocation
	@rm -rf vendor

clean-build: ## Remove all build files from build directory
	@rm -rf dist

# Others
generate-mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	rm -rf mocks/mock_*
	cd mocks; go generate
.PHONY: generate-mocks

