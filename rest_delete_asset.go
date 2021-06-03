package main

import (
	"knowlgraph.com/ent/asset"
	"knowlgraph.com/ent/user"
)

func deleteAsset(c *Context) error {
	var _query struct {
		AssetID int `form:"assetId" binding:"required"`
	}

	err := c.ShouldBindQuery(&_query)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	_userID, _ := c.Get(GinKeyUserID)

	i, err := client.Asset.
		Delete().
		Where(
			asset.ID(_query.AssetID),
			asset.HasUserWith(user.ID(_userID.(int))),
		).
		Exec(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if 0 == i {
		return c.BadRequest("There are no resources to delete")
	}

	return c.Ok(i)
}
