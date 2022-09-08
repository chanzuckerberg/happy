# happy-api
An API to encapsulate Happy functionality

## Running the Server
To run the server, simply run:
```
make dev
```

## Running Tests
To run all tests, simply run:
```
make test
```

To run a single test, run:
```
APP_ENV=test go test -v ./... -run ^TestVersionCheckFail$
```
(replace `TestVersionCheckFail` with the name of the test you want to run)

## Running the Linter
`golangci-lint` is used to lint the code in this repo. To run lint locally, run `make lint`.

## Updating Swagger docs
`Fiber` has swagger support which we use to generate documentation (https://github.com/gofiber/swagger).

More information about the declarative comment format can be found [here](https://github.com/swaggo/swag#declarative-comments-format).

Make sure you have `swag` installed. If not, you can install it with:
```
go install github.com/swaggo/swag/cmd/swag@latest
```

After updating annotations update the docs by running:
```
make update-docs
```
