package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

func getStaticServerFiles(c *gin.Context) {
	_paths := c.Param("paths")[1:]

	var _type string
	if idx := strings.Index(_paths, "/"); -1 < idx {
		_type = _paths[:idx]
		_paths = _paths[idx+1:]
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var output []byte
	for idx := 0; -1 < idx; {
		var _name string

		idx = strings.Index(_paths, ",")
		if -1 == idx {
			_name = _paths
		} else if 0 < idx {
			_name = _paths[:idx]
			_paths = _paths[idx+1:]
		} else {
			continue
		}

		_bs, err := packr.NewBox("./res").Find(_type + "/" + _name)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		switch _type {
		case "js":
			output = append(output, ([]byte)(fmt.Sprintf("\r\n////////////////// %v/%v //////////////////\r\n", _type, _name))...)
		case "css":
			output = append(output, ([]byte)(fmt.Sprintf("\r\n/***************** %v/%v *****************/\r\n", _type, _name))...)
		}

		output = append(output, _bs...)
	}

	if 0 == len(output) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(_paths))
	if "" == ctype {
		ctype = "application/octet-stream"
	}

	c.Status(http.StatusOK)
	c.Writer.Header().Set("Content-Type", ctype)
	io.WriteString(c.Writer, (string)(output))
}
