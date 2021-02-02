package main

import "net/http"

func index(c *Context) (int, string, interface{}) {
	_userID, ok := c.Get(GinKeyUserID)
	if ok {
		_user, err := client.User.Get(ctx, _userID.(int))
		if err != nil {
			return http.StatusOK, TplIndexHTML, nil
		}

		return http.StatusOK, TplMainHTML, _user
	}

	return http.StatusOK, TplIndexHTML, nil
}
