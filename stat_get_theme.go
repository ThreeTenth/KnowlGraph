package main

// import (
// 	"fmt"
// 	"net/http"
// 	"text/template"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gobuffalo/packr/v2"
// )

// func getStaticTheme(c *gin.Context) {
// 	_id := c.Param("id")
// 	var _box *packr.Box
// 	if "" == _id {
// 		_box = packr.NewBox("./res/theme")
// 	} else {
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	}

// 	_data := make(map[string]string)

// 	for _, v := range _box.List() {
// 		if TplThemeJS == v {
// 			continue
// 		}
// 		_html, _ := _box.FindString(v)
// 		_data[v[:len(v)-5]] = _html
// 	}

// 	_json, _ := packr.NewBox("./res").FindString("languages.json")

// 	_data["languages"] = _json
// 	_data["restfulDomain"] = config.Rad
// 	_data["staticDomain"] = config.Ssd

// 	_themeJS, _ := _box.FindString(TplThemeJS)

// 	c.Status(http.StatusOK)
// 	header := c.Writer.Header()
// 	if val := header[HeaderContentType]; len(val) == 0 {
// 		header[HeaderContentType] = []string{"text/javascript; charset=utf-8"}
// 	}

// 	tpl := template.Must(template.New("").Parse(_themeJS))
// 	if err := tpl.Execute(c.Writer, &_data); err != nil {
// 		fmt.Print(err.Error())
// 	}
// 	c.Abort()
// }
