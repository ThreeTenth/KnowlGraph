package main

import "net/http"

func articleHTML(c *Context) (int, string, interface{}) {
	_userID, ok := c.Get(GinKeyUserID)

	var _data struct {
		Logined bool
		User    interface{}
		Article interface{}
	}

	var _query struct {
		ID int `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&_query); err != nil {
		return http.StatusNotFound, TplIndexHTML, _data
	}

	if ok {
		_user, err := client.User.Get(ctx, _userID.(int))

		if err != nil {
			return http.StatusOK, TplIndexHTML, _data
		}

		_article, err := GetArticle(_user.ID, _query.ID, false, 0)

		if err != nil {
			return http.StatusNotFound, TplIndexHTML, _data
		}

		_data.Logined = true
		_data.User = _user
		_data.Article = _article.Edges.Versions[0]

		return http.StatusOK, TplIndexHTML, _data
	}

	return http.StatusOK, TplIndexHTML, _data
}
