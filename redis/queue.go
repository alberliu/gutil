package redis

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// AddTask 添加任务，timestamp是任务执行的时间戳
func AddTask(redisCli *redis.Client, key, data string, timestamp int64) (int64, error) {
	result, err := redisCli.ZAdd(key, redis.Z{
		Score:  float64(timestamp),
		Member: data,
	}).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// CancelTask 取消任务
func CancelTask(redisCli *redis.Client, key, data string) error {
	_, err := redisCli.ZRem(key, data).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetTasks 获取可执行的任务
func GetTasks(redisCli *redis.Client, key string) ([]string, bool) {
	// 任务加锁，通知保证只有一个任务在执行
	locked := Lock(redisCli, key+"_lock", 5*time.Minute)
	if !locked {
		return nil, false
	}
	defer Unlock(redisCli, key+"_lock")

	// 获取可执行任务
	max := strconv.FormatInt(time.Now().Unix(), 10)
	datas, err := redisCli.ZRangeByScore(key,
		redis.ZRangeBy{Min: "-inf", Max: max}).Result()
	if err != nil {
		return nil, false
	}
	if len(datas) == 0 {
		return nil, false
	}

	// 删除可执行任务
	_, err = redisCli.ZRemRangeByScore(key, "-inf", max).Result()
	if err != nil {
		return nil, false
	}
	return datas, true
}

// ConsumeTask 消费任务
func ConsumeTask(redisCli *redis.Client, key string, f func(data string)) {
	for {
		datas, ok := GetTasks(redisCli, key)
		if !ok {
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		for _, data := range datas {
			f(data)
		}
	}
}
