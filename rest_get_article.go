package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getArticle(c *Context) error {
	var _query struct {
		ID           int    `form:"id" binding:"required"`
		Lang         string `form:"lang"`
		NeedVersions bool   `form:"needVersions"`
		VersionID    int    `form:"versionID"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_article, err := client.Asset.Query().
		Where(asset.And(
			asset.HasUserWith(user.ID(_userID.(int))),
			asset.HasArticleWith(article.ID(_query.ID)))).
		QueryArticle().
		WithVersions(func(vq *ent.VersionQuery) {
			if 0 < _query.VersionID {
				vq.Where(version.ID(_query.VersionID))
			} else if _query.NeedVersions {
				vq.Order(ent.Desc(version.FieldCreatedAt))
			} else {
				vq.Order(ent.Desc(version.FieldCreatedAt)).Limit(1)
			}
			vq.WithContent()
		}).
		Only(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_article)
}
