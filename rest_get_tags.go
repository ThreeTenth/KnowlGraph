package main

import (
	"knowlgraph.com/ent/tag"
	"knowlgraph.com/ent/user"
)

func getTags(c *Context) error {
	_status := c.Query("status")
	_userID, ok := c.Get(GinKeyUserID)
	if _status == tag.StatusPrivate.String() && !ok {
		return c.Unauthorized("Unauthorized")
	}

	_publicTagCreate := client.Tag.Query().Where(tag.StatusEQ(tag.StatusPublic)).Select(tag.FieldName)
	_privateTagCreate := client.User.Query().
		Where(user.IDEQ(_userID.(int))).
		QueryTags().
		Select(tag.FieldName)

	var _tags []string
	var err error
	switch _status {
	case tag.StatusPrivate.String():
		_tags, err = _privateTagCreate.Strings(ctx)
	case tag.StatusPublic.String():
		_tags, err = _publicTagCreate.Strings(ctx)
	default:
		if _tags, err = _privateTagCreate.Strings(ctx); err == nil {
			if _publicTags, err1 := _publicTagCreate.Strings(ctx); err1 == nil {
				_tags = mergeArrAndUnique(_tags, _publicTags)
			}
		}
	}

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_tags)
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
