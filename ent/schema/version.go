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
		field.Enum("state").
			Values("review", "release", "reject").
			Default("review").
			Comment("release: all users can access.").
			Comment("review: only voters can access.").
			Comment("reject: nobody can access."),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("content", Content.Type).Unique().Required(),
		edge.To("lang", Language.Type).Unique().Required(),
		edge.To("tags", Tag.Type),
		edge.From("article", Article.Type).
			Ref("versions").
			Unique().
			Required(),
	}
}
