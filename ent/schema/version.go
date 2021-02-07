package schema

import (
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
		field.Enum("status").
			Values("editing", "voting", "published").
			Default("editing").
			Comment("editing: The article version is in the editing state. Other users can't access and except the editor.").
			Comment("voting: The article version is in voting status. Except for editors and voters, no other users can access. Editors can't edit, but they can withdraw publish.").
			Comment("published: The article version is in the published state. If the article is public, everyone can access it, and if the article is private, only the source user can access it."),
		field.String("seo").Optional(),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("content", Content.Type).Unique(),
		edge.From("tags", Tag.Type).Ref("versions"),
		edge.From("article", Article.Type).Ref("versions").Unique().Required(),
		edge.To("lang", Language.Type).Unique().Required(),
	}
}
