package main

import (
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
)

func putArticleContent(c *Context) error {
	var _data struct {
		Body    string `json:"body"`
		LastID  int    `json:"last_id"`
		DraftID int    `json:"draft_id" binding:"required"`
	}

	err := c.ShouldBindJSON(&_data)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_draft, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.ID(_data.DraftID)).
		Only(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if _draft.State == draft.StateRead {
		return c.Unauthorized("Cannot be modified because it is read-only")
	}

	_contentCreate := client.Content.Create()
	if 0 != _data.LastID {
		_contentCreate.SetLastID(_data.LastID)
	}

	_content, err := _contentCreate.
		SetBody(_data.Body).
		SetBrancheID(_data.DraftID).
		Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_content)
}
