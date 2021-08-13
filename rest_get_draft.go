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
		Where(draft.And(draft.ID(_query.ID), draft.StateIn(draft.StateWrite))).
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
		WithOriginal(func(vq *ent.VersionQuery) {
			vq.WithContent().WithKeywords()
		}).
		Only(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	snapshots := _draft.Edges.Snapshots
	if nil == snapshots {
		snapshots = make([]*ent.Content, 0)
	}

	if 0 == len(snapshots) {
		snapshots = append(snapshots, _draft.Edges.Original.Edges.Content)
		_draft.Edges.Snapshots = snapshots
	}

	return c.Ok(&_draft)
}
