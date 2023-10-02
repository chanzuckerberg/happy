package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type AppConfig struct {
	ent.Schema
}

func (AppConfig) Fields() []ent.Field {
	return []ent.Field{
		field.
			Uint("id").
			SchemaType(map[string]string{"postgres": "bigserial"}).
			Annotations(
				entproto.Field(1),
			),
		field.
			Time("created_at").
			Optional().
			Annotations(
				entproto.Field(2),
			),
		field.
			Time("updated_at").
			Optional().
			Annotations(
				entproto.Field(3),
			),
		field.
			Time("deleted_at").
			Optional().
			Annotations(
				entproto.Field(4),
			),
		field.
			String("app_name").
			Optional().
			Annotations(
				entproto.Field(5),
			),
		field.
			String("environment").
			Optional().
			Annotations(
				entproto.Field(6),
			),
		field.
			String("stack").
			Optional().
			Annotations(
				entproto.Field(7),
			),
		field.
			String("key").
			Optional().
			Annotations(
				entproto.Field(8),
			),
		field.
			String("value").
			Optional().
			Annotations(
				entproto.Field(9),
			),
	}
}

func (AppConfig) Edges() []ent.Edge {
	return nil
}

func (AppConfig) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(
			entproto.PackageName("hapi"),
		),
	}
}
