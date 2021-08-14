package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/asset"
)

// QueryNoPrivateArticlesSQL is query no private article
const QueryNoPrivateArticlesSQL = `select 
"tmp".id as version_id,
"tmp".name as version_name,
"tmp".comment as version_comment,
"tmp".title as version_title,
"tmp".gist as version_gist,
"tmp".state as version_state,
"tmp".lang as version_lang,
"tmp".created_at as version_createdAt,
"articles".id as article_id,
"articles".status as article_status
from (
	select "id", "name", "comment", "title", "gist", "state", "lang", "created_at", "article_versions", row_number() 
	over (
		partition by article_versions order by created_at desc)
	as rownum 
	from  "versions"
	WHERE "state" = 'release' %v
) tmp
inner join "articles"
	on "articles".id = "tmp".article_versions and "articles".status = 'public'
where "tmp".rownum < 2  order by "tmp".created_at desc offset $1 limit $2`

func queryNoPrivateArticles(db *sql.DB, lang string, offset, limit int) ([]*ent.Version, error) {
	_sql := QueryNoPrivateArticlesSQL
	_args := []interface{}{
		offset, limit,
	}
	if lang != "" {
		_sql = fmt.Sprintf(_sql, `and "versions"."lang" = $3`)
		_args = append(_args, lang)
	} else {
		_sql = fmt.Sprintf(_sql, "")
	}
	if config.Debug {
		log.Printf("sql:Query: query=%v args=%v", _sql, _args)
	}
	rows, err := db.Query(_sql, _args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleRows(rows)
}

// QueryUserAssetsSQL is query user article assets
const QueryUserAssetsSQL = `select 
"tmp".id as version_id,
"tmp".name as version_name,
"tmp".comment as version_comment,
"tmp".title as version_title,
"tmp".gist as version_gist,
"tmp".state as version_state,
"tmp".lang as version_lang,
"tmp".created_at as version_createdAt,
"articles".id as article_id,
"articles".status as article_status
from (
	SELECT DISTINCT "versions"."id" as "id", "name", "comment", "title", "gist", "state", "lang", "versions"."created_at" as "created_at", "article_versions", "assets"."created_at" as "joined_at", 
	row_number() over (
			partition by "versions".article_versions order by "versions"."created_at" desc
	) as rownum 
	FROM "versions" 
	INNER JOIN "assets" ON "assets"."owner_id" = $1 and "assets"."status" = $2 and "versions"."article_versions" = "assets"."article_assets"
	WHERE "versions"."state" = 'release' %v
) tmp
inner join "articles"
	on "articles".id = "tmp".article_versions
where tmp.rownum < 2 order by tmp.created_at desc offset $3 limit $4;`

func queryUserAssets(db *sql.DB, userID int, status asset.Status, lang string, offset, limit int) ([]*ent.Version, error) {
	_sql := QueryUserAssetsSQL
	_args := []interface{}{
		userID, status, offset, limit,
	}
	if lang != "" {
		_sql = fmt.Sprintf(_sql, `and "versions"."lang" = $5`)
		_args = append(_args, lang)
	} else {
		_sql = fmt.Sprintf(_sql, "")
	}
	if config.Debug {
		log.Printf("sql:Query: query=%v args=%v", _sql, _args)
	}

	rows, err := db.Query(_sql, _args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleRows(rows)
}

const queryReadListSQL = `SELECT *
FROM "versions" 
INNER JOIN "assets" ON "assets"."owner_id" = 51539607552 and "assets"."status" = 'browse' and "versions"."id" = "assets"."asset_version"
inner join "articles" on "articles".id = "versions".article_versions
WHERE "versions"."state" = 'release'
order by versions.created_at desc offset 0 limit 10;`

func queryReadList(db *sql.DB, userID int, lang string, offset, limit int) ([]*ent.Version, error) {
	_sql := QueryUserAssetsSQL
	_args := []interface{}{
		userID, offset, limit,
	}
	if lang != "" {
		_sql = fmt.Sprintf(_sql, `and "versions"."lang" = $5`)
		_args = append(_args, lang)
	} else {
		_sql = fmt.Sprintf(_sql, "")
	}
	if config.Debug {
		log.Printf("sql:Query: query=%v args=%v", _sql, _args)
	}

	rows, err := db.Query(_sql, _args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleRows(rows)
}

// QueryNodeArticles is query node articles
const QueryNodeArticles = `select 
"tmp".id as version_id,
"tmp".name as version_name,
"tmp".comment as version_comment,
"tmp".title as version_title,
"tmp".gist as version_gist,
"tmp".state as version_state,
"tmp".lang as version_lang,
"tmp".created_at as version_createdAt,
"articles".id as article_id,
"articles".status as article_status
from (
	SELECT DISTINCT "versions"."id" as "id", "name", "comment", "title", "gist", "state", "lang", "versions"."created_at" as "created_at", "article_versions", 
	row_number() over (
			partition by "versions".article_versions order by "versions"."created_at" desc
	) as rownum 
	FROM "versions" 
	INNER JOIN "node_articles" ON "node_articles"."node_id" = $1 and "versions"."article_versions" = "node_articles"."article_id"
	WHERE "versions"."state" = 'release' %v
) tmp
inner join "articles"
	on "articles".id = "tmp".article_versions and "articles"."status" = 'public'
where tmp.rownum < 2 order by tmp.created_at desc offset $2 limit $3`

func queryNodeArticles(db *sql.DB, nodeID int, lang string, offset, limit int) ([]*ent.Version, error) {
	_sql := QueryUserAssetsSQL
	_args := []interface{}{
		nodeID, offset, limit,
	}
	if lang != "" {
		_sql = fmt.Sprintf(_sql, `and "versions"."lang" = $4`)
		_args = append(_args, lang)
	} else {
		_sql = fmt.Sprintf(_sql, "")
	}
	if config.Debug {
		log.Printf("sql:Query: query=%v args=%v", _sql, _args)
	}

	rows, err := db.Query(_sql, _args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleRows(rows)
}

// QueryKeywordArticles is query keyword articles
const QueryKeywordArticles = `select 
"tmp".id as version_id,
"tmp".name as version_name,
"tmp".comment as version_comment,
"tmp".title as version_title,
"tmp".gist as version_gist,
"tmp".state as version_state,
"tmp".lang as version_lang,
"tmp".created_at as version_createdAt,
"articles".id as article_id,
"articles".status as article_status
from (
	SELECT DISTINCT "versions"."id" as "id", "name", "comment", "title", "gist", "state", "lang", "versions"."created_at" as "created_at", "article_versions", 
	row_number() over (
			partition by "versions".article_versions order by "versions"."created_at" desc
	) as rownum 
	FROM "versions" 
	INNER JOIN "version_keywords" ON "version_keywords"."word_id" = $1 and "versions"."id" = "version_keywords"."version_id"
	WHERE "versions"."state" = 'release' %v
) tmp
inner join "articles"
	on "articles".id = "tmp".article_versions and "articles"."status" = 'public'
where tmp.rownum < 2 order by tmp.created_at desc offset $2 limit $3;`

func queryKeywordArticles(db *sql.DB, wordID int, lang string, offset, limit int) ([]*ent.Version, error) {
	_sql := QueryUserAssetsSQL
	_args := []interface{}{
		wordID, offset, limit,
	}
	if lang != "" {
		_sql = fmt.Sprintf(_sql, `and "versions"."lang" = $4`)
		_args = append(_args, lang)
	} else {
		_sql = fmt.Sprintf(_sql, "")
	}
	if config.Debug {
		log.Printf("sql:Query: query=%v args=%v", _sql, _args)
	}

	rows, err := db.Query(_sql, _args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticleRows(rows)
}

func scanArticleRows(rows *sql.Rows) ([]*ent.Version, error) {
	var _versions []*ent.Version
	var err error
	for rows.Next() {
		var _version ent.Version
		var _article ent.Article

		err = rows.Scan(
			&_version.ID,
			&_version.Name,
			&_version.Comment,
			&_version.Title,
			&_version.Gist,
			&_version.State,
			&_version.Lang,
			&_version.CreatedAt,
			&_article.ID,
			&_article.Status)
		if err != nil {
			return nil, err
		}
		_version.Edges.Article = &_article
		_versions = append(_versions, &_version)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if 0 == len(_versions) {
		return nil, errors.New("No more")
	}

	return _versions, nil
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
