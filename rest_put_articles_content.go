package main

import (
	"regexp"
	"strings"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"

	stripmd "github.com/writeas/go-strip-markdown"
)

// putDraftContent creates a content version for the specified article
func putDraftContent(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	var _content struct {
		Title       string `json:"title"`
		Body        string `json:"body" binding:"required"`
		VersionName string `json:"version_name"`
		DraftID     int    `json:"draftID" binding:"required,gt=0"`
	}
	err := c.ShouldBindJSON(&_content)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_version, err := client.User.
		Query().
		Where(user.IDEQ(_userID.(int))).
		QueryDrafts().
		Where(version.IDEQ(_content.DraftID)).
		First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_contentCreate := tx.Content.Create().
			SetBody(_content.Body).
			SetVersion(_version)

		if _content.Title != "" {
			_contentCreate.SetTitle(_content.Title)
		}

		if _content.VersionName != "" {
			_contentCreate.SetVersionName(_content.VersionName)
		}

		_content, err := _contentCreate.Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.Ok(_content.ID)
	})
}

func seo(title string, gist string, content string) string {
	if len(gist) != 0 {
		return gist
	}

	// Golang string is saved in UTF-8 format.
	// For languages that use multiple bytes to express one character,
	// such as Chinese characters, when substring the string,
	// invalid character encoding will appear. We need to convert string to []rune type.
	// Rune is stored in unicode encoding, it supports characters with long byte encoding.
	_rune := []rune(stripmd.Strip(content))

	var _seo string
	if len(_rune) <= MaxSeoLen {
		_seo = _seo + " " + string(_rune)
	} else {
		_seo = _seo + " " + string(_rune[:MaxSeoLen]) + "..."
	}

	re := regexp.MustCompile(`\r?\n`)
	_seo = re.ReplaceAllString(_seo, " ")
	return strings.Trim(_seo, " ")
}
