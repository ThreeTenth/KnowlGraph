package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Star holds the schema definition for the Star entity.
type Star struct {
	ent.Schema
}

// Fields of the Star.
func (Star) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("browse", "star", "watch").Default("browse"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Star.
func (Star) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("stars").Unique().Required(),
		edge.From("article", Article.Type).
			Ref("stars").
			Unique(),
		edge.From("node", Node.Type).
			Ref("stars").
			Unique(),
		edge.To("nodes", Star.Type).
			From("articles"),
	}
}
