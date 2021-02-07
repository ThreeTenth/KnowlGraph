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
		field.Int("id"),
		field.Enum("status").
			Values("private", "public").
			Default("public").
			Comment("public: all users can access; private: only the source user can access."),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Version.Type),
		edge.To("reactions", Reaction.Type),
		edge.To("quote", Quote.Type).Unique(),
		edge.To("assets", Asset.Type),
	}
}
