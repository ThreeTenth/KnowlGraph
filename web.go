package main

import (
	"github.com/gin-gonic/gin"
)

// Header is html header data
type Header struct {
	GithubClientID string
	StaticDomain   string
	RestfulDomain  string
	UserLang       string
	Debug          bool
	VersionName    string
}

func html(fn func(p *Context) (int, string, interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{c}
		status, name, body := fn(ctx)

		_lang := ctx.UserLang()

		header := header()
		header.UserLang = _lang

		data := struct {
			Lang   string
			Header Header
			Body   interface{}
		}{
			Lang:   _lang,
			Header: header,
			Body:   body,
		}

		c.HTML(status, name, data)
		c.Abort()
	}
}

func header() Header {
	var header Header
	header.GithubClientID = config.Gci
	header.StaticDomain = config.Ssd
	header.RestfulDomain = config.Rad
	header.Debug = config.Debug
	header.VersionName = VersionName
	return header
}
