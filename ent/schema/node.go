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
		field.String("cover").Optional(),
		field.Int("level"),
	}
}

// Edges of the Node.
func (Node) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tag", Tag.Type).Ref("nodes").Unique().Required(),
		edge.To("fork", Node.Type).
			From("form").
			Unique(),
		edge.To("nodes", Node.Type).
			From("root").
			Unique(),
		edge.To("stars", Star.Type),
	}
}
