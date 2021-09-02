package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// RAS holds the schema definition for the RAS(Random anonymous space) entity.
type RAS struct {
	ent.Schema
}

// Fields of the RAS.
func (RAS) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("allowed", "rejected", "abstained").Optional().
			Comment("Indicates the status of the voting result of the current space"),
		field.String("comment").Optional().
			Comment("Postscript when creating the space"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the RAS.
func (RAS) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("version", Version.Type).Unique().Required().
			Comment("The content is being voted, if allowed, it will be redirected to the content where published to the public space").
			StructTag(`json:"version,omitempty"`),
		edge.To("voters", Voter.Type).
			Comment("Voters are randomly selected and they will vote anonymously.").
			StructTag(`json:"voters,omitempty"`),
	}
}
