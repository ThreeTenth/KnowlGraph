package main

import (
	"fmt"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/story"
)

// publishStory pushes a personal story content to a public channel
func publishStory(c *Context) error {
	var _body struct {
		StoryID   int `json:"storyID" binding:"required,gt=0,storyExist"`
		ContentID int `json:"ContentID" binding:"required,gt=0"`
	}
	err := c.ShouldBindJSON(&_body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	///////////////// User authentication //////////////////
	//
	//
	////////////////////////////////////////////////////////

	_story, err := client.Story.Get(ctx, _body.StoryID)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if _story.Status != story.StatusPrivate {
		return c.Forbidden("The operation was blocked because the content is not private")
	}

	//////////////// Story authentication //////////////////
	//
	//
	////////////////////////////////////////////////////////

	_content, err := _story.QueryVersions().Where(content.ID(_body.ContentID)).First(ctx)
	if err != nil {
		return c.Forbidden("The operation was blocked because the content does not belong to the story")
	}

	_published, _ := _story.QueryMain().First(ctx)

	_lang, _ := _content.QueryLang().FirstID(ctx)

	if _published != nil {
		ok, _ := _published.QueryVersions().
			Where(
				content.And(
					content.VersionGT(_content.Version),
					content.HasLangWith(language.ID(_lang)))).
			Exist(ctx)
		if ok {
			return c.Created("Please update the latest version before publishing, because the content version is too low")
		}
	}

	_count, err := _story.QueryVersions().Count(ctx)
	if err != nil {
		_count = 1
	}

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		if _published == nil {
			if _, err := tx.Story.UpdateOneID(_body.StoryID).SetMain(_published).Save(ctx); err != nil {
				return c.InternalServerError(err.Error())
			}
		}

		_tags, _ := _content.QueryTags().IDs(ctx)

		_publishedContent, err := tx.Content.Create().
			SetTitle(_content.Title).
			SetGist(_content.Gist).
			SetContent(_content.Content).
			SetVersion(_content.Version + 1).
			SetVersionName(fmt.Sprintf("v%v.%v.%v", (_content.Version + 1), _count, _lang)).
			SetStory(_published).
			SetLangID(_lang).
			AddTagIDs(_tags...).
			Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		if _, err = _content.Update().SetMain(_publishedContent).Save(ctx); err != nil {
			return c.InternalServerError(err.Error())
		}

		///////////// Voting and user notification //////////////
		//
		//
		////////////////////////////////////////////////////////

		return c.NoContent()
	})
}
