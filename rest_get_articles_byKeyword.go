package main

func getKeywordArticles(c *Context) error {
	var _query struct {
		WordID int    `form:"nodeId" binding:"required"`
		Lang   string `form:"lang"`
		Limit  int    `form:"limit,default=10"`
		Offset int    `form:"Offset,default=0"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_versions, err := queryKeywordArticles(db, _query.WordID, _query.Lang, _query.Offset, _query.Limit)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_versions)
}
