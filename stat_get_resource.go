package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ResourceReadDir return `[]fs.DirEntry` from `src`` dir.
func ResourceReadDir(src string) ([]fs.DirEntry, error) {
	if config.Debug {
		return os.ReadDir(filepath.Join("./", src))
	}
	return resource.ReadDir(src)
}

// ResourceReadFile return `[]byte` from `src`` file.
func ResourceReadFile(src string) ([]byte, error) {
	if config.Debug {
		return os.ReadFile(filepath.Join("./", src))
	}
	return resource.ReadFile(src)
}

var mainCSSBuffer *bytes.Buffer

func getMainCSS(c *gin.Context) {
	if mainCSSBuffer == nil {
		mainCSSBuffer = new(bytes.Buffer)
	}
	if 0 == mainCSSBuffer.Len() || config.Debug {
		mainCSSBuffer.Reset()

		_, err := writeEmbedDir(ResourceReadDir, ResourceReadFile, "res/css", mainCSSBuffer, nil)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "text/css; charset=utf-8")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(mainCSSBuffer.Len()))
	c.Writer.Write(mainCSSBuffer.Bytes())
	c.Abort()
}

var appJSBuffer *bytes.Buffer

func getAppJS(c *gin.Context) {
	if appJSBuffer == nil {
		appJSBuffer = new(bytes.Buffer)
	}
	if 0 == appJSBuffer.Len() || config.Debug {
		appJSBuffer.Reset()
		languagesFormat := "// %v, %v \nconst languages = %v"
		defaultStringsFormat := "// %v, %v \nconst defaultLang = %v"

		_, err := writeEmbedFile(ResourceReadFile, "res/strings/languages.json", appJSBuffer, languagesFormat)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_, err = writeEmbedFile(ResourceReadFile, "res/strings/strings-en.json", appJSBuffer, defaultStringsFormat)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_, err = writeEmbedDir(ResourceReadDir, ResourceReadFile, "res/js", appJSBuffer, []string{"res/js/app.js"})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		_, err = writeEmbedFile(ResourceReadFile, "res/js/app.js", appJSBuffer)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(appJSBuffer.Len()))
	c.Writer.Write(appJSBuffer.Bytes())
	c.Abort()
}

var themeBuffer *bytes.Buffer

func getDefaultThemes(c *gin.Context) {
	if themeBuffer == nil {
		themeBuffer = new(bytes.Buffer)
	}
	if 0 == themeBuffer.Len() || config.Debug {
		themeBuffer.Reset()
		themeFormat := "// %v \nconst %v = `%v`"

		var err error
		_, err = writeEmbedDir(ResourceReadDir, ResourceReadFile, "res/theme", themeBuffer, nil, themeFormat)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.Status(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(themeBuffer.Len()))
	c.Writer.Write(themeBuffer.Bytes())
	c.Abort()
}

func writeEmbedDir(
	readDir func(src string) ([]fs.DirEntry, error),
	readFile func(src string) ([]byte, error),
	src string,
	dst io.Writer,
	ignores []string,
	format ...string,
) (int, error) {
	themes, err := readDir(src)
	if err != nil {
		return 0, err
	}

	count := 0
	dir := src + "/"
	for _, theme := range themes {
		length := 0
		path := dir + theme.Name()
		skip := false
		if nil != ignores {
			for _, ignore := range ignores {
				if path == ignore {
					skip = true
					break
				}
			}
		}
		if skip {
			continue
		}

		if theme.IsDir() {
			length, err = writeEmbedDir(readDir, readFile, path, dst, ignores, format...)
		} else {
			length, err = writeEmbedFile(readFile, path, dst, format...)
		}

		count += length
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func writeEmbedFile(
	readFile func(src string) ([]byte, error),
	src string,
	dst io.Writer,
	format ...string,
) (int, error) {
	bs, err := readFile(src)
	if err != nil {
		return 0, err
	}

	if 1 == len(format) {
		_, fullname := filepath.Split(src)
		var extension = filepath.Ext(fullname)
		var filename = fullname[0 : len(fullname)-len(extension)]
		bs = []byte(fmt.Sprintf(format[0], fullname, filename, string(bs)))
	}

	return dst.Write(bs)
}
