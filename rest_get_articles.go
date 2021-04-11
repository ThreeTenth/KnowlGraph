package main

import (
	"knowlgraph.com/ent"
)

func getArticles(c *Context) error {
	var _query struct {
		Limit  int `form:"limit,default=10"`
		Offset int `form:"Offset,default=0"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	rows, err := queryNoPrivateArticles(db, _query.Offset, _query.Limit)
	if err != nil {
		return c.NotFound(err.Error())
	}
	defer rows.Close()

	var _versions []*ent.Version
	for rows.Next() {
		var _version ent.Version
		var _article ent.Article

		err := rows.Scan(
			&_version.ID,
			&_version.Name,
			&_version.Comment,
			&_version.Title,
			&_version.Gist,
			&_version.State,
			&_version.CreatedAt,
			&_article.ID,
			&_article.Status)
		if err != nil {
			return c.InternalServerError("rows scan failed: " + err.Error())
		}
		_version.Edges.Article = &_article
		_versions = append(_versions, &_version)
	}

	if err = rows.Err(); err != nil {
		return c.InternalServerError("rows error: " + err.Error())
	}

	return c.Ok(_versions)
}
