package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("private", "public").Default("private"),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Content.Type),
		edge.To("main", Article.Type).Unique(),
		edge.To("reactions", Reaction.Type),
		edge.To("quote", Quote.Type).Unique(),
		edge.To("stars", Star.Type),
	}
}
