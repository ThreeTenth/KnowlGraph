package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
)

func putArticleNew(c *Context) error {
	var _query struct {
		Status article.Status `form:"status,default=public"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_draft, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.And(
			draft.Not(draft.HasSnapshots()),
			draft.Not(draft.HasHead()),
			draft.StateEQ(draft.StateWrite))).
		Limit(1).
		WithArticle().
		First(ctx)

	if err == nil {
		return c.Ok(struct {
			ArticleID int
			DraftID   int
		}{
			ArticleID: _draft.Edges.Article.ID,
			DraftID:   _draft.ID,
		})
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		_article, err := tx.Article.Create().SetStatus(_query.Status).Save(ctx)
		if err != nil {
			return err
		}

		_branche, err := tx.Draft.Create().SetArticle(_article).SetUserID(_userID.(int)).SetState(draft.StateWrite).Save(ctx)
		if err != nil {
			return err
		}

		return c.Ok(struct {
			ArticleID int
			DraftID   int
		}{
			ArticleID: _article.ID,
			DraftID:   _branche.ID,
		})
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return err
}
