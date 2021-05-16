package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/word"
)

func getKeywords(c *Context) error {
	_status := c.Query("status")
	_userID, ok := c.Get(GinKeyUserID)
	if _status == word.StatusPrivate.String() && !ok {
		return c.Unauthorized("Unauthorized")
	}

	_publicCreate := client.Word.Query().Where(word.StatusEQ(word.StatusPublic))

	var _privateCreate *ent.WordQuery

	if ok {
		_privateCreate = client.User.Query().
			Where(user.IDEQ(_userID.(int))).
			QueryWords().
			QueryWord()
	}

	var _words []*ent.Word
	var err error
	switch _status {
	case word.StatusPrivate.String():
		_words, err = _privateCreate.All(ctx)
	case word.StatusPublic.String():
		_words, err = _publicCreate.All(ctx)
	default:
		if _words, err = _publicCreate.All(ctx); err == nil && ok {
			if _privateWords, err1 := _privateCreate.All(ctx); err1 == nil {
				_words = mergeArrAndUnique(_words, _privateWords)
			}
		}
	}

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_words)
}

func mergeArrAndUnique(a []*ent.Word, b []*ent.Word) []*ent.Word {
	check := make(map[int]*ent.Word)
	d := append(a, b...)
	res := make([]*ent.Word, 0)
	for _, val := range d {
		check[val.ID] = val
	}

	for letter := range check {
		res = append(res, check[letter])
	}

	return res
}
