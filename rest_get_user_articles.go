package main

import (
	"time"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/user"
)

func getUserArticles(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)

	var _versions []*struct {
		ID        int       `json:"id"`
		Article   int       `json:"article_versions"`
		Seo       string    `json:"seo"`
		CreatedAt time.Time `json:"created_at"`
	}

	_query := client.User.Query().
		Where(user.IDEQ(_userID.(int))).
		QueryArticles().
		Where(article.HasVersions()).
		QueryVersions()

	if _lang := c.Query(QueryLang); _lang != "" {
		_query.Where(content.HasLangWith(language.IDEQ(_lang)))
	}

	err := _query.Order(ent.Desc(content.FieldID)).
		GroupBy(content.ArticleColumn).
		Aggregate(
			ent.As(ent.Max(content.FieldID), content.FieldID)).
		Scan(ctx, &_versions)
	if err != nil {
		return c.InternalServerError(err.Error())
	}
	if len(_versions) == 0 {
		return c.Ok(_versions)
	}

	_ids := make([]int, len(_versions))

	for i, _version := range _versions {
		_ids[i] = _version.ID
	}

	_contents, err := client.Content.Query().Where(content.IDIn(_ids...)).All(ctx)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	for _, _content := range _contents {
		for _, _version := range _versions {
			if _content.ID == _version.ID {
				_version.Seo = _content.Seo
				_version.CreatedAt = _content.CreatedAt
			}
		}
	}

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_versions)
}
