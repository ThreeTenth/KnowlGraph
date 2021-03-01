package main

import (
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
)

func putArticleContent(c *Context) error {
	var _data struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		LastID  int    `json:"last_id"`
		DraftID int    `json:"draft_id"`
	}

	err := c.ShouldBindJSON(&_data)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_, err = client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.ID(_data.DraftID)).
		Only(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_content, err := client.Content.Create().
		SetTitle(_data.Title).
		SetBody(_data.Body).
		SetLastID(_data.LastID).
		SetBrancheID(_data.DraftID).
		Save(ctx)

	return c.Ok(_content)
}
