package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Language holds the schema definition for the Language entity.
// see https://www.loc.gov/standards/iso639-2/php/code_list.php
type Language struct {
	ent.Schema
}

// Fields of the Language.
func (Language) Fields() []ent.Field {
	return []ent.Field{
		field.String("code"),
	}
}

// Edges of the Language.
func (Language) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("languages").
			Unique().
			StructTag(`json:"user,omitempty"`),
	}
}
