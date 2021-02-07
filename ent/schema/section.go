package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
			Comment("belong: belong to which pact").
			Unique(),
	}
}
