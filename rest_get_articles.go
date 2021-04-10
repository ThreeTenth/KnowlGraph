package main

import "fmt"

// QueryNoPrivateArticlesSQL is query no private article
const QueryNoPrivateArticlesSQL = `select * from (
	select  id, name, state, created_at, article_versions, row_number() over (partition by article_versions order by created_at) as rownum from  "versions" WHERE "state" = 'release'
) tmp
where rownum < 2  order by created_at desc limit ? offset ?;`

func getArticles(c *Context) error {
	var _query struct {
		Limit  int
		Offset int
	}

	var id int
	rows, err := db.Query(QueryNoPrivateArticlesSQL, _query.Limit, _query.Offset)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return c.Ok(id)
}
