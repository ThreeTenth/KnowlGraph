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
		field.String("body").Default(""),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Content.
func (Content) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("last", Content.Type).Unique().
			StructTag(`json:"last,omitempty"`),
		edge.To("sections", Section.Type).
			StructTag(`json:"sections,omitempty"`),
		edge.From("branche", Draft.Type).
			Ref("snapshots").
			Unique().
			Required().
			Comment("Snapshot of a branch.").
			StructTag(`json:"branche,omitempty"`),
		edge.From("version", Version.Type).
			Ref("content").
			Unique().
			Comment("The content of a published version. If the article is public, a vote is required. If the article is private, it is created directly.").
			StructTag(`json:"version,omitempty"`),
	}
}
