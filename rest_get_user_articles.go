package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getUserArticles(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	var _query struct {
		Status asset.Status `form:"status"`
	}

	c.ShouldBindQuery(_query)
	if _query.Status == "" {
		_query.Status = asset.StatusSelf
	}

	_articles, err := client.Asset.Query().
		Where(asset.And(
			asset.HasUserWith(
				user.IDEQ(_userID.(int))),
			asset.StatusEQ(_query.Status))).
		QueryArticle().
		WithVersions(func(vq *ent.VersionQuery) {
			vq.Where(version.StatusEQ(version.StatusRelease)).Order(ent.Desc(version.FieldID)).Limit(1)
		}).
		All(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_articles)
}
