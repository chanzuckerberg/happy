package schema

import (
	"context"
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	gen "github.com/chanzuckerberg/happy/api/pkg/ent"
	"github.com/chanzuckerberg/happy/api/pkg/ent/appconfig"
	"github.com/chanzuckerberg/happy/api/pkg/ent/hook"
)

type AppConfig struct {
	ent.Schema
}

func (AppConfig) Fields() []ent.Field {
	return []ent.Field{
		field.
			Uint("id").
			SchemaType(map[string]string{"postgres": "bigserial"}),
		field.
			Time("created_at").
			Default(time.Now).
			Immutable(),
		field.
			Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.
			Time("deleted_at").
			Optional().
			Nillable().
			Default(nil),
		field.
			String("app_name"),
		field.
			String("environment"),
		field.
			String("stack").
			Default(""),
		field.
			String("key"),
		field.
			Text("value"),
		field.
			Enum("source").
			Values("stack", "environment").
			Default("environment").
			Comment("'stack' if the config is for a specific stack or 'environment' if available to all stacks in the environment"),
	}
}

func (AppConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Fields("deleted_at"),
		index.
			Fields("app_name", "environment", "stack", "key").
			Unique(),
	}
}

func (AppConfig) Edges() []ent.Edge {
	return nil
}

func (AppConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Make this readonly for now
		entoas.DeleteOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.CreateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.UpdateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
	}
}

func (AppConfig) Hooks() []ent.Hook {
	return []ent.Hook{
		// hook to populate the source field
		func(next ent.Mutator) ent.Mutator {
			return hook.AppConfigFunc(func(ctx context.Context, m *gen.AppConfigMutation) (ent.Value, error) {
				source := appconfig.SourceEnvironment
				if stack, ok := m.Stack(); ok && stack != "" {
					source = appconfig.SourceStack
				}
				m.SetSource(source)
				return next.Mutate(ctx, m)
			})
		},
	}
}
