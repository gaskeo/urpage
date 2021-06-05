package redis_api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

var rdb *redis.Client

func Connect(address string, password string, db int) (bool, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	pong, err := rdb.Ping(ctx).Result()

	return pong == "PONG", err
}

func Set(key string, value string, expiredDate time.Time) error {
	var expiredTime time.Duration

	zeroTime := time.Time{}

	if expiredDate == zeroTime {
		expiredTime = time.Duration(0)
	} else {
		expiredTime = expiredDate.Sub(time.Now())
	}
	fmt.Println(expiredTime)

	_, err := rdb.Set(ctx, key, value, expiredTime).Result()

	return err
}

func Get(key string) (string, error) {
	result, err := rdb.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return result, nil
}
