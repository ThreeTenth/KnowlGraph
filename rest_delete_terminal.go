package main

import (
	"github.com/go-redis/redis/v8"
	"knowlgraph.com/ent/terminal"
	"knowlgraph.com/ent/user"
)

func deleteTerminal(c *Context) error {
	var form struct {
		ID int `form:"id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	userID, _ := c.Get(GinKeyUserID)
	token := getRequestToken(c.Context)

	terminalMap := make(map[int]string)
	if err := GetV4Redis(RUser(userID.(int)), &terminalMap); err != nil {
		return c.ServiceUnavailable(err.Error())
	}

	var currentTerminalID int
	for k, v := range terminalMap {
		if v == token {
			currentTerminalID = k
			break
		}
	}
	if 0 == currentTerminalID {
		return c.Unauthorized("The current terminal is not authorized")
	}

	if currentTerminalID == form.ID {
		// 终端不能在有其他终端存在时，删除自己，所以
		// 当请求删除的终端 ID 与当前操作终端的 ID 相同时，拒绝删除请求，
		// 除非当前账号有且只有一个终端，即当前操作终端时，删除有效。
		count, err := client.Terminal.
			Query().
			Where(terminal.HasUserWith(user.ID(userID.(int)))).
			Count(ctx)

		if err != nil {
			return c.InternalServerError(err.Error())
		}

		if 1 < count {
			return c.PreconditionFailed("Please delete other terminals first")
		}
	}

	deleteTerminalToken := terminalMap[form.ID]
	delete(terminalMap, form.ID)
	d := rdb.TTL(ctx, RUser(userID.(int))).Val()

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if err1 := SetV2RedisPipe(pipe, RUser(userID.(int)), &terminalMap, d); err1 != nil {
			return err1
		}
		pipe.Del(ctx, RToken(deleteTerminalToken))
		return nil
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	count, err := client.Terminal.
		Delete().
		Where(terminal.ID(form.ID), terminal.HasUserWith(user.ID(userID.(int)))).
		Exec(ctx)

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	if 0 == count {
		return c.NotFound("The terminal ID was not found")
	}

	return c.Ok(true)
}
