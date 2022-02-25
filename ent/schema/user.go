package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
		field.Int("github_id").Optional().StructTag(`json:"-"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("assets", Asset.Type).
			StorageKey(edge.Column("owner_id")).
			StructTag(`json:"assets,omitempty"`),
		edge.To("archives", Archive.Type).
			StorageKey(edge.Column("owner_id")).
			StructTag(`json:"archives,omitempty"`),
		edge.To("drafts", Draft.Type).
			StorageKey(edge.Column("editor_id")).
			StructTag(`json:"drafts,omitempty"`),
		edge.To("dict", Dict.Type).
			StorageKey(edge.Column("user_dict")).
			StructTag(`json:"dict,omitempty"`),
		edge.To("languages", Language.Type).
			StructTag(`json:"languages,omitempty"`),
		edge.To("terminals", Terminal.Type).
			StructTag(`json:"terminals,omitempty"`),
	}
}
