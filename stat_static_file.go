package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const jsdelivrCDN = "https://cdn.jsdelivr.net/"
const jsdelivrFilepath = "gh/ThreeTenth/KnowlGraph@%v/build/%v"

// FileStat is static file info
type FileStat struct {
	Content []byte
	Type    string
}

func getStaticFile(c *gin.Context) {
	_filepath := c.Param("filepath")[1:]

	_fs, err := getSF4Redis(_filepath)

	if err != nil {
		_fs, err = getJsdeliverFile(_filepath)

		if err == nil {
			go saveSF2Redis(_filepath, _fs)
		}
	}

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.Status(http.StatusOK)
		c.Writer.Header().Set("Content-Type", _fs.Type)
		c.Writer.Header().Set("Cache-Control", "max-age=2592000")
		c.Writer.Write(_fs.Content)
	}
}

func getJsdeliverFile(_filepath string) (*FileStat, error) {
	_filenames := strings.Split(_filepath, ",")

	_jsdeliverFilepath := make([]string, len(_filenames))
	for i, _fn := range _filenames {
		_jsdeliverFilepath[i] = fmt.Sprintf(jsdelivrFilepath, VersionName, strings.Trim(_fn, " "))
	}

	_filepath = strings.Join(_jsdeliverFilepath, ",")
	if 1 < len(_filenames) {
		_filepath = "combine/" + _filepath
	}
	_url := jsdelivrCDN + _filepath

	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, err
	}

	_httpClient := &http.Client{}
	resp, err := _httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf(string(body))
	}

	_fileStat := FileStat{
		Content: body,
		Type:    resp.Header.Get("Content-Type"),
	}

	return &_fileStat, nil
}

func saveSF2Redis(filename string, _fs *FileStat) error {
	return rdb.Set(ctx, SF+filename, string(_fs.Content), ExpireTimeToken).Err()
}

func getSF4Redis(filename string) (*FileStat, error) {
	_cmd := rdb.Get(ctx, SF+filename)
	if err := _cmd.Err(); err != nil {
		return nil, err
	}

	_text := _cmd.Val()
	_type := mime.TypeByExtension(filepath.Ext(filename))
	if "" == _type {
		_type = "application/octet-stream"
	}

	return &FileStat{
		Content: []byte(_text),
		Type:    _type,
	}, nil
}
