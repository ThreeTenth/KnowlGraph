package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/archive"
	"knowlgraph.com/ent/predicate"
	"knowlgraph.com/ent/user"
)

func getArchives(c *Context) error {
	var _query struct {
		Status archive.Status `form:"status"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	var _archiveWhere []predicate.Archive
	_archiveWhere = append(
		_archiveWhere,
		archive.HasUserWith(
			user.ID(_userID.(int))))
	if _query.Status.String() != "" {
		_archiveWhere = append(_archiveWhere, archive.StatusEQ(_query.Status))
	}

	_archive, err := client.Archive.
		Query().
		Where(_archiveWhere...).
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
