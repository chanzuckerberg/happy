test: ## Run tests for shared module
	go test ./... -v -run ^$(name)

lint: ## Lint the current project
	@golangci-lint run


generate-mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	rm -rf aws/interfaces/mock_*
	rm -rf backend/aws/interfaces/mock_*
	rm -rf workspace_repo/mock_*
	rm -rf workspace_repo/interfaces/mock_*
	rm -rf stack/mock_*
	go generate
	cd backend/aws/interfaces; go generate
	cd workspace_repo; go generate
	cd backend/aws/interfaces; go generate
.PHONY: generate-mocks
