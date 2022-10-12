
dev:
	CONFIG_YAML_DIRECTORY=./ TZ=utc APP_ENV=development go run main.go

test:
	CONFIG_YAML_DIRECTORY=../.. TZ=utc APP_ENV=test go test -v ./... -run ^$(name)

lint:
	golangci-lint run

update-docs:
	swag init --dir=./,../shared/model && swag fmt --dir=./,../shared/model
