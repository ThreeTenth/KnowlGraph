package main

import "net/http"

func articleHTML(c *Context) (int, string, interface{}) {
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

	_userID := c.GetInt(GinKeyUserID)

	_article, status, err := GetArticle(_userID, _query.ID, 0)

	if err != nil {
		return status, TplIndexHTML, _data
	}

	if _userID > 0 {
		_user, _ := client.User.Get(ctx, _userID)

		_data.User = _user
		_data.Logined = true
	}

	_data.Article = _article.Edges.Versions[0]

	return http.StatusOK, TplIndexHTML, _data
}
