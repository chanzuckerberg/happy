# Install
install: install-coverage
.PHONY: install

install-coverage:
	@go install github.com/blend/go-sdk/cmd/coverage@latest
.PHONY: install-coverage

# Build
build:
	go build .
.PHONY: build

# Test
test:
	@ go test ./...

coverage: install-coverage
	@coverage --keep-coverage-out --exclude="., mocks/*" --covermode=atomic --coverprofile=coverage.txt
.PHONY: coverage

coverage-update: install-coverage
	@coverage --update --keep-coverage-out --exclude="., mocks/*" --covermode=atomic --coverprofile=coverage.txt
.PHONY: install-coverage

lint:
	@ golangci-lint run -E whitespace --exclude-use-default
	@ go vet ./...
	@ goimports -w .
.PHONY: lint

# Others
generate-mocks:
	@go install github.com/golang/mock/mockgen@latest
	@go get -u ./...
	@rm -rf mocks/mock_*
	@cd mocks; go generate
	@go mod tidy
.PHONY: generate-mocks
