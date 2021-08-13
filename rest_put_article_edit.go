package main

import (
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func editArticle(c *Context) error {
	var _query struct {
		ID int `form:"id" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_version, err := client.Version.Query().
		Where(version.ID(_query.ID)).
		WithArticle().
		WithContent().
		Only(ctx)

	if err != nil {
		return err
	}

	_draft, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.And(
			draft.HasOriginalWith(version.ID(_version.ID)),
			draft.StateEQ(draft.StateWrite))).
		First(ctx)

	if err != nil {
		_draft, err = client.Draft.Create().
			SetArticle(_version.Edges.Article).
			SetOriginal(_version).
			SetUserID(_userID.(int)).
			Save(ctx)

		if err != nil {
			return c.InternalServerError(err.Error())
		}

	}

	return c.Ok(&_draft)
}
