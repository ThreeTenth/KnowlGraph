package main

import (
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/reaction"
	"knowlgraph.com/ent/user"
)

func putReaction(c *Context) error {
	var _query struct {
		ArticleID int             `form:"articleId" binding:"required"`
		Reaction  reaction.Status `form:"reaction" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	ok, err := client.Article.
		Query().
		Where(article.Or(
			article.StatusNEQ(article.StatusPrivate),
			article.HasAssetsWith(
				asset.HasUserWith(
					user.ID(_userID.(int)))))).
		Exist(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if !ok {
		return c.NotFound("NotFound")
	}

	_reaction, err := client.Reaction.
		Query().
		Where(
			reaction.StatusEQ(_query.Reaction),
			reaction.HasArticleWith(article.ID(_query.ArticleID))).
		Only(ctx)

	if err != nil {
		_reaction, err = client.Reaction.
			Create().
			SetArticleID(_query.ArticleID).
			SetStatus(_query.Reaction).
			SetCount(1).
			Save(ctx)
	} else {
		_reaction, err = client.Reaction.
			UpdateOne(_reaction).
			AddCount(1).
			Save(ctx)
	}

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_reaction)
}
