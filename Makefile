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
	@ go test -timeout 10s ./...

coverage: install-coverage
	@coverage --keep-coverage-out --exclude="., mocks/*, pkg/backend/aws/testbackend/mock_*, pkg/stack_mgr/mock_*, pkg/workspace_repo/mock_*" --covermode=atomic --coverprofile=coverage.txt --enforce
.PHONY: coverage

coverage-update: install-coverage
	@coverage --update --keep-coverage-out --exclude="., mocks/*, pkg/backend/aws/testbackend/mock_*, pkg/stack_mgr/mock_*, pkg/workspace_repo/mock_*" --covermode=atomic --coverprofile=coverage.txt
.PHONY: install-coverage

lint:
	@ golangci-lint run --verbose
	@ go vet ./...
	@ goimports -w .
.PHONY: lint

fmt:
	@ golangci-lint run --fix

# Others
generate-mocks:
	go install github.com/golang/mock/mockgen@latest
	rm -rf mocks/mock_*
	rm -rf pkg/backend/aws/interfaces/mock_*
	rm -rf pkg/workspace_repo/mock_*
	cd mocks; go generate
	go generate
	cd pkg/workspace_repo; go generate
.PHONY: generate-mocks
