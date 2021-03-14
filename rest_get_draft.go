package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
)

func getDraft(c *Context) error {
	var _query struct {
		ID          int  `form:"id" binding:"required"`
		NeedHistory bool `form:"needHistory"`
		HistoryID   int  `form:"historyID"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_draft, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.And(draft.ID(_query.ID), draft.StateIn(draft.StateRead, draft.StateWrite))).
		WithSnapshots(func(cq *ent.ContentQuery) {
			if 0 < _query.HistoryID {
				cq.Where(content.ID(_query.HistoryID))
			} else if _query.NeedHistory {
				cq.Order(ent.Desc(content.FieldCreatedAt))
			} else {
				cq.Order(ent.Desc(content.FieldCreatedAt)).Limit(1)
			}
		}).
		WithArticle().
		Only(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_draft)
}
