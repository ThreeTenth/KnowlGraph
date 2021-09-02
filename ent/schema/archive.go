package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Archive holds the schema definition for the Archive entity.
type Archive struct {
	ent.Schema
}

// Fields of the Archive.
func (Archive) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values("star", "watch").
			Default("star").
			Comment("star: a node favored by the user; watch: a node followed by the user."),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Archive.
func (Archive) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("archives").
			Unique().
			Required().
			Comment("The owner of the archive").
			StructTag(`json:"user,omitempty"`),
		edge.From("node", Node.Type).
			Ref("archives").
			Unique().
			Required().
			StructTag(`json:"node"`),
	}
}
