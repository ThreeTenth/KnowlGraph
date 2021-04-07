package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/ras"
	"knowlgraph.com/ent/user"
)

func getVote(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	_vote, err := client.RAS.Query().
		Where(ras.HasVotersWith(
			user.ID(_userID.(int)))).
		WithVersion(func(vq *ent.VersionQuery) {
			vq.WithContent().
				WithKeywords()
		}).
		First(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_vote)
}
