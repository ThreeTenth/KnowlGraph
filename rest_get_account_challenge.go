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

type terminal struct {
	Nickname   string `json:"nickname"`
	OS         string `json:"os"`
	Device     string `json:"device"`
	DeviceType string `json:"deviceType"`
	Name       string `json:"name"`
}

func getAccountTerminals(c *Context) error {
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47",
		"Dalvik/2.1.0 (Linux; U; Android 10; Pixel XL Build/QP1A.191005.007.A3)",
		"Mozilla/5.0 (Linux; Android 10; Pixel XL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.116 Mobile Safari/537.36 EdgA/46.03.4.5155",
		"Mozilla/5.0 (Linux; Android 10; Pixel XL Build/QP1A.191005.007.A3; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.120 MQQBrowser/6.2 TBS/045713 Mobile Safari/537.36 MMWEBID/6006 MicroMessenger/8.0.3.1880(0x2800033F) Process/tools WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64",
	}

	terminals := make([]terminal, 0)

	for _, s := range userAgents {
		ua := ua.Parse(s)
		var t terminal
		t.Nickname = New4BitID()
		t.Name = ua.Name
		t.OS = ua.OS
		t.Device = ua.Device

		if 0 == len(t.Device) {
			if ua.Mobile {
				t.Device = "Mobile"
			} else if ua.Tablet {
				t.Device = "Tablet"
			} else if ua.Desktop {
				t.Device = "Desktop"
			} else if ua.Bot {
				t.Device = "Bot"
			}
		}

		if ua.Mobile {
			t.DeviceType = "smartphone"
		} else if ua.Tablet {
			t.DeviceType = "tablet"
		} else if ua.Desktop {
			t.DeviceType = "desktop"
		} else if ua.Bot {
			t.DeviceType = "devices"
		}

		if 0 == len(t.DeviceType) {
			t.DeviceType = "smartphone"
		}

		if "Dalvik" == t.Name {
			t.Name = "KnowlGraph client"
		}

		terminals = append(terminals, t)
	}

	return c.Ok(&terminals)
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

	// Main 值不可向客户端暴露，
	// 因为 Main 是授权方的 token。
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
