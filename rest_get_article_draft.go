package main

import (
	"knowlgraph.com/ent/user"
)

func getArticleBranche(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	_drafts, err := client.User.Query().Where(user.ID(_userID.(int))).QueryDrafts().All(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_drafts)
}
