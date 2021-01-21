package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Story holds the schema definition for the Story entity.
type Story struct {
	ent.Schema
}

// Fields of the Story.
func (Story) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("state").Values("private", "protected", "public").Default("private"),
	}
}

// Edges of the Story.
func (Story) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Content.Type),
		edge.To("fork", Story.Type).Unique(),
		edge.To("reactions", Reaction.Type),
		edge.To("quote", Quote.Type).Unique(),
		edge.To("stars", Star.Type),
	}
}
