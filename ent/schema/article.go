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
			Values("private", "protected", "public").
			Default("public").
			Comment("public: all users can access and edit.").
			Comment("protected: all users just can access, used to built-in protocol.").
			Comment("private: only the source user can access and edit."),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", Version.Type),
		edge.To("branches", Draft.Type),
		edge.To("reactions", Reaction.Type),
		edge.To("assets", Asset.Type),
		edge.To("quotes", Quote.Type),
		edge.From("nodes", Node.Type).
			Ref("articles"),
		edge.From("word", Word.Type).
			Ref("definition").
			Unique(),
		edge.From("quote", Quote.Type).
			Ref("response").
			Unique(),
	}
}
