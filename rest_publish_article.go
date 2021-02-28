package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/user"
)

func publishArticle(c *Context) error {
	var _data struct {
		Name      string `json:"name"`
		Comment   string `json:"comment"`
		Title     string `json:"title"`
		Gist      string `json:"gist"`
		Lang      string `json:"lang"`
		ContentID int    `json:"content_id"`
		ArticleID int    `json:"article_id"`
	}

	err := c.ShouldBindJSON(&_data)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_lang, err := client.Language.Query().Where(language.CodeEQ(_data.Lang)).First(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_branche, err := client.Draft.Query().
		Where(draft.HasUserWith(user.ID(_userID.(int)))).
		WithSnapshots(func(cq *ent.ContentQuery) {
			cq.Where(content.ID(_data.ContentID))
		}).
		Only(ctx)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_content := _branche.Edges.Snapshots[0]

	if _data.Title == "" {
		_data.Title = _content.Title
	}

	if _data.Gist == "" {
		_data.Gist = _content.Body
	}

	_version, err := client.Version.Create().
		SetComment(_data.Comment).
		SetTitle(_data.Title).
		SetGist(_data.Gist).
		SetLangID(_lang.ID).
		SetContentID(_data.ContentID).
		SetArticleID(_data.ArticleID).
		Save(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(&_version)
}
