package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/node"
	"knowlgraph.com/ent/word"
)

// 权限规则是当节点路径上有一个为私有节点时，其后所有节点均为私有。
//
// 获取指定节点的下级节点(nexts)。
// 需要对指定节点进行鉴权，以及对 nexts 鉴权。
func getNodes(c *Context) error {
	var _query struct {
		ID int `form:"id"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if 0 == _query.ID {
		_nodes, err := client.Node.Query().
			Where(node.Level(0)).
			QueryNexts().
			WithWord().
			All(ctx)
		if err != nil {
			return c.NotFound(err.Error())
		}

		_node := ent.Node{
			Edges: ent.NodeEdges{
				Nexts: _nodes,
			},
		}

		return c.Ok(_node)
	}

	_nodeQuery := client.Node.Query()
	if 0 < _query.ID {
		_node, err := _nodeQuery.
			Where(node.ID(_query.ID)).
			WithPath(func(nq *ent.NodeQuery) {
				nq.WithWord()
			}).
			First(ctx)
		if err != nil {
			return c.NotFound(err.Error())
		}

		_path := _node.Edges.Path
		_status := word.StatusPublic
		for _, _pn := range _path {
			if _pn.Edges.Word.Status == word.StatusPrivate {
				_status = word.StatusPrivate
				break
			}
		}

		if _status == word.StatusPrivate {

		}
	}

	_node, err := client.Node.
		Query().
		Where(node.Level(0)).
		First(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_node)
}
