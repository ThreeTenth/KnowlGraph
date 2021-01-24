package main

import (
	"fmt"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/story"
	"knowlgraph.com/ent/tag"
)

// NewStory creates a new story, the default language is en
func NewStory(c *Context) error {

	///////////////// Story authentication //////////////////
	//
	//
	////////////////////////////////////////////////////////

	_storyCreater := client.Story.Create().SetStatus(story.StatusPrivate)

	_story, err := _storyCreater.Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}
	return c.Ok(&_story.ID)
}

// PutStoryContent creates a content version for the specified story
func PutStoryContent(c *Context) error {

	///////////////// Story authentication //////////////////
	//
	//
	////////////////////////////////////////////////////////

	var _body struct {
		Title       string   `json:"title,omitempty"`
		Gist        string   `json:"gist,omitempty"`
		Content     string   `json:"content"`
		VersionName string   `json:"versionName,omitempty"`
		StoryID     int      `json:"storyID"`
		Tags        []string `json:"tags,omitempty" comment:"5 at most"`
	}
	err := c.ShouldBindJSON(&_body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	if nil != _body.Tags && 5 < len(_body.Tags) {
		return c.BadRequest("Up to 5 tags")
	}

	_story, err := client.Story.Get(ctx, _body.StoryID)
	if err != nil {
		return c.NotFound(err.Error())
	}

	_lang := c.DefaultQuery("lang", DefaultLanguage)

	var columns []struct {
		Lang string `json:"content_lang"`
		Max  int    `json:"max"`
	}

	_story.QueryVersions().
		Where(content.HasLangWith(language.ID(_lang))).
		GroupBy(content.LangColumn).
		Aggregate(ent.Max(content.FieldVersion)).
		Scan(ctx, &columns)

	maxVersion := 1
	if 1 == len(columns) && maxVersion <= columns[0].Max {
		maxVersion = columns[0].Max
	}

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_contentCreate := tx.Content.Create().
			SetTitle(_body.Title).
			SetGist(_body.Gist).
			SetContent(_body.Content).
			SetStory(_story).
			SetVersion(maxVersion).
			SetLangID(_lang)

		if nil != _body.Tags && 0 < len(_body.Tags) {
			_tags := make([]*ent.Tag, len(_body.Tags))
			for i, _tag := range _body.Tags {
				if "" == _tag {
					continue
				}
				_tags[i], err = tx.Tag.Query().Where(tag.NameContainsFold(_tag)).First(ctx)
				if err != nil {
					_tags[i], err = tx.Tag.Create().SetName(_tag).Save(ctx)
					if nil != err {
						return c.InternalServerError(err.Error())
					}
				}
			}
			_contentCreate.AddTags(_tags...)
		}

		_content, err := _contentCreate.Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.Ok(_content.ID)
	})
}

// PublishStory pushes a personal story content to a public channel
func PublishStory(c *Context) error {
	var _body struct {
		StoryID   int `json:"storyID"`
		ContentID int `json:"ContentID"`
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

// GetStory returns 200 and if an story is found, if the request fails, it returns a non-200 code
func GetStory(c *Context) error {
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
