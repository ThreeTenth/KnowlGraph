package main

func isAuthorized(c *Context) error {
	return c.Ok(true)
}
