package main

import (
	"sort"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
)

// Drafts is user draft list
type Drafts []*ent.Draft

func (d Drafts) Len() int { return len(d) }

func (d Drafts) Less(i, j int) bool {
	return d[i].Edges.Snapshots[0].CreatedAt.After(d[j].Edges.Snapshots[0].CreatedAt)
}

func (d Drafts) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func getArticleDrafts(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	_drafts, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.And(
			draft.StateIn(draft.StateRead, draft.StateWrite),
			draft.HasSnapshots())).
		WithSnapshots(func(cq *ent.ContentQuery) {
			cq.Order(ent.Desc(content.FieldID))
		}).
		WithArticle().
		All(ctx)
	if err != nil {
		return c.NotFound(err.Error())
	}

	for _, v := range _drafts {
		_snapshots := make([]*ent.Content, 0)
		if 0 < len(v.Edges.Snapshots) {
			_snapshots = v.Edges.Snapshots[:1]
		}
		v.Edges.Snapshots = _snapshots
	}

	sort.Sort(Drafts(_drafts))

	return c.Ok(&_drafts)
}
