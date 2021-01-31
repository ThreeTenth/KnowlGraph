package main

func index(c *Context) error {
	return c.OkHTML(TplIndexHTML, config.Gci)
}
