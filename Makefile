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
	@coverage --keep-coverage-out --exclude="., mocks/*" --covermode=atomic --coverprofile=coverage.txt --enforce
.PHONY: coverage

coverage-update: install-coverage
	@coverage --update --keep-coverage-out --exclude="., mocks/*" --covermode=atomic --coverprofile=coverage.txt
.PHONY: install-coverage

lint:
	@ golangci-lint run
	@ go vet ./...
	@ goimports -w .
.PHONY: lint

fmt:
	@ golangci-lint run --fix

# Others
generate-mocks:
	go install github.com/golang/mock/mockgen@latest
	rm -rf mocks/mock_*
	rm -rf pkg/backend/aws/testbackend/mock_*
	cd mocks; go generate
	cd pkg/backend/aws/testbackend; go generate
.PHONY: generate-mocks
