package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getFavicon(c *gin.Context) {
	_favicon, err := resource.ReadFile("res/favicon/favicon.png")
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Status(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(_favicon)))
	c.Writer.Write(_favicon)
	c.Abort()
}
