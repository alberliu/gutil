package redis

import (
	"time"

	"github.com/go-redis/redis"
)

func Lock(redisCli *redis.Client, key string, duration time.Duration) bool {
	ok, err := redisCli.SetNX(key, "0", duration).Result()
	if err != nil {
		return false
	}
	if ok {
		return true
	}
	return false
}

func Unlock(redisCli *redis.Client, key string) error {
	_, err := redisCli.Del(key).Result()
	return err
}
