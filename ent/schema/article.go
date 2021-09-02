package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
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
		edge.To("versions", Version.Type).
			StructTag(`json:"versions,omitempty"`),
		edge.To("branches", Draft.Type).
			StructTag(`json:"branches,omitempty"`),
		edge.To("reactions", Reaction.Type).
			StructTag(`json:"reactions"`),
		edge.To("assets", Asset.Type).
			StructTag(`json:"assets"`),
		edge.To("quotes", Quote.Type).
			StructTag(`json:"quotes,omitempty"`),
		edge.From("nodes", Node.Type).
			Ref("articles").
			StructTag(`json:"nodes"`),
		edge.From("word", Word.Type).
			Ref("definition").
			Unique().StructTag(`json:"word,omitempty"`),
		edge.From("quote", Quote.Type).
			Ref("response").
			Unique().StructTag(`json:"quote,omitempty"`),
	}
}
