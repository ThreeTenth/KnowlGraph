package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Reaction holds the schema definition for the Reaction entity.
type Reaction struct {
	ent.Schema
}

// Fields of the Reaction.
func (Reaction) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values("up", "down", "laugh", "hooray", "confused", "heart", "rocket", "eyes").
			Comment("Refer to the Reaction list on GitHub"),
		field.Int("count"),
	}
}

// Edges of the Reaction.
func (Reaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("article", Article.Type).
			Ref("reactions").
			Unique().
			StructTag(`json:"article,omitempty"`),
	}
}
