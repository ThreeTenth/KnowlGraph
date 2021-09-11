package main

import (
	"math/rand"
	"net/http"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func setAuthorization(c *gin.Context, userID int) error {
	_token := New64BitID()

	err := rdb.Set(ctx, _token, userID, ExpireTimeToken).Err()
	if err != nil {
		return err
	}

	c.SetCookie(CookieToken, _token, int(30*24*3600), "/", "", !config.Debug, false)
	return nil
}

func authentication(c *gin.Context) {
	_token := getRequestToken(c)

	if _token != "" {
		// radis 中 token 对应的值，有三种情况：
		// 1. 空闲状态，对应的值是 TokenStateIdle;
		// 2. 已激活状态，对应的值是 TokenStateActivated;
		// 3. 已授权状态，对应的值才是 user id.
		// 具体可见 [TokenStateIdle], [TokenStateActivated], [TokenStateAuthorized]
		if userID, err := rdb.Get(ctx, _token).Int(); err == nil && TokenStateActivated < userID {
			c.Set(GinKeyUserID, userID)
		}
	}

	c.Next()
}

func deauthorize(c *gin.Context) {
	_token := getRequestToken(c)

	if _token != "" {
		if rdb.Del(ctx, _token).Err() != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.SetCookie(CookieToken, "", -1, "/", "", !config.Debug, false)
	c.Next()
}

func authorizeRequired(c *gin.Context) {
	_token := getRequestToken(c)

	if _token != "" {
		if userID, err := rdb.Get(ctx, _token).Int(); err == nil && 0 < userID {
			c.Set(GinKeyUserID, userID)
			c.Next()
			return
		}
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

func getRequestToken(c *gin.Context) (_token string) {
	_token = c.GetHeader(HeaderAuthorization)
	if _token == "" {
		_token, _ = c.Cookie(CookieToken)
	}
	return
}

const digits = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz_."

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandNdigMbitString returns a randomly generated string with n digits and m base,
// the string range is: a-z, A-Z, 0-9 and'_','.' symbols
// see: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go?answertab=votes#tab-top
func RandNdigMbitString(n int, m ...int) string {
	b := make([]byte, n)
	dict := digits
	if 1 == len(m) {
		dict = digits[0:m[0]]
	} else if 2 == len(m) {
		dict = digits[m[0] : m[0]+m[1]]
	}

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(dict) {
			b[i] = dict[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// New128BitID Get a 128-base random string,
// 10 numbers + 26 lowercase letters + 26 uppercase letters + (=, _), length 64
func New128BitID() string {
	return RandNdigMbitString(128)
}

// New64BitID Get a 64-base random string,
// 10 numbers + 26 lowercase letters + 26 uppercase letters + (=, _), length 64
func New64BitID() string {
	return RandNdigMbitString(64)
}

// New32BitID Get a 32-base random string,
// 10 numbers + 26 uppercase letters
func New32BitID() string {
	return RandNdigMbitString(32, 36)
}

// New32bitID Get a 32-base random string,
// 10 numbers + 26 lowercase letters
func New32bitID() string {
	return RandNdigMbitString(32, 26, 36)
}

// New16BitID Get a 16-base random string,
// 10 numbers + 26 uppercase letters
func New16BitID() string {
	return RandNdigMbitString(16, 36)
}

// New16bitID Get a 16-base random string,
// 10 numbers + 26 lowercase letters
func New16bitID() string {
	return RandNdigMbitString(16, 26, 36)
}

// New4BitID Get a 4-base random string,
// 10 numbers + 26 uppercase letters
func New4BitID() string {
	return RandNdigMbitString(4, 36)
}

// New4bitID Get a 4-base random string,
// 10 numbers + 26 lowercase letters
func New4bitID() string {
	return RandNdigMbitString(4, 26, 36)
}
