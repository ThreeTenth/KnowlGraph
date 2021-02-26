package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Version holds the schema definition for the Version entity.
type Version struct {
	ent.Schema
}

// Fields of the Version.
func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Unique(),
		field.String("comment").Optional().Unique(),
		field.String("title").Optional(),
		field.String("seo").Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("content", Content.Type).Unique().Required(),
		edge.From("article", Article.Type).
			Ref("versions").
			Unique().
			Required(),
	}
}
