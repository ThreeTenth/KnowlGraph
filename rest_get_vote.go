package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/voter"
)

func getVote(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	_vote, err := client.Voter.Query().
		Where(
			voter.VotedEQ(false),
			voter.HasUserWith(user.ID(_userID.(int)))).
		QueryRas().
		WithVersion(func(vq *ent.VersionQuery) {
			vq.WithContent().
				WithKeywords()
		}).
		First(ctx)
	if err != nil {
		return c.NoContent()
	}

	return c.Ok(&_vote)
}
