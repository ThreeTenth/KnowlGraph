package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
)

// Language holds the schema definition for the Language entity.
// see https://www.loc.gov/standards/iso639-2/php/code_list.php
type Language struct {
	ent.Schema
}

// Fields of the Language.
func (Language) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			MaxLen(3).
			NotEmpty().
			Unique().
			Immutable().Comment("Codes for the Representation of Names of Languages"),
		field.String("name"),
		field.Enum("direction").Values("ltr", "rtl").Default("ltr").Comment("Text writing direction"),
	}
}

// Edges of the Language.
func (Language) Edges() []ent.Edge {
	return nil
}
