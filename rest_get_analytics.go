package main

import (
	"time"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/analytics"
	"knowlgraph.com/ent/predicate"
)

// AnalyticsResult is analytics result
type AnalyticsResult struct {
	ent.Analytics
	Count int
}

func getAnalytics(c *Context) error {
	var form struct {
		Event      string     `form:"event" binding:"required"`
		Category   []string   `form:"category"`
		Label      []string   `form:"label"`
		City       []string   `form:"city"`
		Timezone   []string   `form:"timezone"`
		Referrer   []string   `form:"referrer"`
		Lang       []string   `form:"lang"`
		Device     []string   `form:"device"`
		DeviceType []string   `form:"device_type"`
		Display    []string   `form:"display"`
		Os         []string   `form:"os"`
		OsVer      []string   `form:"os_version"`
		Platform   []string   `form:"platform"`
		PlatVer    []string   `form:"platform_version"`
		Browser    []string   `form:"browser"`
		BrowsVer   []string   `form:"browser_version"`
		StartTime  *time.Time `form:"start_time"`
		GroupBy    string     `form:"group_by"`
	}
	if err := c.ShouldBindQuery(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	predicates := make([]predicate.Analytics, 0, 2)
	predicates = append(predicates, analytics.EventEQ(form.Event))
	if form.Category != nil {
		if len(form.Category) == 1 {
			predicates = append(predicates, analytics.CategoryEQ(form.Category[0]))
		} else {
			predicates = append(predicates, analytics.CategoryIn(form.Category...))
		}
	}
	if form.Label != nil {
		if len(form.Label) == 1 {
			predicates = append(predicates, analytics.LabelEQ(form.Label[0]))
		} else {
			predicates = append(predicates, analytics.LabelIn(form.Label...))
		}
	}
	if form.City != nil {
		if len(form.City) == 1 {
			predicates = append(predicates, analytics.CityEQ(form.City[0]))
		} else {
			predicates = append(predicates, analytics.CityIn(form.City...))
		}
	}
	if form.Timezone != nil {
		if len(form.Timezone) == 1 {
			predicates = append(predicates, analytics.TimezoneEQ(form.Timezone[0]))
		} else {
			predicates = append(predicates, analytics.TimezoneIn(form.Timezone...))
		}
	}
	if form.Referrer != nil {
		if len(form.Referrer) == 1 {
			predicates = append(predicates, analytics.ReferrerEQ(form.Referrer[0]))
		} else {
			predicates = append(predicates, analytics.ReferrerIn(form.Referrer...))
		}
	}
	if form.Lang != nil {
		if len(form.Lang) == 1 {
			predicates = append(predicates, analytics.LangEQ(form.Lang[0]))
		} else {
			predicates = append(predicates, analytics.LangIn(form.Lang...))
		}
	}
	if form.Device != nil {
		if len(form.Device) == 1 {
			predicates = append(predicates, analytics.DeviceEQ(form.Device[0]))
		} else {
			predicates = append(predicates, analytics.DeviceIn(form.Device...))
		}
	}
	if form.DeviceType != nil {
		if len(form.DeviceType) == 1 {
			predicates = append(predicates, analytics.DeviceTypeEQ(form.DeviceType[0]))
		} else {
			predicates = append(predicates, analytics.DeviceTypeIn(form.DeviceType...))
		}
	}
	if form.Display != nil {
		if len(form.Display) == 1 {
			predicates = append(predicates, analytics.DisplayEQ(form.Display[0]))
		} else {
			predicates = append(predicates, analytics.DisplayIn(form.Display...))
		}
	}
	if form.Os != nil {
		if len(form.Os) == 1 {
			predicates = append(predicates, analytics.OsEQ(form.Os[0]))
			if form.OsVer != nil { // 只可以选择查看一个系统的某几个版本，以下同理
				predicates = append(predicates, analytics.OsVersionIn(form.OsVer...))
			}
		} else {
			predicates = append(predicates, analytics.OsIn(form.Os...))
		}
	}
	if form.Platform != nil {
		if len(form.Platform) == 1 {
			predicates = append(predicates, analytics.PlatformEQ(form.Platform[0]))
			if form.PlatVer != nil {
				predicates = append(predicates, analytics.PlatformVersionIn(form.PlatVer...))
			}
		} else {
			predicates = append(predicates, analytics.PlatformIn(form.Platform...))
		}
	}
	if form.Browser != nil {
		if len(form.Browser) == 1 {
			predicates = append(predicates, analytics.BrowserEQ(form.Browser[0]))
			if form.BrowsVer != nil && len(form.Browser) == 1 {
				predicates = append(predicates, analytics.BrowserVersionIn(form.BrowsVer...))
			}
		} else {
			predicates = append(predicates, analytics.BrowserIn(form.Browser...))
		}
	}
	if form.StartTime != nil {
		predicates = append(predicates, analytics.StartTimeGTE(*form.StartTime))
	} else {
		// 默认获取最近一天的统计数据
		predicates = append(predicates, analytics.StartTimeGTE(time.Now().AddDate(0, 0, -1)))
	}

	result := make([]*AnalyticsResult, 0)
	query := client.Analytics.Query().Where(analytics.And(predicates...))
	if form.GroupBy != "" && form.GroupBy != "start_time" {
		err := query.GroupBy(form.GroupBy).Aggregate(ent.Count()).Scan(ctx, &result)
		if err != nil {
			return c.InternalServerError(err.Error())
		}
		return c.Ok(&result)
	}
	count, err := query.Count(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	result = append(result, &AnalyticsResult{Analytics: ent.Analytics{Event: form.Event}, Count: count})

	return c.Ok(&result)
}
