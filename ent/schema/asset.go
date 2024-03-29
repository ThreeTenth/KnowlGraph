package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Fields of the Asset.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values("browse", "star", "watch", "self").
			Default("browse").
			Comment("browse: Articles visited by users.").
			Comment("star: a article favored by the user.").
			Comment("watch: a article followed by the user.").
			Comment("self: a private article created by the user."),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Asset.
func (Asset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("assets").
			Unique().
			Required().
			Comment("The owner of the asset").
			StructTag(`json:"user,omitempty"`),
		edge.From("article", Article.Type).
			Ref("assets").
			Unique().
			Required().
			StructTag(`json:"article,omitempty"`),
		edge.To("version", Version.Type).
			Unique().
			StructTag(`json:"version,omitempty"`),
	}
}
