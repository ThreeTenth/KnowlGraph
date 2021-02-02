package schema

import (
	"time"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Content holds the schema definition for the Content entity.
type Content struct {
	ent.Schema
}

// Fields of the Content.
func (Content) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").Optional(),
		field.String("gist").Optional(),
		field.String("content").Default(""),
		field.Int("version").Default(1),
		field.String("versionName").Optional().Unique(),
		field.String("seo").Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Content.
func (Content) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", Tag.Type).Ref("contents"),
		edge.From("story", Story.Type).Ref("versions").Unique().Required(),
		edge.To("lang", Language.Type).Unique().Required(),
		edge.To("main", Content.Type).Unique().Comment("Main is a copy of the main branch of the current copy."),
	}
}
