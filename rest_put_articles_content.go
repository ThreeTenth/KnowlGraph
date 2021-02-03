package main

import (
	"regexp"
	"strings"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"

	stripmd "github.com/writeas/go-strip-markdown"
)

// putArticleContent creates a content version for the specified article
func putArticleContent(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	var _body struct {
		Title       string   `json:"title"`
		Gist        string   `json:"gist"`
		Content     string   `json:"content" binding:"required"`
		VersionName string   `json:"versionName"`
		ArticleID   int      `json:"articleID" binding:"required,gt=0,articleExist"`
		Tags        []string `json:"tags" binding:"max=5" comment:"5 at most"`
	}
	err := c.ShouldBindJSON(&_body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if nil != _body.Tags && 5 < len(_body.Tags) {
		return c.BadRequest("Up to 5 tags")
	}

	_article, err := client.User.
		Query().
		Where(user.IDEQ(_userID.(int))).
		QueryArticles().
		Where(article.IDEQ(_body.ArticleID)).
		First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	_lang := c.Query(QueryLang)

	var columns []struct {
		Lang string `json:"content_lang"`
		Max  int    `json:"max"`
	}

	_article.QueryVersions().
		Where(content.HasLangWith(language.ID(_lang))).
		GroupBy(content.LangColumn).
		Aggregate(ent.Max(content.FieldVersion)).
		Scan(ctx, &columns)

	maxVersion := 1
	if 1 == len(columns) && maxVersion <= columns[0].Max {
		maxVersion = columns[0].Max
	}

	seo := seo(_body.Title, _body.Gist, _body.Content)

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_contentCreate := tx.Content.Create().
			SetTitle(_body.Title).
			SetGist(_body.Gist).
			SetContent(_body.Content).
			SetArticle(_article).
			SetVersion(maxVersion).
			SetSeo(seo).
			SetLangID(_lang)

		if nil != _body.Tags && 0 < len(_body.Tags) {
			_tags := make([]*ent.Tag, len(_body.Tags))
			for i, _tag := range _body.Tags {
				if "" == _tag {
					continue
				}
				_tags[i], err = tx.Tag.Query().Where(tag.NameContainsFold(_tag)).First(ctx)
				if err != nil {
					_tags[i], err = tx.Tag.Create().SetName(_tag).Save(ctx)
					if nil != err {
						return c.InternalServerError(err.Error())
					}
				}
			}
			_contentCreate.AddTags(_tags...)
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
