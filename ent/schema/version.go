package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Version holds the schema definition for the Version entity.
type Version struct {
	ent.Schema
}

// Fields of the Version.
func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default(""),
		field.String("comment").Default(""),
		field.String("title").Default(""),
		field.String("gist"),
		field.String("cover").Default(""),
		field.String("lang").Default("en"),
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
		edge.To("content", Content.Type).
			Unique().
			Required().
			StructTag(`json:"content,omitempty"`),
		edge.To("keywords", Word.Type).
			StructTag(`json:"keywords,omitempty"`),
		edge.To("quotes", Quote.Type).
			StructTag(`json:"quotes,omitempty"`),
		edge.From("article", Article.Type).
			Ref("versions").
			Unique().
			Required().
			StructTag(`json:"article,omitempty"`),
		edge.From("assets", Asset.Type).
			Ref("version"),
	}
}
