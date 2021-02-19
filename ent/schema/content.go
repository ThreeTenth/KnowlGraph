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
		field.String("body").Default(""),
		field.String("versionName").Optional().Unique(),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Content.
func (Content) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("version", Version.Type).Ref("history").Unique().Required(),
		edge.To("sections", Section.Type),
		edge.To("lang", Language.Type),
	}
}
