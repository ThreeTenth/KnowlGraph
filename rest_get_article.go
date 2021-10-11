package main

import (
	"net/http"

	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/archive"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/node"
	"knowlgraph.com/ent/quote"
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

	_article, status, err := GetArticle(
		ok,
		_userID,
		_query.ID,
		_query.NeedVersions,
		_query.VersionID)

	switch status {
	case http.StatusOK:
		return c.Ok(&_article)
	case http.StatusNotFound:
		return c.NotFound(err.Error())
	case http.StatusUnauthorized:
		return c.Unauthorized(err.Error())
	}

	return c.NotFound(err.Error())
}

// GetArticle is get article
func GetArticle(isLogin bool, _userID interface{}, articleID int, needVersions bool, versionID int) (*ent.Article, int, error) {
	_article, err := client.Article.
		Query().
		Where(article.ID(articleID)).
		First(ctx)

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	if _article.Status == article.StatusPrivate {
		if !isLogin {
			return nil, http.StatusUnauthorized, errors.New("Unauthorized")
		}

		ok, _ := _article.QueryAssets().Where(asset.HasUserWith(user.ID(_userID.(int)))).Exist(ctx)
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("Unauthorized")
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
	_version, err := vq.WithContent().WithKeywords().First(ctx)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	_reactions, _ := _article.QueryReactions().All(ctx)

	_quotes, _ := _article.
		QueryQuotes().
		Where(quote.HasArticleWith(article.StatusEQ(article.StatusPublic))).
		WithArticle().
		All(ctx)

	_nodesWhere := node.StatusEQ(node.StatusPublic)
	if isLogin {
		// 已登录用户的私有节点查询
		_nodesWhere = node.And(
			_nodesWhere,
			node.HasArchivesWith(
				archive.HasUserWith(
					user.ID(_userID.(int)))))
	}
	_nodesQuery := _article.QueryNodes().Where(_nodesWhere).WithWord()

	if isLogin {
		// 已登录用户的归档查询
		_nodesQuery.WithArchives(func(aq *ent.ArchiveQuery) {
			aq.Where(archive.HasUserWith(user.ID(_userID.(int))))
		})
	}

	_nodes, _ := _nodesQuery.All(ctx)

	if isLogin {
		_assets, _ := _article.
			QueryAssets().
			Where(asset.HasUserWith(user.ID(_userID.(int)))).
			All(ctx)

		_article.Edges.Assets = _assets
	}

	_article.Edges.Versions = []*ent.Version{_version}
	_article.Edges.Reactions = _reactions
	_article.Edges.Nodes = _nodes
	_article.Edges.Quotes = _quotes

	return _article, http.StatusOK, err
}
