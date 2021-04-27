package main

import "knowlgraph.com/ent/archive"

func putArchive(c *Context) error {
	var _query struct {
		NodeID int            `form:"nodeId" binding:"required"`
		Status archive.Status `form:"status" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_node, ok, err := getPublicNodeIfExist(_query.NodeID)
	if err != nil {
		return c.NotFound(err.Error())
	} else if !ok {
		return c.MethodNotAllowed("No access to private methods")
	}

	_userID, _ := c.Get(GinKeyUserID)

	_archive, err := client.Archive.
		Create().
		SetNode(_node).
		SetUserID(_userID.(int)).
		SetStatus(_query.Status).
		Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_archive)
}
