package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/ras"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
	"knowlgraph.com/ent/vote"
	"knowlgraph.com/ent/voter"
	"knowlgraph.com/ent/word"
)

// Post vote. Users vote for an article,
// and if the vote is passed, the article will be published.
//
// This method must be called when the user is authorized
func postVote(c *Context) error {
	// ID: ras id
	// Status: vote status
	// Code: Violating code
	// Comment: A description of the violation
	var _query struct {
		ID     int         `json:"id" binding:"required"`
		Status vote.Status `json:"status" binding:"required"`
		Code   []string    `json:"code"`
		// Comment string      `json:"comment"`
	}

	err := c.ShouldBindJSON(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if _query.Code != nil && 0 < len(_query.Code) {
		if _query.Status != vote.StatusOverruled {
			return c.BadRequest("Invalid code")
		}
	} else if _query.Status == vote.StatusOverruled {
		return c.BadRequest("Missing code")
	}

	_userID, _ := c.Get(GinKeyUserID)

	// Match whether the specified user has the right
	// to vote in RAS, and return the matching result.
	// If the match fails, it returns authentication failure (401)
	_ras, err := client.RAS.Query().
		Where(
			ras.ID(_query.ID),
			ras.HasVotersWith(
				voter.Voted(false),
				voter.HasUserWith(user.ID(_userID.(int))))).
		Only(ctx)

	if err != nil {
		return c.Unauthorized(err.Error())
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		// 设置投票结果，不记录表决者，匿名投票
		_, err = tx.Vote.Create().
			SetRasID(_ras.ID).
			SetStatus(_query.Status).
			SetCode(_query.Code).
			// SetComment(_query.Comment).
			Save(ctx)
		if err != nil {
			return err
		}

		// 设置表决者投票状态为已投票。
		// 不记录表决结果，仅记录是否投票。
		// 匿名投票。
		_, err = tx.Voter.Update().
			Where(
				voter.HasRasWith(ras.ID(_query.ID)),
				voter.HasUserWith(user.ID(_userID.(int)))).
			SetVoted(true).
			Save(ctx)
		if err != nil {
			return err
		}

		// 获取未投票的数量
		_unvoteCount, err := _ras.QueryVoters().
			Where(voter.Voted(false)).
			Count(ctx)
		if err != nil {
			return err
		}
		// 在数据库事务处理中，
		// 由于操作数据没有真正被更新，
		// 所以需要手动减一。
		_unvoteCount--

		if 0 == _unvoteCount {
			return votingSettlement(tx, _ras)
		}

		return nil
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(&_query)
}

func votingSettlement(tx *ent.Tx, _ras *ent.RAS) error {
	// 如果所有表决者已投票，则进入表决结算阶段

	_votes, err := tx.Vote.Query().
		Where(
			vote.HasRasWith(ras.ID(_ras.ID))).
		All(ctx)

	if err != nil {
		return err
	}

	// 允许通过的票数
	_allowedCount := 0
	// 反对通过的票数
	_rejectedCount := 0
	// 放弃票决的票数
	_abstainedCount := 0
	for _, _vote := range _votes {
		switch _vote.Status {
		case vote.StatusAllowed:
			_allowedCount++
		case vote.StatusOverruled:
			_rejectedCount++
		case vote.StatusAbstained:
			_abstainedCount++
		}
	}

	// 票数多者为表决结果
	// 三种状态的票数一致时，重新开启一个随机匿名空间
	_voteStatus := ras.StatusAllowed

	if _allowedCount < _rejectedCount && _abstainedCount < _rejectedCount {
		_voteStatus = ras.StatusRejected
	} else if (_allowedCount < _abstainedCount && _rejectedCount < _abstainedCount) ||
		(_allowedCount == _rejectedCount && _allowedCount == _abstainedCount) {
		_voteStatus = ras.StatusAbstained
	}

	// 获取表决文章的版本
	_version, err := _ras.QueryVersion().First(ctx)
	if err != nil {
		return err
	}
	_ras.Edges.Version = _version

	if ras.StatusAbstained == _voteStatus {
		// 如果表决无效，则重新开启一个随机匿名空间

		// 获取原表决者
		// todo 1.0 争议性方案
		// 重新开启的随机匿名空间，
		// 其表决者应是原空间的表决者（以下称原表决者）？
		// 还是应该重新随机选择表决者（以下称新表决者）？
		// 两者各有好处：
		// 原表决者对表决文章已有一定的认识，
		// 在此基础上，再结合上次投票结果，
		// 作出更合理的决定，并加快表决进程。
		// 采用新表决者，则旁观者清，
		// 根据原表决者的投票结果，加上自己的理解，
		// 可以更加理智的作出决定。
		// 两种方案各有优劣，
		// 这里暂时选择第一个方案，
		// 即原空间的表决者。
		_voterIDs, err := _ras.QueryVoters().Select(user.FieldID).Ints(ctx)
		if err != nil {
			return err
		}

		// 开启一个随机匿名空间
		err = openRandomAnonymousSpace(tx, _ras.Comment, _version, _voterIDs)
		if err != nil {
			return err
		}

	} else if ras.StatusAllowed == _voteStatus {
		// 表决结果为允许通过，
		// 即文章可以公示在公共空间。
		// 则删除文章草稿
		// 并更新文章状态，
		// 同时通知所有参与者（发起人和表决者）

		// 获取该版本与关联的草稿
		_branche, err := _version.QueryContent().QueryBranche().First(ctx)
		if err != nil {
			return err
		}

		// 删除草稿
		err = tx.Draft.DeleteOne(_branche).Exec(ctx)
		if err != nil {
			return err
		}

		// 更新文章状态为正式版
		_, err = tx.Version.
			UpdateOne(_version).
			SetState(version.StateRelease).
			Save(ctx)
		if err != nil {
			return err
		}

		_wordIDs, err := _version.QueryKeywords().Select(word.FieldID).Ints(ctx)
		_, err = tx.Word.Update().Where(word.IDIn(_wordIDs...)).SetStatus(word.StatusPublic).Save(ctx)
		if err != nil {
			return err
		}

		// todo 更新文章关键字及与用户的关系
		// todo 通知所有参与者

	} else if ras.StatusRejected == _voteStatus {
		// 表决结果为反对通过，
		// 文章版本状态变更为 "reject"，
		// 且文章草稿由只读变为可写，
		// 并通知所有参与者（发起人和表决者），
		// 同时触发警示系统。

		// 更新文章状态为驳回
		_, err = tx.Version.
			UpdateOne(_version).
			SetState(version.StateReject).
			Save(ctx)
		if err != nil {
			return err
		}

		// 获取与文章版本关联的草稿
		_branche, err := _version.QueryContent().QueryBranche().First(ctx)
		if err != nil {
			return err
		}

		// 设置草稿状态变更为可写
		_, err = tx.Draft.
			UpdateOne(_branche).
			SetState(draft.StateWrite).
			Save(ctx)
		if err != nil {
			return err
		}

		// todo 通知所有参与者
		// todo 触发警示系统
	}

	// 关闭空间（更新随机匿名空间状态），并清除所有表决者
	//
	// todo 暂时方案，此处与 1.0 争议性方案联动
	// 由于重新开启空间时录入的是原表决者，
	// 所以此处清除所有表决者，不会影响后续表决成功后的通知（参与者）效果。
	// 如果采用新表决者，那么此时清除原表决者，
	// 则会出现表决成功后，无法通知原表决者的情况，
	// 故，在切换方案时，一定要考虑此类场景的处理方式。
	_, err = tx.RAS.
		UpdateOne(_ras).
		SetStatus(_voteStatus).
		ClearVoters().
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}
