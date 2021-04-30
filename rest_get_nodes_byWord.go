package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/node"
	"knowlgraph.com/ent/word"
)

func getWordNodes(c *Context) error {
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

	_nodes, err := client.Node.
		Query().
		Where(
			node.HasWordWith(word.StatusEQ(word.StatusPublic)),
			node.HasWordWith(_where)).
		WithPath(func(nq *ent.NodeQuery) {
			nq.WithWord()
		}).
		WithWord().
		All(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	for _, _node := range _nodes {
		if _node.Edges.Word.Status == word.StatusPrivate {
			_node.Edges.Word = nil
			_node.Edges.Path = nil
		}
		for _, m := range _node.Edges.Path {
			if m.Edges.Word.Status == word.StatusPrivate {
				_node.Edges.Word = nil
				_node.Edges.Path = nil
				break
			}
		}
	}

	return c.Ok(_nodes)
}
