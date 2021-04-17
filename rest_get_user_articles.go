package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/asset"
)

// Articles are used to sort articles
type Articles []*ent.Article

func (d Articles) Len() int { return len(d) }

func (d Articles) Less(i, j int) bool {
	return d[i].Edges.Versions[0].CreatedAt.
		After(
			d[j].Edges.Versions[0].CreatedAt)
}

func (d Articles) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func getUserArticles(c *Context) error {
	var _query struct {
		Status asset.Status `form:"status"`
		Lang   string       `form:"lang"`
		Limit  int          `form:"limit,default=10"`
		Offset int          `form:"Offset,default=0"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_versions, err := queryUserAssets(db, _userID.(int), _query.Status, _query.Lang, _query.Offset, _query.Limit)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_versions)
}
