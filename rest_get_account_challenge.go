package main

func getAccountChallenge(c *Context) error {
	challenge := RandNdigMbitString(26, 62)
	err := rdb.Set(ctx, challenge, true, ExpireTimeChallenge).Err()

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(challenge)
}
