package main

import (
	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/archive"
	"knowlgraph.com/ent/node"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/userword"
	"knowlgraph.com/ent/word"
)

func putNode(c *Context) error {
	// 创建节点的接口参数
	// WordID: 节点指向的关键字
	// NodeID: 上级节点，若无，则创建的节点自动为根节点
	var _query struct {
		WordID int `json:"wordId" binding:"required"`
		NodeID int `json:"nodeId"`
	}

	err := c.ShouldBindJSON(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	// 获取关键字信息
	// 如果关键字是私有的，那么需要对关键字鉴权
	_word, err := client.Word.Get(ctx, _query.WordID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	_userID, ok := c.Get(GinKeyUserID)
	_nodeStatus := node.StatusPublic

	if _word.Status == word.StatusPrivate {
		// 关键字是私有的

		if !ok {
			// 未登录
			return c.Unauthorized("Unlogin")
		}

		// 查询指定的私有关键字是否为用户所有
		ok, err := client.UserWord.Query().
			Where(
				userword.HasUserWith(
					user.ID(_userID.(int))),
				userword.HasWordWith(
					word.ID(_query.WordID))).
			Exist(ctx)

		if err != nil {
			return c.InternalServerError(err.Error())
		}

		if !ok {
			// 鉴权失败

			return c.Unauthorized("Unauthorized")
		}

		// 当节点指向的关键字为私有时，节点便为私有
		_nodeStatus = node.StatusPrivate
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		// 创建节点

		_nodeCreate := tx.Node.Create().SetWordID(_query.WordID)

		if 0 < _query.NodeID {
			// 如果存在上级节点，则需要鉴权

			_prev, ok, err := getPublicNodeIfExist(_query.NodeID)
			if err != nil {
				return c.NotFound(err.Error())
			}

			if !ok {
				// 如果上级节点的路径中存在私有节点

				// 查询用户是否拥有指定节点
				ok, err = client.Archive.
					Query().
					Where(
						archive.HasNodeWith(node.ID(_prev.ID)),
						archive.HasUserWith(user.ID(_userID.(int)))).
					Exist(ctx)
				if err != nil {
					return err
				}

				if !ok {
					return errors.New("No access to private node")
				}
			}

			_path, _ := _prev.QueryPath().All(ctx)

			// 为创建的节点添加节点路径
			_path = append(_path, _prev)

			_nodeCreate.SetPrev(_prev).AddPath(_path...)

			if node.StatusPublic == _nodeStatus {
				_nodeStatus = _prev.Status
			}
		}

		// 开始创建节点
		_node, err := _nodeCreate.SetStatus(_nodeStatus).Save(ctx)

		if err != nil {
			return err
		}

		// 为用户添加节点
		_archive, err := tx.Archive.
			Create().
			SetNode(_node).
			SetUserID(_userID.(int)).
			Save(ctx)

		if err != nil {
			return err
		}

		_archive.Edges.Node = _node

		return c.Ok(_archive)
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return nil
}
