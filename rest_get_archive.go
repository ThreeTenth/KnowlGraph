package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/archive"
	"knowlgraph.com/ent/user"
)

func getArchive(c *Context) error {
	var _query struct {
		ID     int            `form:"id" binding:"required"`
		status archive.Status `form:"status" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_archive, err := client.Archive.
		Query().
		Where(
			archive.StatusEQ(_query.status),
			archive.HasUserWith(
				user.ID(_userID.(int)))).
		WithNode(func(nq *ent.NodeQuery) {
			nq.WithWord().
				WithPath(func(nq *ent.NodeQuery) {
					nq.WithWord()
				})
		}).
		First(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_archive)
}
