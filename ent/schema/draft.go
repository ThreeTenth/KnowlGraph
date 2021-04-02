package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Draft holds the schema definition for the Draft entity.
type Draft struct {
	ent.Schema
}

// Fields of the Draft.
func (Draft) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("state").Values("read", "write").Default("write"),
	}
}

// Edges of the Draft.
func (Draft) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("snapshots", Content.Type).
			StructTag(`json:"snapshots,omitempty"`),
		edge.To("original", Version.Type).Unique().
			StructTag(`json:"original,omitempty"`),
		edge.From("user", User.Type).Ref("drafts").Unique().Required().
			StructTag(`json:"user,omitempty"`),
		edge.From("article", Article.Type).Ref("branches").Unique().Required().
			StructTag(`json:"article,omitempty"`),
	}
}
