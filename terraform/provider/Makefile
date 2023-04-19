
test:
	go test -cover -v ./... -run ^$(name)

lint:
	golangci-lint run

update-docs:
	go generate ./...
