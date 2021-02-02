package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
)

// newArticle creates a new article, the default language is en
func newArticle(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_user, err := tx.User.Get(ctx, _userID.(int))
		if err != nil {
			return c.Forbidden(err.Error())
		}

		// Return articles without any content first
		_articleID, err := _user.QueryArticles().Where(article.Not(article.HasVersions())).FirstID(ctx)
		if err == nil {
			return c.Ok(&_articleID)
		}

		_articleCreater := tx.Article.Create().SetStatus(article.StatusPrivate)

		_article, err := _articleCreater.Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_userID, err = tx.User.UpdateOneID(_userID.(int)).AddArticles(_article).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.Ok(&_article.ID)
	})
}
