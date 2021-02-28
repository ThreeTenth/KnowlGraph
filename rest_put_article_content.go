package main

func putArticleContent(c *Context) error {
	var _data struct {
		Title     string
		Body      string
		LastID    int
		BrancheID int
	}

	err := c.ShouldBindJSON(&_data)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_content, err := client.Content.Create().
		SetTitle(_data.Title).
		SetBody(_data.Body).
		SetLastID(_data.LastID).
		SetBrancheID(_data.BrancheID).
		Save(ctx)

	return c.Ok(_content)
}
