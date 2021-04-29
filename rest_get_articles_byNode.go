package main

func getNodeArticles(c *Context) error {
	var _query struct {
		NodeID int    `form:"nodeId" binding:"required"`
		Lang   string `form:"lang"`
		Limit  int    `form:"limit,default=10"`
		Offset int    `form:"Offset,default=0"`
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

	_versions, err := queryNodeArticles(db, _node.ID, _query.Lang, _query.Offset, _query.Limit)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_versions)
}
