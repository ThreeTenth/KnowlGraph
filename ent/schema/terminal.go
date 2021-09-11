package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Terminal holds the schema definition for the Terminal entity.
type Terminal struct {
	ent.Schema
}

// Fields of the Terminal.
func (Terminal) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("ua"),
		field.Time("created_at").
			Default(time.Now).
			StructTag(`json:"-"`),
		// field.String("cred_id"),
		// field.String("public_key"),
		// field.String("attestation_type"),
	}
}

// Edges of the Terminal.
func (Terminal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("terminals").
			Unique().
			StructTag(`json:"user,omitempty"`),
	}
}
