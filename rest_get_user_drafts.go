package main

import (
	"knowlgraph.com/ent/predicate"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getUserDrafts(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	var _query struct {
		Status version.Status `form:"status"`
	}

	c.ShouldBindQuery(_query)

	var _statusWhere predicate.Version
	switch _query.Status {
	case version.StatusNew:
		_statusWhere = version.StatusEQ(_query.Status)
	case version.StatusModify:
		_statusWhere = version.StatusEQ(_query.Status)
	case version.StatusTranslate:
		_statusWhere = version.StatusEQ(_query.Status)
	case "":
		_statusWhere = version.StatusIn(version.StatusNew, version.StatusModify, version.StatusTranslate)
	default:
		return c.BadRequest("status invaild")
	}

	_articles, err := client.Version.Query().
		Where(version.And(
			version.HasUserWith(user.IDEQ(_userID.(int))), _statusWhere)).
		All(ctx)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(_articles)
}
