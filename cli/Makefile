
VERSION = `cat ../version.txt`

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

build: ## Build only local OS for testing compilation and running locally
	@go build -o build/happy main.go

build-all: clean-build ## Build the CLI tool for various platforms
	@echo Building windows version
	@GOARCH=amd64 GOOS=windows go build -o build/happy-windows-${VERSION} main.go
	@echo Building mac intel version
	@GOARCH=amd64 GOOS=darwin go build -o build/happy-mac-intel-${VERSION} main.go
	@echo Building mac arm version
	@GOARCH=arm64 GOOS=darwin go build -o build/happy-mac-arm-${VERSION} main.go

clean-all: clean-build clean-vendor ## Clean `build` and `vendor` directories

clean-vendor: ## Clean vendor directory from `go mod vendor` invocation
	@rm -rf vendor

clean-build: ## Remove all build files from build directory
	@rm -rf build

# Others
generate-mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	rm -rf mocks/mock_*
	rm -rf pkg/backend/aws/interfaces/mock_*
	rm -rf pkg/workspace_repo/mock_*
	rm -rf pkg/workspace_repo/stack_mgr/mock_*
	rm -rf pkg/workspace_repo/interfaces/mock_*
	cd mocks; go generate
	go generate
	cd pkg/workspace_repo; go generate
	cd pkg/stack_mgr; go generate
.PHONY: generate-mocks

