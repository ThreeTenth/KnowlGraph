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
		field.String("code").
			StructTag(`json:"-"`),
		field.String("name"),
		field.String("ua"),
		field.Bool("only_once").
			Default(false).
			StructTag(`json:"onlyOnce"`), // 临时授权，仅此一次
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
	// 终端的地理位置信息和创建时间默认不向客户端显示，需要得到用户的授权，才能显示。
	return []ent.Edge{
		// 登记地
		edge.To("registration", Location.Type).
			Unique().
			StructTag(`json:"-"`),
		// 上线地
		edge.To("locations", Location.Type).
			StructTag(`json:"-"`),
		edge.To("credential", Credential.Type).
			StructTag(`json:"-"`).
			Unique(),

		edge.From("user", User.Type).
			Ref("terminals").
			Unique().
			StructTag(`json:"user,omitempty"`),
	}
}
