package main

import "knowlgraph.com/ent/node"

func putNode(c *Context) error {
	// 创建节点的接口参数
	// WordID: 节点指向的关键字
	// NodeID: 上级节点，若无，则创建的节点自动为根节点
	var _query struct {
		WordID int `form:"wordID"`
		NodeID int `form:"nodeId"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_nodeCreate := client.Node.Create().SetWordID(_query.WordID)

	if 0 < _query.NodeID {
		_prev, err := client.Node.Query().Where(node.ID(_query.NodeID)).WithPath().First(ctx)
		if err != nil {
			return c.BadRequest(err.Error())
		}

		_path := append(_prev.Edges.Path, _prev)

		_nodeCreate.SetPrev(_prev).SetLevel(_prev.Level + 1).AddPath(_path...)
	}

	_node, err := _nodeCreate.Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_node)
}
