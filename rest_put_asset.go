package main

import (
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
)

func putAsset(c *Context) error {
	var _query struct {
		ArticleID int          `form:"articleId" binding:"required"`
		VersionID int          `form:"versionId"`
		Status    asset.Status `form:"status" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if asset.StatusSelf == _query.Status {
		return c.Forbidden("Forbid adding articles to my list")
	}

	_article, err := client.Article.Get(ctx, _query.ArticleID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	if _article.Status == article.StatusPrivate {
		// 如果文章是私有的，则判断该文章是否为用户所有

		ok, _ := client.Asset.
			Query().
			Where(
				asset.HasArticleWith(article.ID(_query.ArticleID)),
				asset.HasUserWith(user.ID(_userID.(int))),
				asset.StatusEQ(asset.StatusSelf)).
			Exist(ctx)

		if !ok {
			return c.Unauthorized("No access to private article")
		}
	}

	if asset.StatusBrowse != _query.Status {
		// 已关注或收藏的文章，无需再次关注或收藏
		ok, _ := client.Asset.
			Query().
			Where(
				asset.HasArticleWith(article.ID(_query.ArticleID)),
				asset.HasUserWith(user.ID(_userID.(int))),
				asset.StatusEQ(_query.Status)).
			Exist(ctx)

		if ok {
			return c.NoContent()
		}
	}

	_assetCreate := client.Asset.Create().SetArticle(_article).SetUserID(_userID.(int)).SetStatus(_query.Status)
	if 0 < _query.VersionID {
		_assetCreate.SetVersionID(_query.VersionID)
	}
	_asset, err := _assetCreate.Save(ctx)

	// _asset, err := client.Asset.
	// 	Create().
	// 	SetArticle(_article).
	// 	SetVersionID(_query.VersionID).
	// 	SetUserID(_userID.(int)).
	// 	SetStatus(_query.Status).
	// 	Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_asset)
}
