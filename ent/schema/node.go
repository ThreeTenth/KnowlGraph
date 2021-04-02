package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Node holds the schema definition for the Node entity.
type Node struct {
	ent.Schema
}

// Fields of the Node.
func (Node) Fields() []ent.Field {
	return []ent.Field{
		field.Int("level"),
		field.Enum("status").
			Values("private", "public").
			Default("private").
			Comment("public: all words on the node path are public; private: there is at least one private word on the path."),
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
		edge.To("subnodes", Node.Type).
			Comment("subnodes: the next subdivision node of this node").
			From("belong").
			Comment("belong: belong to which node").
			Unique().
			StructTag(`json:"subnodes,omitempty"`),
		edge.To("branches", Node.Type).
			Comment("branches: root node attribute, which means all branch nodes under the root node.").
			From("root").
			Comment("root: the root node of the current node.").
			Unique().
			StructTag(`json:"branches,omitempty"`),
		edge.To("archives", Archive.Type).
			StructTag(`json:"archives,omitempty"`),
		edge.To("articles", Article.Type).
			StructTag(`json:"articles,omitempty"`),
	}
}
