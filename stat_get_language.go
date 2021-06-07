package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"text/template"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gobuffalo/packr/v2"
// )

// func getStaticLang(c *gin.Context) {
// 	_id := c.Param("id")

// 	var _box *packr.Box

// 	if "" == _id {
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	}

// 	_box = packr.NewBox("./res/strings")

// 	_langStrings, err := _box.FindString(strings.ReplaceAll(_id, ".js", ".json"))
// 	if err != nil {
// 		c.AbortWithError(http.StatusNotFound, err)
// 		return
// 	}

// 	_temp := `const defaultLang = {{ . }}`

// 	c.Status(http.StatusOK)
// 	header := c.Writer.Header()
// 	if val := header[HeaderContentType]; len(val) == 0 {
// 		header[HeaderContentType] = []string{"text/javascript; charset=utf-8"}
// 	}

// 	tpl := template.Must(template.New("").Parse(_temp))
// 	if err := tpl.Execute(c.Writer, &_langStrings); err != nil {
// 		fmt.Print(err.Error())
// 	}
// 	c.Abort()
// }
