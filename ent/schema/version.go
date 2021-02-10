package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
)

// Version holds the schema definition for the Version entity.
type Version struct {
	ent.Schema
}

// Fields of the Version.
func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values("new", "modify", "translate", "review", "release").
			Default("new").
			Comment("new: The first version of the new article. Other users can't access and except the editor.").
			Comment("modify: Modified version for a language version. Other users can't access and except the editor.").
			Comment("translate: Translated version for a language version. Other users can't access and except the editor.").
			Comment("review: After the public article is pushed, the version that receives the vote before it is published. Except for editors and voters, no other users can access. They can't edit, but editor can withdraw publish.").
			Comment("release: If the article is public, everyone can access and edit it. if the article is private, only the source user can access it."),
		field.String("title").Optional(),
		field.String("seo").Optional(),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("content", Content.Type).Unique(),
		edge.To("lang", Language.Type).Unique().Required(),
		edge.To("history", Content.Type),
		edge.To("main", Version.Type).Unique(),
		edge.From("tags", Tag.Type).Ref("versions"),
		edge.From("article", Article.Type).Ref("versions").Unique().Required(),
		edge.From("user", User.Type).Ref("drafts").Unique(),
	}
}

// Indexs of the Version.
func (Version) Indexs() []ent.Index {
	return []ent.Index{
		index.Fields("status").Edges("article").Edges("lang").Edges("user").Unique(),
	}
}
