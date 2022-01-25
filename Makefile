# Install
install: install-coverage
.PHONY: install

install-coverage:
	@go install github.com/blend/go-sdk/cmd/coverage@latest
.PHONY: install-coverage

# Test
coverage: install-coverage
	@coverage --keep-coverage-out --exclude="." --covermode=atomic --coverprofile=coverage.txt
.PHONY: coverage

coverage-update: install-coverage
	@coverage --update --keep-coverage-out --covermode=atomic --coverprofile=coverage.txt
.PHONY: install-coverage

lint:
	golangci-lint run -E whitespace --exclude-use-default
.PHONY: lint

# Others
generate-mocks:
	@go get github.com/golang/mock/mockgen@v1.5.0
	@go get -u github.com/aws/aws-sdk-go/...
	@rm -rf mocks/mock_*
	@cd mocks; go generate
	@go mod tidy
.PHONY: generate-mocks
