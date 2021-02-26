package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/draft"
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

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_article, err := tx.Article.Create().SetStatus(_query.Status).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_branche, err := tx.Draft.Create().SetArticle(_article).SetUserID(_userID.(int)).SetStatus(draft.StatusWrite).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_, err = _article.Update().AddBranches(_branche).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.Ok(struct {
			ArticleID int
			DraftID   int
		}{
			ArticleID: _article.ID,
			DraftID:   _branche.ID,
		})
	})
}
