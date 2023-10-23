//go:build ignore

package main

// NOTE:
// This file is ignored by go build and only used for code generation. Because of this ogent is not in the go.mod file.
// At time of writing this was using ogent v0.0.0-20230621041143-ed3e5d4da458

import (
	"log"

	"ariga.io/ogent"
	"entgo.io/contrib/entoas"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
)

// pagination paramaters
func paginationParameters() []*ogen.Parameter {
	var min int64 = 1
	return []*ogen.Parameter{
		ogen.NewParameter().
			InQuery().
			SetName("page").
			SetDescription("what page to render").
			SetSchema(ogen.Int().SetMinimum(&min)),
		ogen.NewParameter().
			InQuery().
			SetName("itemsPerPage").
			SetDescription("item count to render per page").
			SetSchema(ogen.Int().SetMinimum(&min)),
	}
}

func appEnvStackQueryParameters() []*ogen.Parameter {
	return []*ogen.Parameter{
		ogen.NewParameter().
			InQuery().
			SetName("app_name").
			SetRequired(true).
			SetSchema(ogen.String()),
		ogen.NewParameter().
			InQuery().
			SetName("environment").
			SetRequired(true).
			SetSchema(ogen.String()),
		ogen.NewParameter().
			InQuery().
			SetName("stack").
			SetRequired(false).
			SetSchema(ogen.String()),
	}
}

func main() {
	spec := new(ogen.Spec)
	oas, err := entoas.NewExtension(
		entoas.Spec(spec),
		entoas.Mutations(func(graph *gen.Graph, spec *ogen.Spec) error {
			spec.AddPathItem("/health", ogen.NewPathItem().
				SetDescription("Check if the server is up").
				SetGet(ogen.NewOperation().
					SetOperationID("Health").
					SetSummary("Simple endpoint to check if the server is up").
					AddResponse(
						"200",
						ogen.
							NewResponse().
							SetDescription("Server is reachable").
							SetJSONContent(
								ogen.NewSchema().
									SetType("object").
									AddRequiredProperties(
										ogen.String().ToProperty("status"),
									),
							),
					).
					AddResponse("503", ogen.NewResponse().SetDescription("Server is not reachable")),
				),
			)
			spec.AddPathItem(
				"/app-configs",
				ogen.NewPathItem().
					SetGet(ogen.NewOperation().
						SetOperationID("listAppConfig").
						AddParameters(paginationParameters()...).
						AddParameters(appEnvStackQueryParameters()...).
						AddResponse(
							"200",
							ogen.
								NewResponse().
								SetDescription("result AppConfig list").
								SetJSONContent(ogen.
									NewSchema().
									SetRef("#/components/schemas/AppConfigList").
									AsArray(),
								),
						).
						AddResponse("400", ogen.NewResponse().SetRef("#/components/responses/400")).
						AddResponse("404", ogen.NewResponse().SetRef("#/components/responses/404")).
						AddResponse("409", ogen.NewResponse().SetRef("#/components/responses/409")).
						AddResponse("500", ogen.NewResponse().SetRef("#/components/responses/500")),
					),
			)
			return nil
		}),
	)
	if err != nil {
		log.Fatalf("creating entoas extension: %v", err)
	}
	ogent, err := ogent.NewExtension(spec)
	if err != nil {
		log.Fatalf("creating ogent extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureUpsert,
		},
	}, entc.Extensions(ogent, oas))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
