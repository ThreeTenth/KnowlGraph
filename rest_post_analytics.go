package main

import (
	"net/url"
	"time"

	ua "github.com/mileusna/useragent"
)

func postAnalytics(c *Context) error {
	var body struct {
		Event      string    `json:"event" binding:"required"`
		Category   string    `json:"category" binding:"required"`
		Label      string    `json:"label" binding:"required"`
		Message    string    `json:"message"`
		City       string    `json:"city"`
		Timezone   string    `json:"timezone"`
		Referrer   string    `json:"referrer" binding:"required"`
		URL        string    `json:"url"`
		Path       string    `json:"path"`
		Lang       string    `json:"lang"`
		Device     string    `json:"device"`
		DeviceType string    `json:"device_type"`
		Display    string    `json:"display"`
		Os         string    `json:"os"`
		OsVer      string    `json:"os_version"`
		Platform   string    `json:"platform" binding:"required"`
		PlatVer    string    `json:"platform_version" binding:"required"`
		Browser    string    `json:"browser"`
		BrowsVer   string    `json:"browser_version"`
		StartTime  time.Time `json:"start_time" binding:"required"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	ua := ua.Parse(c.Request.UserAgent())

	if body.URL == "" {
		body.URL = c.Request.Referer()
	}
	if body.Path == "" {
		if _url, err := url.Parse(body.URL); err == nil {
			body.Path = _url.Path
		}
	}
	if body.Browser == "" {
		body.Browser = ua.Name
		body.BrowsVer = ua.Version
	}
	if body.Lang == "" {
		body.Lang = c.UserLang()
	}
	if body.Os == "" {
		body.Os = ua.OS
		body.OsVer = ua.OSVersion
	}
	if body.Device == "" {
		body.Device = ua.Device
	}
	if body.DeviceType == "" {
		if ua.Mobile {
			body.DeviceType = "Mobile"
		} else if ua.Tablet {
			body.DeviceType = "Tablet"
		} else if ua.Desktop {
			body.DeviceType = "Desktop"
		} else {
			body.DeviceType = "Unknown"
		}
	}
	if body.Display == "" {
		body.Display = "Unknown"
	}
	if body.City == "" {
		ipinfo, _ := ForeignIP(c.ClientIP())
		if ipinfo != nil {
			body.City = ipinfo.City
			body.Timezone = ipinfo.Timezone
		}
	}

	analyticsCreater := client.Analytics.Create().
		SetEvent(body.Event).
		SetCategory(body.Category).
		SetLabel(body.Label).
		SetMessage(body.Message).
		SetReferrer(body.Referrer).
		SetURL(body.URL).
		SetPath(body.Path).
		SetLang(body.Lang).
		SetCity(body.City).
		SetTimezone(body.Timezone).
		SetOs(body.Os).
		SetOsVersion(body.OsVer).
		SetPlatform(body.Platform).
		SetPlatformVersion(body.PlatVer).
		SetBrowser(body.Browser).
		SetBrowserVersion(body.BrowsVer).
		SetDevice(body.Device).
		SetDeviceType(body.DeviceType).
		SetDisplay(body.Display).
		SetBot(ua.Bot).
		SetStartTime(body.StartTime).
		SetVersion(1)

	if userID, ok := c.Get(GinKeyUserID); ok {
		analyticsCreater.SetUserID(userID.(int))
	}

	_, err = analyticsCreater.Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}
	return c.NoContent()
}
