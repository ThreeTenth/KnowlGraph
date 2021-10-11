package main

import (
	"knowlgraph.com/ent/terminal"
	"knowlgraph.com/ent/user"
)

func deleteTerminal(c *Context) error {
	var form struct {
		ID int `form:"id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	userID, _ := c.Get(GinKeyUserID)
	token := getRequestToken(c.Context)

	if terminalID, err := rdb.Get(ctx, RTerminal(token)).Int(); err != nil {
		return c.NotFound(err.Error())
	} else if terminalID == form.ID {
		count, err := client.Terminal.
			Query().
			Where(terminal.HasUserWith(user.ID(userID.(int)))).
			Count(ctx)

		if err != nil {
			return c.InternalServerError(err.Error())
		}

		if 1 < count {
			return c.PreconditionFailed("Please delete other terminals first")
		}
	}

	if err := rdb.Del(ctx, RToken(token)).Err(); err != nil {
		return c.InternalServerError(err.Error())
	}

	count, err := client.Terminal.
		Delete().
		Where(terminal.ID(form.ID), terminal.HasUserWith(user.ID(userID.(int)))).
		Exec(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if 0 == count {
		return c.NotFound("The terminal ID was not found")
	}

	return c.Ok(true)
}
