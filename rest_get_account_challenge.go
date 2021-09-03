package main

func getAccountChallenge(c *Context) error {
	challenge := New64BitID()
	err := rdb.Set(ctx, challenge, true, ExpireTimeChallenge).Err()

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(challenge)
}

func getAccountCreate(c *Context) error {
	var body struct {
		Challenge string `json:"challenge"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	err = rdb.Get(ctx, body.Challenge).Err()

	if err != nil {
		return c.Unauthorized(err.Error())
	}

	_user, err := client.User.Create().Save(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	err = rdb.Set(ctx, body.Challenge, _user.ID, ExpireTimeToken).Err()
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(true)
}
