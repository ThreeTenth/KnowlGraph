package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Section holds the schema definition for the Section entity.
type Section struct {
	ent.Schema
}

// Fields of the Section.
func (Section) Fields() []ent.Field {
	return []ent.Field{
		field.Int("index"),
		field.String("text"),
		field.Enum("type").Values("title", "text", "gist", "case").Default("text"),
	}
}

// Edges of the Section.
func (Section) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subsections", Section.Type).
			Comment("subsections: the next subdivision sections of this section").
			From("belong").
			Unique().
			Comment("belong: belong to which pact").
			StructTag(`json:"subsections,omitempty"`),
		edge.From("content", Content.Type).
			Ref("sections").
			Unique().
			Required().
			StructTag(`json:"content,omitempty"`),
	}
}
