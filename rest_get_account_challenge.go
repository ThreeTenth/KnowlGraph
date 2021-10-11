package main

import (
	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
	"github.com/pkg/errors"
)

// import (
// 	"errors"
// 	"net"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"knowlgraph.com/ent"
// 	"knowlgraph.com/ent/terminal"
// 	"knowlgraph.com/ent/user"

// 	ua "github.com/mileusna/useragent"
// )

// // 知识图谱的用户认证授权
// // A: 未授权设备
// // B: 已授权设备
// //
// // 0: A 创建一个新账号
// // 1: A 向 B 展示请求码
// // 2: B 向 A 展示授权码
// // 3: A(web) 向 A(app) 请求授权
// //
// // # A 向 B 展示请求码
// // 1. A 向服务端请求请求码，并轮询请求码状态
// // 2. B 向 A 扫描请求码
// // 3. B 向服务端推送 A 的请求码，并更新 A 的请求码状态
// // 4. B 向服务端确认授权请求
// // 5. 服务端将 A 绑定到 B 绑定的账号
// // 5 -> 1.1. 授权成功
// //
// // # B 向 A 展示授权码
// // 1. B 向服务端请求授权码，并轮询授权码状态
// // 2. A 向 B 扫描授权码
// // 3. A 向服务端推送 B 的授权码
// // 4. 服务端更新 B 的授权码
// // 4 -> 1.1. B 向服务端确认授权请求
// // 5. 授权成功
// //
// // # A(web) 向 A(app) 请求授权
// // 同 "A 向 B 展示请求码"
// //
// // # 代码
// // - 请求码、授权码、创建码
// // - 请求码等待授权
// // - 授权码等待扫描
// // - 创建码等待推送
// //
// // # 状态
// // - 等待授权: 等待已授权的设备认证
// // - 等待扫描: 等待未授权的设备提交
// // - 等待创建: 等待未授权的设备创建一个账号
// // - 已认证: 该代码已经过已授权的设备认证
// // - 已验证: 该代码已经过未授权的设备提交
// // - 已创建: 该代码已经过未授权的设备确认创建

// // Terminal is the carrier of the account,
// // `Name` is terminal name,
// // `State` is terminal state,
// // `Main` is token authorized to the terminal.
// // If `Main` is empty, it means that the terminal is the first terminal of the account.
// type Terminal struct {
// 	Name  string `json:"name"`
// 	State int    `json:"state"`
// 	Main  string `json:"main,omitempty"`
// }

// func getAccountTerminals(c *Context) error {
// 	_userID, _ := c.Get(GinKeyUserID)
// 	ts, err := client.Terminal.Query().
// 		Where(
// 			terminal.HasUserWith(user.ID(_userID.(int))),
// 		).
// 		All(ctx)

// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	if 0 == len(ts) {
// 		return c.NotFound("No terminals")
// 	}

// 	return c.Ok(&ts)
// }

// func checkAccountChallenge(c *Context) error {
// 	challenge := c.Query("challenge")

// 	var terminal Terminal
// 	if err := GetV4Redis(challenge, &terminal); err != nil {
// 		return c.BadRequest(err.Error())
// 	}

// 	// 如果 state 的值是 user id，则表示该 token 已授权
// 	if TokenStateActivated < terminal.State {
// 		terminal.State = TokenStateAuthorized
// 	}

// 	// Main 值不可向客户端暴露，
// 	// 因为 Main 是授权方的 token。
// 	return c.Ok(Terminal{
// 		Name:  terminal.Name,
// 		State: terminal.State,
// 	})
// }

// func putAccountTerminal(c *Context) error {
// 	challenge, _ := c.Get(GinKeyChallenge)

// 	name, err := getTerminalName(c.Context)
// 	if err != nil {
// 		return c.BadRequest(err.Error())
// 	}

// 	var terminal Terminal

// 	if err := GetV4Redis(challenge.(string), &terminal); err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	terminal.Name = name
// 	terminal.State = TokenStateActivated

// 	err = SetV2Redis(challenge.(string), &terminal, ExpireTimeChallengeConfirm)
// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	// 每一秒查询一次终端授权状态，与缓存时间一致。
// 	// 由于查询需要时间，因此第 60 次查询时，缓存一定会超时，将返回 401。
// 	for i := time.Second; i <= ExpireTimeChallengeConfirm; i += time.Second {
// 		time.Sleep(time.Second)

// 		if rdb.Exists(ctx, challenge.(string)).Val() == 0 {
// 			// 表示授权方取消授权，删除终端 token 了。
// 			return c.Unauthorized("Authorization failed")
// 		} else if err = GetV4Redis(challenge.(string), &terminal); err != nil {
// 			// 表示授权方已授权，终端 token 指向的是 user id，无法解析为 Terminal 结构了。
// 			break
// 		}
// 	}

// 	return c.Ok(true)
// }

// func postAccountAuthn(c *Context) error {
// 	_userID, _ := c.Get(GinKeyUserID)
// 	challenge, _ := c.Get(GinKeyChallenge)

// 	var err error
// 	if c.Request.Method == http.MethodPost {
// 		err = rdb.Set(ctx, challenge.(string), _userID.(int), ExpireTimeToken).Err()
// 	} else if c.Request.Method == http.MethodPatch {
// 		err = rdb.Set(ctx, challenge.(string), _userID.(int), ExpireTimeTokenOnce).Err()
// 	} else if c.Request.Method == http.MethodDelete {
// 		err = rdb.Del(ctx, challenge.(string)).Err()
// 	}

// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	return c.Ok(true)
// }

// func getAccountChallenge(c *Context) error {
// 	challenge := New64BitID()
// 	terminal := Terminal{
// 		State: TokenStateIdle,
// 	}
// 	if _, ok := c.Get(GinKeyUserID); ok {
// 		terminal.Main = getRequestToken(c.Context)
// 	}

// 	if err := SetV2Redis(challenge, &terminal, ExpireTimeChallenge); err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	return c.Ok(challenge)
// }

// func putAccountCreate(c *Context) error {
// 	challenge, _ := c.Get(GinKeyChallenge)

// 	ua := ua.Parse(c.GetHeader(HeaderUserAgent))
// 	name := ua.Device
// 	if ua.Bot {
// 		return c.MethodNotAllowed("Bot can't allow this method")
// 	}

// 	if 0 == len(name) {
// 		if ua.Mobile {
// 			name = "Mobile"
// 		} else if ua.Tablet {
// 			name = "Tablet"
// 		} else if ua.Desktop {
// 			name = "Desktop"
// 		}
// 	}

// 	name = name + " 1"

// 	var terminalID int
// 	err := WithTx(ctx, client, func(tx *ent.Tx) error {
// 		_user, err := tx.User.Create().Save(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		_terminal, err := tx.Terminal.Create().
// 			SetCode(New16bitID()).
// 			SetName(name).
// 			SetUa(c.GetHeader(HeaderUserAgent)).
// 			SetUser(_user).
// 			Save(ctx)

// 		if err != nil {
// 			return err
// 		}

// 		terminalID = _terminal.ID

// 		return rdb.Set(ctx, challenge.(string), _user.ID, ExpireTimeToken).Err()
// 	})

// 	if err != nil {
// 		return c.InternalServerError(err.Error())
// 	}

// 	return c.Ok(terminalID)
// }

// // 检测 challenge 是否存在，challenge 的状态是否合法，以及是否存在非法请求者
// func checkChallenge(state int) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var body struct {
// 			Challenge string `json:"challenge"`
// 		}

// 		if err := c.ShouldBindJSON(&body); err != nil {
// 			c.AbortWithError(http.StatusBadRequest, err)
// 			return
// 		}

// 		var terminal Terminal
// 		if err := GetV4Redis(body.Challenge, &terminal); err != nil {
// 			c.AbortWithError(http.StatusBadRequest, err)
// 			return
// 		} else if state&terminal.State != terminal.State {
// 			c.AbortWithStatus(http.StatusConflict)
// 			return
// 		} else if _, ok := c.Get(GinKeyUserID); ok && getRequestToken(c) != terminal.Main {
// 			// 如果请求 token 和终端授权者的 token 不一致，
// 			// 则表示请求者和授权者不是同一人
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 			return
// 		}

// 		c.Set(GinKeyChallenge, body.Challenge)
// 		c.Next()
// 	}
// }

func getTerminalName(c *gin.Context) (string, error) {
	ua := ua.Parse(c.GetHeader("User-Agent"))
	if 0 == len(ua.OS) {
		return "Unknown", errors.New("Invalid User-Agent")
	}

	name := ua.Device

	if 0 == len(name) {
		if ua.Mobile {
			name = "Mobile"
		} else if ua.Tablet {
			name = "Tablet"
		} else if ua.Desktop {
			name = "Desktop"
		} else {
			name = "Unknown"
		}

		name += " 1"
	}

	return name, nil
}
