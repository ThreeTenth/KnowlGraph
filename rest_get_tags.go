package main

import (
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"
)

func getTags(c *Context) error {
	_publicTags, err := client.Article.Query().
		Where(article.StatusEQ(article.StatusPublic)).
		QueryVersions().
		QueryTags().
		Select(tag.FieldName).
		Strings(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	_userID, ok := c.Get(GinKeyUserID)
	if ok {
		_userTags, err := client.User.Query().
			Where(user.IDEQ(_userID.(int))).
			QueryArticles().
			QueryVersions().
			QueryTags().
			Select(tag.FieldName).
			Strings(ctx)

		if err == nil {
			_publicTags = mergeArrAndUnique(_publicTags, _userTags)
		}
	}

	return c.Ok(_publicTags)
}

func mergeArrAndUnique(a []string, b []string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}
