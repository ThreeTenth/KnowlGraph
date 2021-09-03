package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

func getFavicon(c *gin.Context) {
	_box := packr.NewBox("./res/favicon")
	_favicon, _ := _box.Find("favicon.png")

	c.Status(http.StatusOK)
	header := c.Writer.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"image/png"}
	}

	if i, err := c.Writer.Write(_favicon); err != nil {
		fmt.Print(i, err.Error())
	}
	c.Abort()
}
