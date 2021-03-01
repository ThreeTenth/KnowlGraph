package main

import (
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
)

func getUserArticles(c *Context) error {
	var _query struct {
		Status asset.Status `from:"status"`
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
		All(ctx)

	return c.Ok(_articles)
}
