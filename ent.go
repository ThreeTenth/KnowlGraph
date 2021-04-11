package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/pkg/errors"
	"knowlgraph.com/ent"
)

// QueryNoPrivateArticlesSQL is query no private article
const QueryNoPrivateArticlesSQL = `select 
"tmp".id as version_id,
"tmp".name as version_name,
"tmp".comment as version_comment,
"tmp".title as version_title,
"tmp".gist as version_gist,
"tmp".state as version_state,
"tmp".created_at as version_createdAt,
"articles".id as article_id,
"articles".status as article_status
from (
	select "id", "name", "comment", "title", "gist", "state", "created_at", "article_versions", row_number() 
	over (
		partition by article_versions order by created_at)
	as rownum 
	from  "versions"
	WHERE "state" = 'release'
) tmp
inner join "articles"
	on "articles".id = "tmp".article_versions and "articles".status = 'public'
where "tmp".rownum < 2  order by "tmp".created_at desc offset $1 limit $2`

func queryNoPrivateArticles(db *sql.DB, offset, limit int) (*sql.Rows, error) {
	if config.Debug {
		log.Printf("sql:Query: query=%v args=[%v, %v]", QueryNoPrivateArticlesSQL, offset, limit)
	}
	return db.Query(QueryNoPrivateArticlesSQL, offset, limit)
}

// WithTx best Practices, reusable function that runs callbacks in a transaction
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "committing transaction: %v", err)
	}
	return nil
}
