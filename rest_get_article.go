package main

import (
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

// getArticle returns 200 and if an article is found, if the request fails, it returns a non-200 code
func getArticle(c *Context) error {
	var _query struct {
		ID        int    `form:"id" binding:"required"`
		VersionID int    `form:"version_id"`
		Lang      string `form:"lang"`
	}
	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_article, err := client.Article.Get(ctx, _query.ID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	if article.StatusPrivate == _article.Status {
		_userID, _ := c.Get(GinKeyUserID)
		_, err := client.Asset.Query().
			Where(asset.And(
				asset.HasUserWith(user.IDEQ(_userID.(int))),
				asset.HasArticleWith(article.IDEQ(_query.ID)),
				asset.StatusEQ(asset.StatusSelf))).
			First(ctx)

		if err != nil {
			return c.Unauthorized(err.Error())
		}
	}

	_version, err := client.Version.Query().
		Where(version.And(
			version.HasArticleWith(article.IDEQ(_query.ID))),
			version.Or(
				version.IDEQ(_query.VersionID),
				version.HasLangWith(language.CodeEQ(_query.Lang)))).
		WithContent().
		WithTags().
		WithLang().
		First(ctx)

	return c.Ok(_version)
}
