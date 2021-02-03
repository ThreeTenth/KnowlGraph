package main

import (
	"time"

	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/content"
	"knowlgraph.com/ent/language"
)

func getUserArticles(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)
	_user, err := client.User.Get(ctx, _userID.(int))
	if err != nil {
		return c.Forbidden(err.Error())
	}

	var _versions []*struct {
		ID        int       `json:"id"`
		Article   int       `json:"article_versions"`
		Seo       string    `json:"seo"`
		CreatedAt time.Time `json:"created_at"`
	}

	_query := _user.QueryArticles().
		Where(article.HasVersions()).
		QueryVersions()

	if _lang := c.Query(QueryLang); _lang != "" {
		_query.Where(content.HasLangWith(language.IDEQ(_lang)))
	}

	err = _query.Order(ent.Desc(content.FieldID)).
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

	for _, _version := range _versions {
		_content, err := client.Content.Get(ctx, _version.ID)
		if err != nil {
			return c.InternalServerError(err.Error())
		}
		_version.Seo = _content.Seo
		_version.CreatedAt = _content.CreatedAt
	}

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(_versions)
}
