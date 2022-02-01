package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vote holds the schema definition for the Vote entity.
type Vote struct {
	ent.Schema
}

// Fields of the Vote.
func (Vote) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").Values("allowed", "overruled", "abstained"),
		field.Strings("code").Optional(),   // 驳回原因代码
		field.String("comment").Optional(), // 驳回原因详细说明
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Vote.
func (Vote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("ras", RAS.Type).
			Unique().
			Required().
			StructTag(`json:"ras,omitempty"`),
	}
}
