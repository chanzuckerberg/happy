
test:
	TF_ACC=yes HAPPY_API_BASE_URL=https://fake.happy-api.io go test -cover -v ./... -run ^$(name)

lint:
	golangci-lint run

update-docs:
	go generate ./...
