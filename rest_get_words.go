package main

import (
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/word"
)

func getKeywords(c *Context) error {
	_status := c.Query("status")
	_userID, ok := c.Get(GinKeyUserID)
	if _status == word.StatusPrivate.String() && !ok {
		return c.Unauthorized("Unauthorized")
	}

	_publicCreate := client.Word.Query().Where(word.StatusEQ(word.StatusPublic)).Select(word.FieldName)
	_privateCreate := client.User.Query().
		Where(user.IDEQ(_userID.(int))).
		QueryWords().
		QueryWord().
		Select(word.FieldName)

	var _words []string
	var err error
	switch _status {
	case word.StatusPrivate.String():
		_words, err = _privateCreate.Strings(ctx)
	case word.StatusPublic.String():
		_words, err = _publicCreate.Strings(ctx)
	default:
		if _words, err = _privateCreate.Strings(ctx); err == nil {
			if _publicWords, err1 := _publicCreate.Strings(ctx); err1 == nil {
				_words = mergeArrAndUnique(_words, _publicWords)
			}
		}
	}

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_words)
}

func mergeArrAndUnique(a []string, b []string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}
