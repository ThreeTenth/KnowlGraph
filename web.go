package main

import (
	"github.com/gin-gonic/gin"
)

func html(fn func(p *Context) (int, string, interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, name, body := fn(&Context{c})

		data := struct {
			Header interface{}
			Body   interface{}
		}{
			Header: header(),
			Body:   body,
		}

		c.HTML(status, name, data)
		c.Abort()
	}
}

func header() interface{} {
	var header struct {
		GithubClientID string
	}

	return header
}
