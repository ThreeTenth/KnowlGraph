package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Dict holds the schema definition for the Dict entity.
type Dict struct {
	ent.Schema
}

// Fields of the Dict.
func (Dict) Fields() []ent.Field {
	return []ent.Field{
		field.Int("weight").Default(1),
	}
}

// Edges of the Dict.
func (Dict) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("dict").
			Unique().
			Required().
			StructTag(`json:"user,omitempty"`),
		edge.To("word", Word.Type).
			Unique().
			Required().
			StructTag(`json:"word,omitempty"`),
	}
}

// Indexes of the Dict.
func (Dict) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user").
			Edges("word").
			Unique(),
	}
}
