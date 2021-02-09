package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/version"
)

// newDraft creates a new article, the default language is en
func newDraft(c *Context) error {
	var _query struct {
		Status article.Status `form:"status" binding:"required"`
		Lang   article.Status `form:"lang,default=en"`
	}
	_userID, _ := c.Get(GinKeyUserID)

	if err := c.ShouldBindQuery(&_query); err != nil {
		return c.BadRequest(err.Error())
	}

	if err := article.StatusValidator(_query.Status); err != nil {
		return c.BadRequest(err.Error())
	}

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_user, err := tx.User.Get(ctx, _userID.(int))
		if err != nil {
			return c.Forbidden(err.Error())
		}

		// Return articles without any content first
		_id, err := _user.QueryDrafts().
			Where(version.And(
				version.Not(version.HasContent()),
				version.StatusEQ(version.DefaultStatus),
				version.HasArticleWith(article.StatusEQ(_query.Status)))).
			FirstID(ctx)
		if err == nil {
			return c.Ok(&_id)
		}

		_articleCreater := tx.Article.Create().SetStatus(_query.Status)

		_article, err := _articleCreater.Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		var _langID int
		for _, _lang := range langs {
			if _lang.Code == string(_query.Lang) {
				_langID = _lang.ID
			}
		}

		_version, err := tx.Version.Create().SetStatus(version.DefaultStatus).SetLangID(_langID).SetArticle(_article).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_userID, err = tx.User.UpdateOneID(_userID.(int)).AddDrafts(_version).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.Ok(&_article.ID)
	})
}
