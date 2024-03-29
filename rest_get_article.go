package main

import (
	"net/http"

	"github.com/pkg/errors"
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/archive"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/node"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

func getArticle(c *Context) error {
	var _query struct {
		ID        int    `form:"id" binding:"required"`
		Lang      string `form:"lang"`
		VersionID int    `form:"versionID"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID := c.GetInt(GinKeyUserID)

	_article, status, err := GetArticle(
		_userID,
		_query.ID,
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
func GetArticle(_userID int, articleID int, versionID int) (*ent.Article, int, error) {
	_article, err := client.Article.
		Query().
		Where(article.ID(articleID)).
		First(ctx)

	if err != nil {
		return nil, http.StatusNotFound, err
	}

	isLogin := 0 < _userID

	if _article.Status == article.StatusPrivate {
		if !isLogin {
			return nil, http.StatusUnauthorized, errors.New("Unauthorized")
		}

		ok, _ := _article.QueryAssets().Where(asset.HasUserWith(user.ID(_userID))).Exist(ctx)
		if !ok {
			return nil, http.StatusUnauthorized, errors.New("Unauthorized")
		}
	}

	vq := _article.QueryVersions()

	if 0 < versionID {
		vq.Where(version.ID(versionID))
	} else {
		vq.Order(ent.Desc(version.FieldCreatedAt))
	}
	_version, err := vq.WithContent().WithKeywords().First(ctx)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	_reactions, _ := _article.QueryReactions().All(ctx)

	_nodesWhere := node.StatusEQ(node.StatusPublic)
	if isLogin {
		// 已登录用户的私有节点查询
		_nodesWhere = node.Or(
			_nodesWhere,
			node.HasArchivesWith(
				archive.HasUserWith(
					user.ID(_userID))))
	}
	_nodes, _ := _article.
		QueryNodes().
		Where(_nodesWhere).
		WithWord().
		WithPath(func(nq *ent.NodeQuery) {
			nq.WithWord()
		}).
		All(ctx)

	if isLogin {
		_assets, _ := _article.
			QueryAssets().
			Where(
				asset.HasUserWith(user.ID(_userID)),
				asset.StatusNEQ(asset.StatusBrowse),
			).
			All(ctx)

		_article.Edges.Assets = _assets

		_archives, _ := _article.
			QueryArchives().
			Where(archive.HasUserWith(user.ID(_userID))).
			WithNode(func(nq *ent.NodeQuery) {
				nq.WithWord()
			}).
			All(ctx)

		_article.Edges.Archives = _archives
	}

	_article.Edges.Versions = []*ent.Version{_version}
	_article.Edges.Reactions = _reactions
	_article.Edges.Nodes = _nodes

	return _article, http.StatusOK, err
}
