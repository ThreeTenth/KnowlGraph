package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/archive"
	"knowlgraph.com/ent/article"
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/node"
	"knowlgraph.com/ent/user"
)

func putNodeArticle(c *Context) error {
	var _query struct {
		NodeID    int `json:"nodeId" binding:"required"`
		ArticleID int `json:"articleId" binding:"required"`
	}

	err := c.ShouldBindJSON(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	userID := c.GetInt(GinKeyUserID)

	ok, err := client.Node.
		Query().
		Where(
			node.ID(_query.NodeID),
			node.Or(
				node.StatusEQ(node.StatusPublic),
				node.HasArchivesWith(archive.HasUserWith(user.ID(userID))),
			)).
		Exist(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	} else if !ok {
		return c.Unauthorized("No access to node")
	}

	ok, err = client.Article.
		Query().
		Where(
			article.ID(_query.ArticleID),
			article.Or(
				article.StatusNEQ(article.StatusPrivate),
				article.HasAssetsWith(asset.HasUserWith(user.ID(userID))),
			)).
		Exist(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	} else if !ok {
		return c.Unauthorized("No access to node")
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		_, err = tx.Node.
			Update().
			Where(node.ID(_query.NodeID)).
			RemoveArticleIDs(_query.ArticleID).
			AddArticleIDs(_query.ArticleID).
			AddWeight(1).
			Save(ctx)

		if err != nil {
			return err
		}

		_archive, err := tx.Archive.Query().Where(
			archive.StatusEQ(archive.StatusStar),
			archive.HasUserWith(user.ID(userID)),
			archive.HasNodeWith(node.ID(_query.NodeID)),
		).Only(ctx)

		if err != nil && !ent.IsNotFound(err) {
			return err
		}

		if _archive == nil {
			_archive, err = tx.Archive.
				Create().
				SetStatus(archive.StatusStar).
				SetNodeID(_query.NodeID).
				SetUserID(userID).
				Save(ctx)

			if err != nil {
				return err
			}
		}

		err = _archive.
			Update().
			RemoveArticleIDs(_query.ArticleID).
			AddArticleIDs(_query.ArticleID).
			Exec(ctx)

		if err != nil {
			return err
		}

		ok, err := tx.Asset.Query().Where(
			asset.StatusEQ(asset.StatusStar),
			asset.HasUserWith(user.ID(userID)),
			asset.HasArticleWith(article.ID(_query.ArticleID)),
		).Exist(ctx)

		if err != nil {
			return err
		}

		if !ok {
			err = tx.Asset.
				Create().
				SetStatus(asset.StatusStar).
				SetArticleID(_query.ArticleID).
				SetUserID(userID).
				Exec(ctx)
		}

		return err
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.NoContent()
}
