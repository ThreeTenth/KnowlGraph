package main

import (
	"knowlgraph.com/ent"
	"knowlgraph.com/ent/draft"
	"knowlgraph.com/ent/user"
	"knowlgraph.com/ent/version"
)

// 编辑已发布文章的流程
// 参数 id 是文章版本的 ID
// 0. 查询版本关联的文章和内容
// 1. 创建一个指向该文章的草稿
// 2. 复制该文章版本的内容
// 3. 关联复制版到草稿
// 4. 关联草稿到用户
// 5. 返回草稿
func editArticle(c *Context) error {
	var _query struct {
		ID int `form:"id" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	_version, err := client.Version.Query().
		Where(version.ID(_query.ID)).
		WithArticle().
		WithContent().
		Only(ctx)

	if err != nil {
		return err
	}

	_draft, err := client.User.Query().
		Where(user.ID(_userID.(int))).
		QueryDrafts().
		Where(draft.And(
			draft.HasOriginalWith(version.ID(_version.ID)),
			draft.StateEQ(draft.StateWrite))).
		Limit(1).
		WithArticle().
		WithOriginal(func(vq *ent.VersionQuery) {
			vq.WithContent()
		}).
		WithSnapshots().
		First(ctx)

	if err == nil {
		_draft.Edges.Snapshots = make([]*ent.Content, 0)
		return c.Ok(_draft)
	}

	err = WithTx(ctx, client, func(tx *ent.Tx) error {
		_article, err := tx.Article.Create().
			SetStatus(_version.Edges.Article.Status).
			Save(ctx)

		if err != nil {
			return err
		}

		_draft, err := tx.Draft.Create().
			SetArticle(_article).
			SetOriginal(_version).
			SetUserID(_userID.(int)).
			Save(ctx)

		return c.Ok(&_draft)
	})

	return err
}
