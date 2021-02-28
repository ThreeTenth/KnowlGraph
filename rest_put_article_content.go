package main

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

	_content, err := client.Content.Create().
		SetTitle(_data.Title).
		SetBody(_data.Body).
		SetLastID(_data.LastID).
		SetBrancheID(_data.DraftID).
		Save(ctx)

	return c.Ok(_content)
}
