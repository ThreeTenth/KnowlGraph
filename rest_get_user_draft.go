package main

import (
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getUserDraft(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	id, err := c.GetQueryInt("id")
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_draft, err := client.Version.Query().
		Where(version.And(
			version.HasUserWith(user.IDEQ(_userID.(int))),
			version.IDEQ(id))).
		WithContent().
		Only(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	if _draft.Status == version.StatusReview || _draft.Status == version.StatusRelease {
		return c.Conflict("The draft has been published")
	}

	return c.Ok(_draft)
}
