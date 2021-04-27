package main

import (
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
)

func putAsset(c *Context) error {
	var _query struct {
		ArticleID int          `form:"articleID" binding:"required"`
		Status    asset.Status `form:"status" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_article, err := client.Article.Get(ctx, _query.ArticleID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	if _article.Status == article.StatusPrivate {
		return c.Unauthorized("No access to private article")
	}

	_userID, _ := c.Get(GinKeyUserID)

	_asset, err := client.Asset.
		Create().
		SetArticle(_article).
		SetUserID(_userID.(int)).
		SetStatus(_query.Status).
		Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_asset)
}
