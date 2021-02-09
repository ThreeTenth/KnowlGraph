package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional(),
		field.String("email").Optional(),
		field.String("avatar").Optional(),
		field.Int("github_id").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assets", Asset.Type).StorageKey(edge.Column("owner_id")),
		edge.To("archives", Archive.Type).StorageKey(edge.Column("owner_id")),
		edge.To("drafts", Version.Type).StorageKey(edge.Column("editor_id")),
		edge.To("tags", Tag.Type),
		edge.To("languages", Language.Type),
		edge.From("rass", RAS.Type).Ref("voters"),
	}
}
