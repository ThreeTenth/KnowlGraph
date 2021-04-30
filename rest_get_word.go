package main

import "knowlgraph.com/ent/word"

func getWord(c *Context) error {
	var _query struct {
		ID   int    `form:"id"`
		Name string `form:"id"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_where := word.ID(_query.ID)
	if 0 == _query.ID {
		_where = word.Name(_query.Name)
	}
	_word, err := client.Word.
		Query().
		Where(_where).
		WithDefinition().
		First(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_word)
}
