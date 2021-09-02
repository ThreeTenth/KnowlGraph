package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Word holds the schema definition for the Word entity.
type Word struct {
	ent.Schema
}

// Fields of the Word.
func (Word) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Enum("status").
			Values("private", "public").
			Default("private").
			Comment("public: words associated with public articles; private: words associated with private articles and not associated with public articles."),
	}
}

// Edges of the Word.
func (Word) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("nodes", Node.Type).
			StructTag(`json:"nodes,omitempty"`),
		edge.To("definition", Article.Type).
			Unique().
			StructTag(`json:"definition,omitempty"`),
		edge.From("versions", Version.Type).
			Ref("keywords").
			StructTag(`json:"versions,omitempty"`),
	}
}
