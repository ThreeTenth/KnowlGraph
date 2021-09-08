package main

import (
	"encoding/json"
	"time"
)

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
