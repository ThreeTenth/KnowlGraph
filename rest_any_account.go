package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"knowlgraph.com/ent"
)

// Terminal is user's device Terminal
type Terminal struct {
	Name        string
	UA          string
	ClientState string // 一个不可猜测的随机字符串。它用于防止跨站点请求伪造攻击。
	State       int
	UserID      int
	OnlyOnce    bool // 是否仅此一次授权
}

func (t *Terminal) authorized() bool {
	return t.State == TokenStateAuthorized
}

func getChallenge(c *Context) error {
	state := c.Query("state")
	challenge := New16BitID()

	t := Terminal{
		State:       TokenStateIdle,
		ClientState: state,
	}
	if userID, ok := c.Get(GinKeyUserID); ok {
		t.UserID = userID.(int)
	} else {
		t.UA = c.GetHeader(HeaderUserAgent)
		t.Name, _ = getTerminalName(c.Context)
	}

	if err := SetV2Redis(RChallenge(challenge), &t, ExpireTimeChallenge); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(challenge)
}

// 如果已有账号，则该凭据关联账号，若无，则创建一个新账号
func makeCredential(c *Context) error {
	var form struct {
		State     string `form:"state" binding:"required"`
		Challenge string `form:"challenge" binding:"required"`
	}

	if err := c.ShouldBindQuery(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	var t Terminal
	if err := GetV4Redis(RChallenge(form.Challenge), &t); err != nil {
		return c.BadRequest(err.Error())
	}
	if form.State != t.ClientState {
		return c.BadRequest("Bad state")
	}
	if 0 != t.UserID && !t.authorized() {
		// 如果已绑定认证账号，但认证账号未授权，则返回“没有权限”

		return c.Unauthorized("Permission denied")
	}

	id := 0
	token := New64BitID()
	d := ExpireTimeToken
	if t.OnlyOnce {
		d = ExpireTimeTokenOnce
	}

	terminalMap := make(map[int]string)

	err := WithTx(ctx, client, func(tx *ent.Tx) error {
		if 0 == t.UserID {
			user, err1 := tx.User.Create().Save(ctx)
			if err1 != nil {
				return err1
			}
			t.UserID = user.ID
		} else {
			GetV4Redis(RUser(t.UserID), &terminalMap)
		}

		terminal, err1 := tx.Terminal.
			Create().
			SetCode(New16bitID()).SetName(t.Name).SetUa(t.UA).SetUserID(t.UserID).
			Save(ctx)

		if err1 != nil {
			return err1
		}

		id = terminal.ID

		terminalMap[id] = token

		_, err1 = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if err1 = SetV2RedisPipe(pipe, RUser(t.UserID), &terminalMap, d); err1 != nil {
				return err1
			}
			pipe.Set(ctx, RToken(token), t.UserID, d)
			pipe.Del(ctx, RChallenge(form.Challenge))
			return nil
		})

		return err1
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(gin.H{"id": id, "token": token, "onlyOnce": t.OnlyOnce})
}

func activateTerminal(c *Context) error {
	challenge := c.Query("challenge")

	var t Terminal
	if err := GetV4Redis(RChallenge(challenge), &t); err != nil {
		return c.BadRequest(err.Error())
	}
	if userID, ok := c.Get(GinKeyUserID); ok {
		if 0 != t.UserID {
			// 认证用户不得激活已关联认证用户的 challenge
			return c.Unauthorized("Invalid user")
		}
		t.UserID = userID.(int)
	} else if 0 == t.UserID {
		return c.Unauthorized("Invalid user")
	} else {
		t.UA = c.GetHeader(HeaderUserAgent)
		t.Name, _ = getTerminalName(c.Context)
	}

	t.State = TokenStateActivated

	if err := SetV2Redis(RChallenge(challenge), &t, ExpireTimeChallengeConfirm); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(gin.H{"name": t.Name, "state": t.State})
}

func authorizeTerminal(c *Context) error {
	var form struct {
		Challenge string `form:"challenge" binding:"required"`
		OnlyOnce  bool   `form:"onlyOnce"`
	}

	if err := c.ShouldBindQuery(&form); err != nil {
		return c.BadRequest(err.Error())
	}

	var t Terminal
	if err := GetV4Redis(RChallenge(form.Challenge), &t); err != nil {
		return c.BadRequest(err.Error())
	}
	if t.State < TokenStateActivated {
		return c.PreconditionFailed("Please activate the terminal first")
	}
	if userID, _ := c.Get(GinKeyUserID); userID != t.UserID {
		return c.Unauthorized("Invalid user")
	}

	t.OnlyOnce = form.OnlyOnce
	t.State = TokenStateAuthorized

	if err := SetV2Redis(RChallenge(form.Challenge), &t, ExpireTimeChallengeConfirm); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(gin.H{"name": t.Name, "state": t.State})
}

func cancelActivateTerminal(c *Context) error {
	challenge := c.Query("challenge")

	if err := rdb.Del(ctx, RChallenge(challenge)).Err(); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(true)
}

func scanChallenge(c *Context) error {
	challenge := c.Query("challenge")

	var t Terminal
	if err := GetV4Redis(RChallenge(challenge), &t); err != nil {
		return c.BadRequest(err.Error())
	}

	return c.Ok(gin.H{"name": t.Name, "state": t.State})
}
