package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Analytics holds the schema definition for the Analytics entity.
type Analytics struct {
	ent.Schema
}

// Fields of the Analytics.
func (Analytics) Fields() []ent.Field {
	return []ent.Field{
		// 事件名称
		field.String("event").Default(""),
		// 事件主体
		field.String("category").Default(""),
		field.String("label").Default(""),
		field.String("message").Default(""),
		// 地点、时区
		field.String("city").Default(""),
		field.String("timezone").Default(""),
		// 来源和路径
		field.String("referrer").Default(""),
		field.String("url").Default(""),
		field.String("path").Default(""),
		// 语言
		field.String("lang").Default(""),
		// 设备（型号和类型，PC 或 phone 或 pad）
		// 设备分辨率
		// 系统（操作系统，和版本号）
		// 平台（Web 或 移动应用程序或 PC 应用或其他，如手表、物联网设备等，以及他们的版本号）
		// 浏览器（Chrome 或 edge 或客户端，以及他们的版本号）
		field.String("device").Default(""),
		field.String("device_type").Default(""),
		field.String("display").Default(""),
		field.String("os").Default(""),
		field.String("os_version").Default(""),
		field.String("platform").Default(""),
		field.String("platform_version").Default(""),
		field.String("browser").Default(""),
		field.String("browser_version").Default(""),
		// 统计数据接口版本号
		field.Int("version"),
		// 发起统计时间
		field.Time("start_time").StructTag(`json:"-"`),
		field.Time("created_at").StructTag(`json:"-"`).Default(time.Now),
	}
}

// Edges of the Analytics.
func (Analytics) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			StructTag(`json:"user,omitempty"`),
		edge.To("terminal", Terminal.Type).
			Unique().
			StructTag(`json:"terminal,omitempty"`),
	}
}
