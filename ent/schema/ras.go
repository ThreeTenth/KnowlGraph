package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
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
			Comment("The content is being voted, if allowed, it will be redirected to the content where published to the public space"),
		edge.To("voters", User.Type).Required().
			Comment("Voters are randomly selected and they will vote anonymously."),
	}
}
