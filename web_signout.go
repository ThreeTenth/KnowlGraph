package main

func signout(c *Context) error {
	return c.TemporaryRedirect("/")
}
