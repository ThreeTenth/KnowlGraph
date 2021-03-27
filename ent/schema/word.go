package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
		edge.To("nodes", Node.Type),
		edge.To("definition", Article.Type).Unique(),
		edge.From("users", User.Type).Ref("words"),
		edge.From("versions", Version.Type).Ref("keywords"),
	}
}
