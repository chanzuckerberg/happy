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
make test name=TestVersionCheckFail
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

## Go Conventions (09/09/2022)

### Overview

Use these resources as a general guide:

* https://github.com/uber-go/guide/blob/master/style.md
* https://go.dev/doc/effective_go
* https://github.com/golang/go/wiki/CommonMistakes
* https://github.com/golang/go/wiki/CodeReviewComments

### Pointers

#### Recievers
For receivers, [default to pointer receivers](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type) unless you are doing performance optimizations:

~~~go
type A struct {

}
func (a *A) myFunc() {

}
~~~

#### Structs

* Until our structs get bigger and more complicated, pass values, not pointers. If you need to modify a struct, consider using a pointer receiver:

~~~go
type myStruct struct {
    value string
}
func myFunc(s myStruct) {

}

func (m  *myStruct) editMyStruct {
    m.value = "blah"
}
~~~

* https://stackoverflow.com/questions/23542989/pointers-vs-values-in-parameters-and-return-values: "Slices, maps, channels, strings, function values, and interface values are implemented with pointers internally, and a pointer to them is often redundant."
* Don't use a pointer to avoid memory allocations or for performance reasons. When we measure out program and find performance issues, we will optimize them based on the usecase.
