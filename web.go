package main

func signout(c *Context) error {
	return c.TemporaryRedirect("/")
}

func index(c *Context) error {
	_userID, ok := c.Get(GinKeyUserID)
	if ok {
		_user, err := client.User.Get(ctx, _userID.(int))
		if err != nil {
			return c.OkHTML(TplIndexHTML, config.Gci)
		}

		return c.OkHTML(TplMainHTML, &_user)
	}

	return c.OkHTML(TplIndexHTML, config.Gci)
}
