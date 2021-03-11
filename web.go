package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Header is html header data
type Header struct {
	GithubClientID string
	StaticDomain   string
	RestfulDomain  string
	UserLang       string
	Debug          bool
}

func html(fn func(p *Context) (int, string, interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, name, body := fn(&Context{c})

		_lang, err := c.Cookie(CookieUserLang)
		if err != nil {
			_lang = c.GetHeader(HeaderAcceptLanguage)
			idx := strings.Index(_lang, ",")
			if 0 < idx {
				_lang = _lang[:idx]
			}
			if idx = strings.Index(_lang, ";"); 0 < idx {
				_lang = _lang[:idx]
			}
			if idx = strings.Index(_lang, "-"); 0 < idx {
				_lang = _lang[:idx]
			}
		}
		if _lang == "" {
			_lang = DefaultLanguage
		}

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
	return header
}
