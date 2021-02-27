package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/user"
)

func getArticleBranches(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	_drafts, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		WithSnapshots(func(cq *ent.ContentQuery) {
			cq.Order(ent.Desc(content.FieldCreatedAt)).Limit(1)
		}).
		WithArticle().
		All(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_drafts)
}
