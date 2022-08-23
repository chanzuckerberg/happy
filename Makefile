
dev:
	TZ=utc APP_ENV=development go run main.go

test:
	TZ=utc APP_ENV=test go test -v ./...
