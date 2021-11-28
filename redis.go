package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/go-redis/redis/v8"
)

// RChallenge is challenge key in redis
// It is Terminal
func RChallenge(key string) string {
	return "challenge_" + key
}

// RToken is token key in redis
// It is int, is user id
func RToken(key string) string {
	return "token_" + key
}

// RTerminal is terminal key in redis
func RTerminal(key string) string {
	return "terminal_" + key
}

// RUser is terminal key in redis
// It is map[int]string, map's key is terminal id, and map's value is terminal token
func RUser(userID int) string {
	return "user_" + strconv.Itoa(userID)
}

// SetWebAuthnSession is set webauthn session data
func SetWebAuthnSession(challenge string, seesinData *webauthn.SessionData) error {
	return SetV2Redis("webauthn_"+challenge, seesinData, ExpireTimeChallengeConfirm)
}

// GetWebAuthnSession returns webauthn session data
func GetWebAuthnSession(challenge string) (seesinData webauthn.SessionData, err error) {
	err = GetV4Redis("webauthn_"+challenge, &seesinData)
	return
}

// SetV2RedisPipe encodes the value as json, and
// then stores it in redis in the form of key-value pairs
func SetV2RedisPipe(pipe redis.Pipeliner, key string, v interface{}, duration time.Duration) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return pipe.Set(ctx, key, string(bs), duration).Err()
}

// SetV2Redis encodes the value as json, and
// then stores it in redis in the form of key-value pairs
func SetV2Redis(key string, v interface{}, duration time.Duration) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, string(bs), duration).Err()
}

// GetV4Redis reads the value of the specified key from redis and
// decodes it into the v interface{}
func GetV4Redis(key string, v interface{}) error {
	result, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), v)
}
