package main

import (
	"knowlgraph.com/ent/terminal"
	"knowlgraph.com/ent/user"
)

func getAccountTerminals(c *Context) error {
	_userID := c.GetInt(GinKeyUserID)
	ts, err := client.Terminal.Query().
		Where(
			terminal.HasUserWith(user.ID(_userID)),
		).
		All(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if 0 == len(ts) {
		return c.NotFound("No terminals")
	}

	return c.Ok(&ts)
}
