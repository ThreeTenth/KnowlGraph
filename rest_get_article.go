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

	_userID, ok := c.Get(GinKeyUserID)

	_article, err := GetArticle(
		ok,
		_userID,
		_query.ID,
		_query.NeedVersions,
		_query.VersionID)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.Ok(&_article)
}

// GetArticle is get article
func GetArticle(isLogin bool, _userID interface{}, articleID int, needVersions bool, versionID int) (*ent.Article, error) {
	_article, err := client.Article.
		Query().
		Where(article.ID(articleID)).
		First(ctx)

	if err != nil {
		return nil, err
	}

	if _article.Status == article.StatusPrivate {
		if !isLogin {
			return nil, errors.New("Unauthorized")
		}

		ok, _ := _article.QueryAssets().Where(asset.HasUserWith(user.ID(_userID.(int)))).Exist(ctx)
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

	if isLogin {
		_assets, _ := _article.QueryAssets().
			Where(asset.
				HasUserWith(user.ID(_userID.(int)))).
			WithArchives(func(aq *ent.ArchiveQuery) {
				aq.WithNode(func(nq *ent.NodeQuery) {
					nq.WithWord().
						WithPath(func(nq *ent.NodeQuery) {
							nq.WithWord()
						})
				})
			}).
			All(ctx)

		_article.Edges.Assets = _assets
	}

	_article.Edges.Versions = _versions
	_article.Edges.Reactions = _reactions

	return _article, err
}
