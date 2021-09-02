package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Quote holds the schema definition for the Quote entity.
type Quote struct {
	ent.Schema
}

// Fields of the Response.
func (Quote) Fields() []ent.Field {
	return []ent.Field{
		field.String("text").Optional(),
		field.String("context").Optional(),
		field.String("source"),
		field.Int("highlight").Default(0).Comment("Mark quoted content with colors, lines, etc."),
	}
}

// Edges of the Response.
func (Quote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("response", Article.Type).
			Unique().
			Required().
			StorageKey(edge.Column("response_id")).
			StructTag(`json:"response,omitempty"`),
		edge.From("article", Article.Type).
			Ref("quotes").
			Unique(),
		edge.From("version", Version.Type).
			Ref("quotes").
			Unique().
			StructTag(`json:"version,omitempty"`),
	}
}
