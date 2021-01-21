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
		field.String("content"),
		field.String("lang").Default("EN").Comment("The language version of the content"),
		field.Int("version").Default(1).Comment("The content version number"),
		field.Enum("status").Values("auto", "save").Default("auto").Comment("auto means the version is automatically saved, save means manual save"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Content.
func (Content) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tags", Tag.Type).Ref("contents"),
		edge.From("story", Story.Type).Ref("versions").Unique(),
	}
}
