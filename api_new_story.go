package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/story"
)

// newStory creates a new story, the default language is en
func newStory(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	return WithTx(ctx, client, func(tx *ent.Tx) error {
		_user, err := tx.User.Get(ctx, _userID.(int))
		if err != nil {
			return c.Forbidden(err.Error())
		}

		// Return articles without any content first
		_storyID, err := _user.QueryStories().Where(story.Not(story.HasVersions())).FirstID(ctx)
		if err == nil {
			return c.Ok(&_storyID)
		}

		_storyCreater := tx.Story.Create().SetStatus(story.StatusPrivate)

		_story, err := _storyCreater.Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		_userID, err = tx.User.UpdateOneID(_userID.(int)).AddStories(_story).Save(ctx)
		if err != nil {
			return c.InternalServerError(err.Error())
		}

		return c.Ok(&_story.ID)
	})
}
