package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

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

		fmt.Println(_lang)
		box := packr.NewBox("./res/strings")
		bs, err := box.Find("strings-" + _lang + ".json")
		if err != nil {
			bs, err = box.Find("strings-" + DefaultLanguage + ".json")
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_strings := make(map[string]interface{})
		err = json.Unmarshal(bs, &_strings)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		data := struct {
			Header  interface{}
			Body    interface{}
			Strings map[string]interface{}
		}{
			Header:  header(),
			Body:    body,
			Strings: _strings,
		}

		c.HTML(status, name, data)
		c.Abort()
	}
}

func header() interface{} {
	var header struct {
		GithubClientID string
		StaticDomain   string
		RestfulDomain  string
		Debug          bool
	}

	header.GithubClientID = config.Gci
	header.StaticDomain = config.Ssd
	header.RestfulDomain = config.Rad
	header.Debug = config.Debug
	return header
}
