package main

import "net/http"

func index(c *Context) (int, string, interface{}) {
	_userID, ok := c.Get(GinKeyUserID)
	if ok {
		_user, err := client.User.Get(ctx, _userID.(int))
		if err != nil {
			return http.StatusOK, TplIndexHTML, nil
		}

		_langs, _ := client.Language.Query().All(ctx)

		_data := struct {
			User      interface{}
			Languages interface{}
		}{
			User:      _user,
			Languages: _langs,
		}

		return http.StatusOK, TplMainHTML, _data
	}

	return http.StatusOK, TplIndexHTML, nil
}
