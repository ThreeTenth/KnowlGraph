package main

import (
	"sort"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

// Articles are used to sort articles
type Articles []*ent.Asset

func (d Articles) Len() int { return len(d) }

func (d Articles) Less(i, j int) bool {
	return d[i].Edges.Article.Edges.Versions[0].CreatedAt.
		After(
			d[j].Edges.Article.Edges.Versions[0].CreatedAt)
}

func (d Articles) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func getUserArticles(c *Context) error {
	var _query struct {
		Status asset.Status `form:"status"`
		Lang   string       `form:"lang"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_articles, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryAssets().
		Where(asset.StatusEQ(_query.Status)).
		WithArticle(func(aq *ent.ArticleQuery) {
			aq.WithVersions(func(vq *ent.VersionQuery) {
				if _query.Lang != "" {
					vq = vq.Where(version.Lang(_query.Lang))
				}
				vq.Order(ent.Desc(version.FieldCreatedAt))
			})
		}).
		All(ctx)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	sort.Sort(Articles(_articles))

	return c.Ok(_articles)
}
