package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
)

// UserWord holds the schema definition for the UserWord entity.
type UserWord struct {
	ent.Schema
}

// Fields of the UserWord.
func (UserWord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("weight").Default(1),
	}
}

// Edges of the UserWord.
func (UserWord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("words").
			Unique().
			Required().
			StructTag(`json:"user,omitempty"`),
		edge.To("word", Word.Type).
			Unique().
			Required().
			StructTag(`json:"word,omitempty"`),
	}
}

// Indexes of the UserWord.
func (UserWord) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user").
			Edges("word").
			Unique(),
	}
}
