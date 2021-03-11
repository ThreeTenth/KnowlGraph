package main

import "net/http"

func index(c *Context) (int, string, interface{}) {
	_userID, ok := c.Get(GinKeyUserID)

	var _data struct {
		Logined bool
		User    interface{}
	}

	if ok {
		_user, err := client.User.Get(ctx, _userID.(int))

		if err != nil {
			_data.Logined = false
			return http.StatusOK, TplIndexHTML, _data
		}

		_data.Logined = true
		_data.User = _user

		return http.StatusOK, TplIndexHTML, _data
	}

	return http.StatusOK, TplIndexHTML, _data
}
