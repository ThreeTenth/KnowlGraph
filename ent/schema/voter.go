package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Voter holds the schema definition for the Voter entity.
type Voter struct {
	ent.Schema
}

// Fields of the Voter.
func (Voter) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("voted").Default(false),
	}
}

// Edges of the Voter.
func (Voter) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique(),
		edge.From("ras", RAS.Type).Ref("voters").Unique(),
	}
}
