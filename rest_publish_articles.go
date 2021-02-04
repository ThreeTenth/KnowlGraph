package main

import (
	"fmt"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
)

// publishArticle pushes a personal article content to a public channel
func publishArticle(c *Context) error {
	var _body struct {
		ArticleID int `json:"articleID" binding:"required,gt=0,articleExist"`
		ContentID int `json:"ContentID" binding:"required,gt=0"`
	}
	err := c.ShouldBindJSON(&_body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	///////////////// User authentication //////////////////
	//
	//
	////////////////////////////////////////////////////////

	_article, err := client.Article.Get(ctx, _body.ArticleID)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if _article.Status != article.StatusPrivate {
		return c.Forbidden("The operation was blocked because the content is not private")
	}

	//////////////// Article authentication //////////////////
	//
	//
	////////////////////////////////////////////////////////

	_content, err := _article.QueryVersions().Where(content.ID(_body.ContentID)).First(ctx)
	if err != nil {
		return c.Forbidden("The operation was blocked because the content does not belong to the article")
	}

	_published, _ := _article.QueryMain().First(ctx)

	_lang, _ := _content.QueryLang().FirstID(ctx)

	if _published != nil {
		ok, _ := _published.QueryVersions().
			Where(
				content.And(
					content.VersionGT(_content.Version),
					content.HasLangWith(language.ID(_lang)))).
			Exist(ctx)
		if ok {
			return c.Created("Please update the latest version before publishing, because the content version is too low")
		}
	}

	_count, err := _article.QueryVersions().Count(ctx)
	if err != nil {
		_count = 1
	}

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		if _published == nil {
			if _, err := tx.Article.UpdateOneID(_body.ArticleID).SetMain(_published).Save(ctx); err != nil {
				return c.InternalServerError(err.Error())
			}
		}

		_tags, _ := _content.QueryTags().IDs(ctx)

		_publishedContent, err := tx.Content.Create().
			SetTitle(_content.Title).
			SetGist(_content.Gist).
			SetBody(_content.Body).
			SetVersion(_content.Version + 1).
			SetVersionName(fmt.Sprintf("v%v.%v.%v", (_content.Version + 1), _count, _lang)).
			SetArticle(_published).
			SetLangID(_lang).
			AddTagIDs(_tags...).
			Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		if _, err = _content.Update().SetMain(_publishedContent).Save(ctx); err != nil {
			return c.InternalServerError(err.Error())
		}

		///////////// Voting and user notification //////////////
		//
		//
		////////////////////////////////////////////////////////

		return c.NoContent()
	})
}
