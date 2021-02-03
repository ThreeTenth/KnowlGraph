package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/user"
)

// getArticle returns 200 and if an article is found, if the request fails, it returns a non-200 code
func getArticle(c *Context) error {
	var _query struct {
		ID        int `form:"id" binding:"required"`
		ContentID int `form:"content_id"`
	}
	err := c.ShouldBindQuery(&_query)

	_article, err := client.Article.Get(ctx, _query.ID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	if article.StatusPrivate == _article.Status {
		_userID, ok := c.Get(GinKeyUserID)
		if !ok {
			return c.Unauthorized("No access")
		}
		ok, err = client.User.Query().
			Where(user.IDEQ(_userID.(int))).
			QueryArticles().
			Where(article.IDEQ(_article.ID)).
			Exist(ctx)
		if !ok || err != nil {
			return c.Unauthorized("No access")
		}
	}

	ok := 0 != _query.ContentID // True if content exists
	_version := &ent.Content{}

	if ok {
		_version, err = _article.QueryVersions().Where(content.IDEQ(_query.ContentID)).WithLang().WithTags().First(ctx)
		if err != nil {
			ok = false
		}
	}

	if !ok {
		_versionsCreate := _article.QueryVersions()
		_lang, ok := c.GetQuery(QueryLang)
		if ok {
			_versionsCreate.Where(content.HasLangWith(language.IDEQ(_lang))).WithLang().WithTags()
		}
		_version, err = _versionsCreate.Order(ent.Desc(content.FieldCreatedAt)).First(ctx)
	}

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_version)
}
