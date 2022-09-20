.DEFAULT_GOAL := help

dev:
	TZ=utc APP_ENV=development go run main.go

test:
	TZ=utc APP_ENV=test go test -v ./...

lint:
	golangci-lint run


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'