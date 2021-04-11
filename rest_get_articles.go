package main

import (
	"knowlgraph.com/ent"
)

// QueryNoPrivateArticlesSQL is query no private article
const QueryNoPrivateArticlesSQL = `select * from (
	select  id, name, state, created_at, article_versions, row_number() over (partition by article_versions order by created_at) as rownum from  "versions" WHERE "state" = 'release'
) tmp
where rownum < 2  order by created_at desc offset $1 limit $2 `

func getArticles(c *Context) error {
	var _query struct {
		Limit  int `form:"limit,default=10"`
		Offset int `form:"Offset,default=0"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	rows, err := db.Query(QueryNoPrivateArticlesSQL, _query.Offset, _query.Limit)
	if err != nil {
		return c.NotFound(err.Error())
	}
	defer rows.Close()

	var _versions []*ent.Version
	for rows.Next() {
		var _version ent.Version
		var _article ent.Article
		var rownum int

		err := rows.Scan(&_version.ID, &_version.Name, &_version.State, &_version.CreatedAt, &_article.ID, &rownum)
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
