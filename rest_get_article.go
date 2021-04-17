package main

import (
	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getArticle(c *Context) error {
	var _query struct {
		ID           int    `form:"id" binding:"required"`
		Lang         string `form:"lang"`
		NeedVersions bool   `form:"needVersions"`
		VersionID    int    `form:"versionID"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_article, err := GetArticle(
		_userID.(int),
		_query.ID,
		_query.NeedVersions,
		_query.VersionID)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_article)
}

// GetArticle is get article
func GetArticle(_userID int, articleID int, needVersions bool, versionID int) (*ent.Article, error) {
	_article, err := client.Article.
		Query().
		Where(article.ID(articleID)).
		First(ctx)

	if err != nil {
		return nil, err
	}

	if _article.Status == article.StatusPrivate {
		ok, _ := _article.QueryAssets().Where(asset.HasUserWith(user.ID(_userID))).Exist(ctx)
		if !ok {
			return nil, errors.New("Unauthorized")
		}
	}

	vq := _article.QueryVersions()

	if 0 < versionID {
		vq.Where(version.ID(versionID))
	} else if needVersions {
		vq.Order(ent.Desc(version.FieldCreatedAt))
	} else {
		vq.Order(ent.Desc(version.FieldCreatedAt)).Limit(1)
	}

	_versions, err := vq.WithContent().WithKeywords().WithQuotes().All(ctx)
	if err != nil {
		return nil, err
	}

	_reactions, _ := _article.QueryReactions().All(ctx)

	_article.Edges.Versions = _versions
	_article.Edges.Reactions = _reactions

	return _article, err
}
