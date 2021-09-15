package main

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"knowlgraph.com/ent"

	ua "github.com/mileusna/useragent"
)

// Terminal is the carrier of the account,
// `Name` is terminal name,
// `State` is terminal state,
// `Main` is token authorized to the terminal.
// If `Main` is empty, it means that the terminal is the first terminal of the account.
type Terminal struct {
	Name  string `json:"name"`
	State int    `json:"state"`
	Main  string `json:"main,omitempty"`
}

func checkAccountChallenge(c *Context) error {
	challenge := c.Query("challenge")

	var terminal Terminal
	if err := GetV4Redis(challenge, &terminal); err != nil {
		return c.BadRequest(err.Error())
	}

	// 如果 state 的值是 user id，则表示该 token 已授权
	if TokenStateActivated < terminal.State {
		terminal.State = TokenStateAuthorized
	}

	return c.Ok(Terminal{
		Name:  terminal.Name,
		State: terminal.State,
	})
}

func putAccountTerminal(c *Context) error {
	challenge, _ := c.Get(GinKeyChallenge)

	name, err := getTerminalName(c.Context)
	if err != nil {
		return c.BadRequest(err.Error())
	}

	var terminal Terminal

	if err := GetV4Redis(challenge.(string), &terminal); err != nil {
		return c.InternalServerError(err.Error())
	}

	terminal.Name = name
	terminal.State = TokenStateActivated

	err = SetV2Redis(challenge.(string), &terminal, ExpireTimeChallengeConfirm)
	if err != nil {
		return c.InternalServerError(err.Error())
	}

	// 每一秒查询一次终端授权状态，与缓存时间一致。
	// 由于查询需要时间，因此第 60 次查询时，缓存一定会超时，将返回 401。
	for i := time.Second; i <= ExpireTimeChallengeConfirm; i += time.Second {
		time.Sleep(time.Second)

		if rdb.Exists(ctx, challenge.(string)).Val() == 0 {
			// 表示授权方取消授权，删除终端 token 了。
			return c.Unauthorized("Authorization failed")
		} else if err = GetV4Redis(challenge.(string), &terminal); err != nil {
			// 表示授权方已授权，终端 token 指向的是 user id，无法解析为 Terminal 结构了。
			break
		}
	}

	return c.Ok(true)
}

func postAccountAuthn(c *Context) error {
	_userID, _ := c.Get(GinKeyUserID)
	challenge, _ := c.Get(GinKeyChallenge)

	var err error
	if c.Request.Method == http.MethodPost {
		err = rdb.Set(ctx, challenge.(string), _userID.(int), ExpireTimeToken).Err()
	} else if c.Request.Method == http.MethodPatch {
		err = rdb.Set(ctx, challenge.(string), _userID.(int), ExpireTimeTokenOnce).Err()
	} else if c.Request.Method == http.MethodDelete {
		err = rdb.Del(ctx, challenge.(string)).Err()
	}

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(true)
}

func getAccountChallenge(c *Context) error {
	challenge := New64BitID()
	terminal := Terminal{
		State: TokenStateIdle,
	}
	if _, ok := c.Get(GinKeyUserID); ok {
		terminal.Main = getRequestToken(c.Context)
	}

	if err := SetV2Redis(challenge, &terminal, ExpireTimeChallenge); err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(challenge)
}

func putAccountCreate(c *Context) error {
	challenge, _ := c.Get(GinKeyChallenge)

	err := WithTx(ctx, client, func(tx *ent.Tx) error {
		_user, err := tx.User.Create().Save(ctx)
		if err != nil {
			return err
		}

		return rdb.Set(ctx, challenge.(string), _user.ID, ExpireTimeToken).Err()
	})

	if err != nil {
		return c.InternalServerError(err.Error())
	}

	return c.Ok(true)
}

// 检测 challenge 是否存在，challenge 的状态是否合法，以及是否存在非法请求者
func checkChallenge(state int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Challenge string `json:"challenge"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var terminal Terminal
		if err := GetV4Redis(body.Challenge, &terminal); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		} else if state&terminal.State != terminal.State {
			c.AbortWithStatus(http.StatusConflict)
			return
		} else if _, ok := c.Get(GinKeyUserID); ok && getRequestToken(c) != terminal.Main {
			// 如果请求 token 和终端授权者的 token 不一致，
			// 则表示请求者和授权者不是同一人
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(GinKeyChallenge, body.Challenge)
		c.Next()
	}
}

func getTerminalName(c *gin.Context) (string, error) {
	ua := ua.Parse(c.GetHeader("User-Agent"))
	if 0 == len(ua.OS) {
		return "", errors.New("Invalid User-Agent")
	}

	name := "(" + ua.OS + " " + ua.OSVersion + ")"
	if 0 == len(ua.Device) {
		name = ua.Name + name
	} else {
		name = ua.Device + name
	}

	ip := c.ClientIP()

	city, err := rdb.Get(ctx, ip).Result()
	if err != nil {
		var ipinfo *IPInfo
		if IsPrivateIP(net.ParseIP(ip)) {
			ipinfo, err = MyIP()
		} else {
			ipinfo, err = ForeignIP(ip)
		}
		if err != nil {
			return name, nil
		}
		city = ipinfo.City
		rdb.Set(ctx, ip, city, ExpireTimeIPInfo)
	}

	name += " in " + city

	return name, nil
}
