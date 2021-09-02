package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Node holds the schema definition for the Node entity.
type Node struct {
	ent.Schema
}

// Fields of the Node.
func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("public", "private").Immutable(),
	}
}

// Edges of the Node.
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("word", Word.Type).
			Ref("nodes").
			Unique().
			Required().
			StructTag(`json:"word,omitempty"`),
		edge.To("nexts", Node.Type).
			Comment("nexts: the next subdivision node of this node").
			StructTag(`json:"nexts,omitempty"`).
			From("prev").
			Comment("prev: belong to which node").
			Unique().
			StructTag(`json:"prev,omitempty"`),
		edge.To("path", Node.Type).
			StructTag(`json:"path,omitempty"`),
		edge.To("archives", Archive.Type).
			StructTag(`json:"archives,omitempty"`),
		edge.To("articles", Article.Type).
			StructTag(`json:"articles,omitempty"`),
	}
}
