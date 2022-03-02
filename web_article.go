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

	_article, status, err := GetArticle(ok, _userID.(int), _query.ID, 0)

	if err != nil {
		return status, TplIndexHTML, _data
	}

	if ok {
		_user, _ := client.User.Get(ctx, _userID.(int))

		_data.User = _user
	}

	_data.Logined = ok
	_data.Article = _article.Edges.Versions[0]

	return http.StatusOK, TplIndexHTML, _data
}
