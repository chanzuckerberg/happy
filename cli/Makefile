# Install
install: install-coverage
.PHONY: install

install-coverage:
	@go install github.com/blend/go-sdk/cmd/coverage@latest
.PHONY: install-coverage

# Test
test:
	@ go test -timeout 10s ./...

coverage: install-coverage
	@coverage --keep-coverage-out --exclude="., *mock*, pkg/backend/aws/interfaces/*" --covermode=atomic --coverprofile=coverage.txt --enforce
.PHONY: coverage

coverage-update: install-coverage
	@coverage --update --keep-coverage-out --exclude="., *mock*, pkg/backend/aws/interfaces/*" --covermode=atomic --coverprofile=coverage.txt
.PHONY: coverage-update

coverage-review: install-coverage
	@go tool cover -html=coverage.txt
.PHONY: coverage-review

lint:
	@ golangci-lint run --verbose
	@ go vet ./...
	@ goimports -w .
.PHONY: lint

fmt:
	@ golangci-lint run --fix

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
