package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/language"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getUserArticles(c *Context) error {
	var _query struct {
		Status asset.Status `from:"status"`
		Lang   string       `from:"lang"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_lang := getLanguage(_query.Lang)

	_articles, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryAssets().
		Where(asset.StatusEQ(_query.Status)).
		WithArticle(func(aq *ent.ArticleQuery) {
			aq.WithVersions(func(vq *ent.VersionQuery) {
				if _lang != nil {
					vq = vq.Where(version.HasLangWith(language.Code(_lang.Code)))
				}
				vq.Order(ent.Desc(version.FieldCreatedAt)).
					Limit(1).
					WithContent()
			})
		}).
		All(ctx)

	if err != nil {
		return c.BadRequest(err.Error())
	}

	return c.Ok(_articles)
}
