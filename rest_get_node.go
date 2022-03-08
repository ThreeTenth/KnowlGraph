package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/node"
)

// 权限规则是当节点路径上有一个为私有节点时，其后所有节点均为私有。
//
// 获取指定节点的下级公开节点(nexts)。
//
// 需要对指定节点进行鉴权，如果指定节点为私有，则返回 405（方法拒绝访问）。
func getNode(c *Context) error {
	var _query struct {
		ID int `form:"id"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	// 访问公共节点的查询条件
	_nextsWhere := node.StatusEQ(node.StatusPublic)

	// 如果未指定节点，则返回所有根节点
	if 0 == _query.ID {
		_nexts, err := client.Node.Query().
			Where(_nextsWhere, node.Not(node.HasPrev())).
			WithWord().
			All(ctx)
		if err != nil {
			return c.NotFound(err.Error())
		}

		_node := ent.Node{
			Edges: ent.NodeEdges{
				Nexts: _nexts,
			},
		}

		return c.Ok(_node)
	}

	_node, ok, err := getPublicNodeIfExist(_query.ID)
	if err != nil {
		return c.NotFound(err.Error())
	} else if !ok {
		return c.MethodNotAllowed("No access to private methods")
	}

	// 获取指定节点下，所有的下一节点（nexts）
	_nexts, err := _node.
		QueryNexts().
		Where(_nextsWhere).
		WithWord().
		All(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	_node.Edges.Nexts = _nexts

	return c.Ok(_node)
}

// todo: bug: 此处节点非公开但是用户私有的，那么用户就拥有访问权限
// 应对用户进行鉴权，而不非公开节点统一返回 405 。
func getPublicNodeIfExist(id int) (*ent.Node, bool, error) {
	// 获取指定节点的详细信息
	_node, err := client.Node.Query().
		Where(node.ID(id)).
		WithWord().
		First(ctx)
	if err != nil {
		return nil, false, err
	}

	// 如果指定节点为私有节点，则返回 "MethodNotAllowed(405)"
	if _node.Status == node.StatusPrivate {
		return _node, false, nil
	}

	return _node, true, nil
}
