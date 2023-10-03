package ent

//go:generate go run -mod=mod entc.go
//go:generate go run -mod=mod github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=model -generate=types -o ../../../shared/hapi/model/hapi.gen.go ./openapi.json
//go:generate npm run generate
