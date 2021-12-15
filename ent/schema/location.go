package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Location holds the schema definition for the Location entity.
// The location information must not contain country and province information.
type Location struct {
	ent.Schema
}

// Fields of the Location.
func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.String("city").
			Optional(),
		field.String("timezone"), // 时区
		field.Float("lat").
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(8,6)", // Override MySQL.
				dialect.Postgres: "numeric",      // Override Postgres.
			}),
		field.Float("long").
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "decimal(9,6)", // Override MySQL.
				dialect.Postgres: "numeric",      // Override Postgres.
			}),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Location.
func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("registrant", Terminal.Type).
			Ref("registration").
			Unique(),
		edge.From("terminal", Terminal.Type).
			Ref("locations").
			Unique(),
	}
}
