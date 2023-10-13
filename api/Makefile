VERSION = `cat ../.release-please-manifest.json | jq .api -r`

default:
		@make help

.PHONY: help
help: ## Display makefile target descriptions.
	@printf -- "\Happy Path API - Built to improve your Developer Experience\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

build-linux:
	@echo "Build linux binary"
	@GOARCH=amd64 GOOS=linux go build -o dist/happy-api-linux-${VERSION} main.go

.PHONY: build
build: ## Build only local OS for testing compilation and running locally
	go build -o dist/happy-api-${VERSION}

clean-dist: ## Removing dist directory
	@echo "Cleaning dist directory"
	@rm -rf dist

dev: ## Run app in development mode with current working directory as config
	@CONFIG_YAML_DIRECTORY=./ TZ=utc APP_ENV=development aws-oidc exec --profile czi-si-readonly -- chamber exec happy-ldev-hapi -- go run main.go

dev-docker:
	aws-oidc exec --profile czi-si-readonly -- chamber exec happy-ldev-hapi -- docker-compose --profile hapi up --build

dev-db-shell:
	docker exec -it api-hapi-1 /tools/db-shell.sh

test: ## Run tests for API
	@CONFIG_YAML_DIRECTORY=../.. TZ=utc APP_ENV=test go test ./... -run ^$(name)

lint: ## Lint the current project
	@golangci-lint run

update-docs: ## Update swagger docs
	@echo "Update swagger docs"
	@swag init --dir=./,../shared/model && swag fmt --dir=./,../shared/model

docker: ## Build docker image using scratch and Dockerfile.api in root project directory
	pushd ../ && docker build -t happy-api:dev -f Dockerfile.api .

update-ent-schema:
  ## This will update the ent schema and regenerate the ent code. We need to run `go mod tidy` after the fact to update the go.mod file
	go generate ./pkg/ent && go mod tidy
