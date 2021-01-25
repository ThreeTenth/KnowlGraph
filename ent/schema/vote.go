package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Vote holds the schema definition for the Vote entity.
type Vote struct {
	ent.Schema
}

// Fields of the Vote.
func (Vote) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("allowed", "rejected", "abstained"),
		field.String("response").Optional().Comment("Response to voting content"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Vote.
func (Vote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ras", RAS.Type).Unique(),
	}
}
