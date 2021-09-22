package main

import (
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
)

func putQuote(c *Context) error {
	var _query struct {
		Text       string `json:"text"`
		Context    string `json:"context"`
		Highlight  int    `json:"highlight"`
		ArticleID  int    `json:"articleId" binding:"required"`
		VersionID  int    `json:"versionId"`
		ResponseID int    `json:"responseId" binding:"required"`
	}

	err := c.ShouldBindJSON(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)
	ok, err := client.Article.Query().
		Where(
			article.And(
				article.And(
					article.ID(_query.ResponseID),
					article.Or(
						article.StatusNEQ(article.StatusPrivate),
						article.And(
							article.StatusEQ(article.StatusPrivate),
							article.HasAssetsWith(asset.HasUserWith(user.ID(_userID.(int)))),
						)),
				),
				article.And(
					article.ID(_query.ArticleID),
					article.Or(
						article.StatusNEQ(article.StatusPrivate),
						article.And(
							article.StatusEQ(article.StatusPrivate),
							article.HasAssetsWith(asset.HasUserWith(user.ID(_userID.(int)))),
						)),
				),
			),
		).Exist(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	} else if !ok {
		return c.Unauthorized("Unauthorized")
	}

	_quote, err := client.Quote.Create().
		SetText(_query.Text).
		SetContext(_query.Context).
		SetHighlight(_query.Highlight).
		SetArticleID(_query.ArticleID).
		SetVersionID(_query.VersionID).
		SetResponseID(_query.ResponseID).
		Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_quote.ID)
}
