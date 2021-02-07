package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Tag holds the schema definition for the Tag entity.
type Tag struct {
	ent.Schema
}

// Fields of the Tag.
func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Enum("status").
			Values("private", "public").
			Default("private").
			Comment("public: tags associated with public articles; private: tags associated with private articles and not associated with public articles."),
	}
}

// Edges of the Tag.
func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Version.Type),
		edge.To("nodes", Node.Type),
	}
}
