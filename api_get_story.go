package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/story"
)

// getStory returns 200 and if an story is found, if the request fails, it returns a non-200 code
func getStory(c *Context) error {
	storyID, err := c.GetQueryInt("id")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_story, err := client.Story.Get(ctx, storyID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	if story.StatusPrivate == _story.Status {
		return c.Unauthorized("Unauthorized")
	}

	ok := true // True if content exists
	_version := &ent.Content{}

	contentID, err := c.GetQueryInt("contentID")
	if err != nil {
		ok = false
	} else {
		_version, err = _story.QueryVersions().Where(content.IDEQ(contentID)).First(ctx)
		if err != nil {
			ok = false
		}
	}

	if !ok {
		_versionsCreate := _story.QueryVersions()
		_lang, ok := c.GetQuery("lang")
		if ok {
			_versionsCreate.Where(content.HasLangWith(language.IDEQ(_lang)))
		}
		_version, err = _versionsCreate.Order(ent.Desc(content.FieldCreatedAt)).First(ctx)
	}

	if err != nil {
		return c.NotFound(err.Error())
	}

	lang, err := _version.QueryLang().First(ctx)
	if err == nil {
		_version.Edges.Lang = lang
	}

	_tags, err := _version.QueryTags().All(ctx)
	if err == nil {
		_version.Edges.Tags = _tags
	}

	return c.Ok(_version)
}
