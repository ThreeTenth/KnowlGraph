package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/ras"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/vote"
)

func postVote(c *Context) error {
	var _query struct {
		ID      int
		Status  vote.Status
		Code    vote.Code
		Comment string
	}

	err := c.ShouldBindJSON(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_ras, err := client.RAS.Query().
		Where(
			ras.ID(_query.ID),
			ras.HasVotersWith(
				user.ID(_userID.(int)))).
		Only(ctx)

	if err != nil {
		return c.Unauthorized(err.Error())
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		_, err = client.Vote.Create().
			SetRasID(_ras.ID).
			SetStatus(_query.Status).
			SetCode(_query.Code).
			SetComment(_query.Comment).
			Save(ctx)
		if err != nil {
			return err
		}

		_, err = _ras.Update().
			RemoveVoterIDs(_userID.(int)).
			Save(ctx)
		if err != nil {
			return err
		}

		_voterCount, err := _ras.QueryVoters().Count(ctx)
		if err != nil {
			return err
		}

		if 0 == _voterCount {
			// 统计票数并更新文章状态
		}

		return nil
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(&_query)
}
