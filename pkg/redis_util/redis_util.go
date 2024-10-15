package redis_util

import (
	"assigement_wallet/config"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	Client      *redis.Client
	LockTimeout = errors.New("timeout when obtain lock")
)

func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.ServerConfigs.Redis.HostIP, config.ServerConfigs.Redis.Port),
		Password: config.ServerConfigs.Redis.Password,
	})
	_, err := Client.Ping().Result()
	return err
}

type RedisLock struct {
	c            *redis.Client
	key          string
	expireSecond int64
}

// UnLock
func (l *RedisLock) UnLock() {
	l.c.Del(l.key)
}

// GetLockWithTimeout get a distribute lock
func GetLockWithTimeout(key string, lockerExpireTime time.Duration, maxWaitTime time.Duration) (lock *RedisLock, err error) {
	t := time.Now()
	for time.Now().Sub(t) < maxWaitTime {
		success, err := Client.SetNX(key, 1, lockerExpireTime).Result()
		if err != nil {
			return nil, err
		}
		//try again if fail
		if !success {
			time.Sleep(time.Millisecond * 50)
			continue
		}
		return &RedisLock{
			c:            Client,
			key:          key,
			expireSecond: int64(lockerExpireTime.Seconds()),
		}, err
	}
	return nil, LockTimeout
}
