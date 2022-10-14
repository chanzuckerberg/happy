
VERSION = `cat ../version.txt`

default:
	@make help

.PHONY: help
help: ## Display makefile target descriptions.
	@printf -- "\tHappy Path API - Built to improve your Developer Experience\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

dev: ## Run Dev Environment
	@CONFIG_YAML_DIRECTORY=./ TZ=utc APP_ENV=development go run main.go

test: ## Run go test
	@CONFIG_YAML_DIRECTORY=../.. TZ=utc APP_ENV=test go test ./... -run ^$(name)

lint: ## Run golangci-lint across API project
	@golangci-lint run

update-docs: ## Update swagger docs
	@swag init --dir=./,../shared/model && swag fmt --dir=./,../shared/model

build: ## Build only local OS for testing compilation and running locally
	@go build -o build/happy-api main.go

build-all: clean-build ## Build the CLI tool for various platforms
	@echo Building windows version
	@GOARCH=amd64 GOOS=windows go build -o build/happy-api-windows-${VERSION} main.go
	@echo Building mac intel version
	@GOARCH=amd64 GOOS=darwin go build -o build/happy-api-mac-intel-${VERSION} main.go
	@echo Building mac arm version
	@GOARCH=arm64 GOOS=darwin go build -o build/happy-api-mac-arm-${VERSION} main.go

clean-all: clean-build clean-vendor ## Clean `build` and `vendor` directories

clean-vendor: ## Clean vendor directory from `go mod vendor` invocation
	@rm -rf vendor

clean-build: ## Remove all build files from build directory
	@rm -rf build
