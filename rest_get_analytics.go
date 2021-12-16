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
		Event      string      `form:"event" binding:"required"`
		Category   []string    `form:"category"`
		Label      []string    `form:"label"`
		City       []string    `form:"city"`
		Timezone   []string    `form:"timezone"`
		Referrer   []string    `form:"referrer"`
		Lang       []string    `form:"lang"`
		Device     []string    `form:"device"`
		DeviceType []string    `form:"device_type"`
		Display    []string    `form:"display"`
		Os         []string    `form:"os"`
		OsVer      []string    `form:"os_version"`
		Platform   []string    `form:"platform"`
		PlatVer    []string    `form:"platform_version"`
		Browser    []string    `form:"browser"`
		BrowsVer   []string    `form:"browser_version"`
		Bot        bool        `form:"bot"`
		TimeStart  []time.Time `form:"time_start"`
		TimeUnit   int         `form:"time_unit"`
		GroupBy    []string    `form:"group_by"`
	}
	if err := c.ShouldBindQuery(&form); err != nil {
		return c.BadRequest(err.Error())
	}
	if form.GroupBy != nil {
		for _, v := range form.GroupBy {
			if v == "time_start" {
				return c.BadRequest(`The value of group_by cannot be "time_start"`)
			}
		}
	}

	predicates := make([]predicate.Analytics, 0, 3)
	predicates = append(predicates, analytics.BotEQ(form.Bot))
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
	if form.TimeStart == nil {
		form.TimeStart = make([]time.Time, 1)
	}
	if len(form.TimeStart) == 0 {
		// 默认获取最近一天的统计数据
		_1day := time.Now().AddDate(0, 0, -1)
		form.TimeStart = append(form.TimeStart, _1day)
	}

	resultList := make(map[time.Time][]*AnalyticsResult)
	for _, node := range form.TimeStart {
		startTime := node
		timePredicates := analytics.StartTimeGTE(startTime)
		if 0 < form.TimeUnit {
			// 如果有时间单位（TimeUnit），则仅返回该时间单位内的统计数据，以分钟（Minute）计
			endTime := node.Add(time.Duration(form.TimeUnit) * time.Minute)
			timePredicates = analytics.And(timePredicates, analytics.StartTimeLTE(endTime))
		}
		nodePredicates := append(predicates, timePredicates)
		result, err := getAnalyticsByPredicates(nodePredicates, form.GroupBy)
		if err != nil {
			ar := &AnalyticsResult{Analytics: ent.Analytics{Message: err.Error()}, Count: 1}
			resultList[node] = []*AnalyticsResult{ar}
		} else {
			resultList[node] = result
		}
	}

	return c.Ok(&resultList)
}

func getAnalyticsByPredicates(predicates []predicate.Analytics, groupBy []string) ([]*AnalyticsResult, error) {
	result := make([]*AnalyticsResult, 0)
	query := client.Analytics.Query().Where(analytics.And(predicates...))
	if groupBy != nil {
		analyticsGroupBy := query.GroupBy(groupBy[0])
		if len(groupBy) > 1 {
			analyticsGroupBy = query.GroupBy(groupBy[0], groupBy[1:]...)
		}
		err := analyticsGroupBy.Aggregate(ent.Count()).Scan(ctx, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	count, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	result = append(result, &AnalyticsResult{Count: count})

	return result, nil
}
